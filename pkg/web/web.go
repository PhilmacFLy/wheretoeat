package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui/", http.StatusSeeOther)
}

//GetRouters gives back the router for the api aswell as the frontend
func GetRouters(prefix string) *mux.Router {
	r := mux.NewRouter().PathPrefix(prefix).Subrouter()
	ui := getUIRouter("/ui")
	r.PathPrefix("/ui").Handler(ui)
	api := getAPIRouter("/api")
	r.PathPrefix("/api").Handler(api)
	r.HandleFunc("/", mainHandler)

	return r
}
