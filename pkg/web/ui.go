package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func showtemplate(w http.ResponseWriter, path string, data interface{}) {
	t, err := template.ParseFiles(path)
	if err != nil {
		fmt.Fprintln(w, "Error parsing template:", err)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Fprintln(w, "Error executing template:", err)
		return
	}
}

func mainUIHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Main UI Handler")
	tp := "../../web/templates/main.html"
	showtemplate(w, tp, nil)
}

func getUIRouter(prefix string) *mux.Router {
	r := mux.NewRouter().PathPrefix(prefix).Subrouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/ui/static/", http.FileServer(http.Dir("../../web/static/"))))
	r.HandleFunc("/", mainUIHandler)
	return r
}
