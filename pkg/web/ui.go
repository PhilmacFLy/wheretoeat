package web

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func mainUIHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Main UI Handler")
}

func getUIRouter(prefix string) *mux.Router {
	r := mux.NewRouter().PathPrefix(prefix).Subrouter()
	r.HandleFunc("/", mainUIHandler)
	return r
}
