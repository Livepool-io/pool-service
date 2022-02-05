package common

import (
	"net/http"
)

func HandleInternalError(w http.ResponseWriter, err error) {
	RespondWithError(w, err, http.StatusInternalServerError)
}

func HandleBadRequest(w http.ResponseWriter, err error) {
	RespondWithError(w, err, http.StatusBadRequest)
}

func HandleUnauthorized(w http.ResponseWriter, err error) {
	RespondWithError(w, err, http.StatusUnauthorized)
}

func RespondWithError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

func HandleOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
