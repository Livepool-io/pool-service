package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/livepool-io/pool-service/common"
	"github.com/livepool-io/pool-service/db"
	"github.com/livepool-io/pool-service/middleware"
)

func GetTranscoders(w http.ResponseWriter, r *http.Request) {
	// Make sure DB exists
	if err := db.CacheDB(); err != nil {
		common.HandleInternalError(w, err)
		return
	}

	switch r.Method {
	case "GET":
		handleGetJobs(w, r)
	case "POST":
		handlePostJob(w, r)
	default:
		common.HandleBadRequest(w, errors.New("Must be a GET or POST request"))
		return
	}
}

func handleGetJobs(w http.ResponseWriter, r *http.Request) {
	middleware.HandlePreflightGET(w, r)

	// Todo query parameters and create database filter

	jobs, err := db.Database.GetJobs()
	if err != nil {
		common.HandleInternalError(w, err)
		return
	}

	jobsJSON, err := json.Marshal(jobs)
	if err != nil {
		common.HandleInternalError(w, err)
		return
	}

	common.HandleOK(w)
	w.Write(jobsJSON)
}

func handlePostJob(w http.ResponseWriter, r *http.Request) {
	middleware.HandlePreflightPOST(w, r)

	// Read request body
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		common.HandleBadRequest(w, err)
		return
	}

	// Check client authorisation (HMAC shared secret)
	if ok := middleware.IsAuthorized(
		r.Header.Get("Authorization"),
		body,
	); !ok {
		common.HandleUnauthorized(w, errors.New("Request authentication unsuccesful"))
		return
	}

	// TODO: insert into Jobs table
	// TODO: update transcoder pending balance

	// Respond with 200 OK
	common.HandleOK(w)
}
