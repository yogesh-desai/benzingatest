package api

import (
	"log"
	"net/http"
)

func Init() {

	http.Handle("/healthz", http.HandlerFunc(GetHealth))

	http.Handle("/log", http.HandlerFunc(HandleLog))

	log.Println("[Init] Api routes initialized.")

}
