package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var datafolder string

func mainHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui/", http.StatusSeeOther)
}

//SetupRouters gives back the router for the api aswell as the frontend
func SetupRouters(prefix string, venuedatafolder string) *mux.Router {
	datafolder = venuedatafolder + string(os.PathSeparator)
	r := mux.NewRouter().PathPrefix(prefix).Subrouter()
	ui := getUIRouter("/ui")
	r.PathPrefix("/ui").Handler(ui)
	api := getAPIRouter("/api")
	r.PathPrefix("/api").Handler(api)
	r.HandleFunc("/", mainHandler)
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	})
	return r
}
