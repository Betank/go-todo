package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/gorilla/mux"
)

type storage struct {
	sync.Mutex
	todos []todo
}

type todo struct {
	ID        string
	Title     string
	Completed bool
	Checked   bool
}

var entries *storage

func main() {
	entries = &storage{todos: make([]todo, 0)}

	mux := mux.NewRouter()
	mux.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("/", ServeIndexHTML)

	http.Handle("/", mux)
	http.ListenAndServe(":8080", nil)
}

// ServeIndexHTML renders and serves page
func ServeIndexHTML(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles(filepath.Join("templates", "index.tmpl")))

	entries.Lock()
	defer entries.Unlock()
	t.Execute(w, entries.todos)
}
