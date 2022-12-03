package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/livepool-io/pool-service/common"
	"github.com/livepool-io/pool-service/db"
	"github.com/livepool-io/pool-service/middleware"
)

// TODO rewrite according to readme and other API routes

func GetTranscoder(w http.ResponseWriter, r *http.Request) {
	// Make sure DB exists
	if err := db.CacheDB(); err != nil {
		common.HandleInternalError(w, err)
		return
	}

	// Handle preflight requests
	middleware.HandlePreflightGET(w, r)

	// Check GET
	if r.Method != "GET" {
		common.HandleBadRequest(w, errors.New("must be a GET request"))
		return
	}

	// Get Transcoder param
	tAddr := r.URL.Query().Get("transcoder")

	if tAddr == "" {
		common.HandleBadRequest(w, errors.New("transcoder is a required parameter"))
		return
	}

	t, err := db.Database.GetTranscoder(tAddr)
	if err != nil {
		common.HandleInternalError(w, err)
		return
	}

	tJSON, err := json.Marshal(t)
	if err != nil {
		common.HandleInternalError(w, err)
		return
	}

	common.HandleOK(w)
	w.Write(tJSON)
}
