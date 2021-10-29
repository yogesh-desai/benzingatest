package api

import "net/http"

// GetHealth is a simple health check endpoint
func GetHealth(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("OK")) // This call will automatically sets HTTP 200-OK as header response and body will be OK

}

// HandleLog handles the logs as per env variables
func HandleLog(w http.ResponseWriter, r *http.Request) {

}
