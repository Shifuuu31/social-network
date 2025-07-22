package utils

import (
	"net/http"

	uuid "github.com/satori/go.uuid"
)

// Create unique ID
func UniqueId() string {
	return uuid.NewV4().String()
}

func ConfigHeader(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	return w
}

func ConfigFSHeader(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8080")
		fs.ServeHTTP(w, r)
	}
}
