package middleware

import "net/http"

func HandlePreflightPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodOptions {
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Write(nil)
}

func HandlePreflightGET(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodOptions {
		return
	}
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age=60, public")
}
