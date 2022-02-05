package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/livepool-io/pool-service/common"
	"github.com/livepool-io/pool-service/db"
	"github.com/livepool-io/pool-service/middleware"
	"github.com/livepool-io/pool-service/models"
)

func Jobs(w http.ResponseWriter, r *http.Request) {
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

	// query parameters and create database filter
	query := r.URL.Query()

	// todo require auth if node is defined
	t := query.Get("transcoder")
	n := query.Get("node")
	fromStr := query.Get("from")
	toStr := query.Get("to")

	// If time params aren't defined default to last 24h
	var from int64
	if fromStr == "" {
		from = time.Now().Add(-24 * time.Hour).Unix()
	} else {
		var err error
		from, err = strconv.ParseInt(fromStr, 10, 64)
		if err != nil {
			common.HandleBadRequest(w, err)
		}
	}

	var to int64
	if toStr == "" {
		to = time.Now().Unix()
	} else {
		var err error
		to, err = strconv.ParseInt(toStr, 10, 64)
		if err != nil {
			common.HandleBadRequest(w, err)
		}
	}

	jobs, err := db.Database.GetJobs(t, n, from, to)
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

	var job *models.Job

	// Unmarshal the json, return 400 if error
	if err := json.Unmarshal([]byte(body), job); err != nil {
		common.HandleBadRequest(w, err)
		return
	}

	if err := db.Database.CreateJob(job); err != nil {
		common.HandleInternalError(w, err)
	}

	// TODO: update transcoder pending balance

	// Respond with 200 OK
	common.HandleOK(w)
}
