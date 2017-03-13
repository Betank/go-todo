package main

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func main() {

	mux := mux.NewRouter()
	mux.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("/", ServeIndexHTML)

	http.Handle("/", mux)
	http.ListenAndServe(":8080", nil)
}

// ServeIndexHTML renders and serves page
func ServeIndexHTML(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(filepath.Join("templates", "index.tmpl")))
	t.Execute(w, nil)
}
