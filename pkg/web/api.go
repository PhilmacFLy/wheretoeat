package web

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func mainAPIHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Main API Handler")
}

func listVenuesAPIHandler(w http.ResponseWriter, r *http.Request) {

}

func getVenueAPIHandler(w http.ResponseWriter, r *http.Request) {

}

func postVenueAPIHandler(w http.ResponseWriter, r *http.Request) {

}

func patchVenueAPIHander(w http.ResponseWriter, r *http.Request) {

}

func deleteVenueAPIHandler(w http.ResponseWriter, r *http.Request) {

}

func addVisitAPIHandler(w http.ResponseWriter, r *http.Request) {

}

func getAPIRouter(prefix string) *mux.Router {
	r := mux.NewRouter().PathPrefix(prefix).Subrouter()
	r.HandleFunc("/", mainAPIHandler)
	r.HandleFunc("/venue", postVenueAPIHandler).Methods("POST")
	r.HandleFunc("/venue/list", listVenuesAPIHandler).Methods("GET")
	r.HandleFunc("/venue/{ID}", getVenueAPIHandler).Methods("GET")
	r.HandleFunc("/venue/{ID}", patchVenueAPIHander).Methods("PATCH")
	r.HandleFunc("/venue/{ID}", deleteVenueAPIHandler).Methods("DELETE")
	r.HandleFunc("/venue/{ID}/addvisits", addVisitAPIHandler).Methods("POST")

	return r
}
