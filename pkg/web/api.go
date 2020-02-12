package web

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func mainAPIHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Main API Handler")
}

func getAPIRouter(prefix string) *mux.Router {
	r := mux.NewRouter().PathPrefix(prefix).Subrouter()
	r.HandleFunc("/", mainAPIHandler)
	return r
}
