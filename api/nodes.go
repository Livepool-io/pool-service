package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/livepool-io/pool-service/common"
	"github.com/livepool-io/pool-service/db"
	"github.com/livepool-io/pool-service/middleware"
	"github.com/livepool-io/pool-service/models"
)

func Nodes(w http.ResponseWriter, r *http.Request) {
	// Make sure DB exists
	if err := db.CacheDB(); err != nil {
		common.HandleInternalError(w, err)
		return
	}

	switch r.Method {
	case "GET":
		handleGetNodes(w, r)
	case "POST":
		handlePostNode(w, r)
	default:
		common.HandleBadRequest(w, errors.New("must be a GET or POST request"))
		return
	}
}

func handleGetNodes(w http.ResponseWriter, r *http.Request) {
	middleware.HandlePreflightGET(w, r)

	query := r.URL.Query()

	t := query.Get("transcoder")
	region := query.Get("region")

	if !middleware.IsECDSAAuthorized(t, r.Header.Get("Authorization"), []byte(query.Encode())) {
		common.HandleUnauthorized(w, errors.New("ECDSA authentication failed"))
		return
	}

	// get from DB
	nodes, err := db.Database.GetNodes(t, region)
	if err != nil {
		common.HandleInternalError(w, err)
		return
	}

	nodesJSON, err := json.Marshal(nodes)
	if err != nil {
		common.HandleInternalError(w, err)
		return
	}

	// return data and 200 OK
	common.HandleOK(w)
	w.Write(nodesJSON)
}

func handlePostNode(w http.ResponseWriter, r *http.Request) {
	middleware.HandlePreflightPOST(w, r)

	// Read request body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		common.HandleBadRequest(w, err)
		return
	}

	// Check client authorisation (HMAC shared secret)
	if ok := middleware.IsHMACAuthorized(
		r.Header.Get("Authorization"),
		body,
	); !ok {
		common.HandleUnauthorized(w, errors.New("request authentication unsuccesful"))
		return
	}

	var node *models.Node

	// Unmarshal the json, return 400 if error
	if err := json.Unmarshal([]byte(body), node); err != nil {
		common.HandleBadRequest(w, err)
		return
	}

	if err := db.Database.AddNode(node); err != nil {
		common.HandleInternalError(w, err)
		return
	}
	// Respond with 200 OK
	common.HandleOK(w)

}
