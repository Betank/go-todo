package main

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type storage struct {
	sync.Mutex
	todos map[int]*todo
}

type todo struct {
	ID        int
	Title     string
	Completed bool
}

var entries = &storage{todos: make(map[int]*todo)}
var index int

func serveIndexHTML(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("index").Parse(indexHTML))
	t.Execute(w, entries.todos)
}

func showActive(w http.ResponseWriter, r *http.Request) {
	activeTodos := make([]todo, 0)
	for _, todo := range entries.todos {
		if !todo.Completed {
			activeTodos = append(activeTodos, *todo)
		}
	}

	t := template.Must(template.New("index").Parse(indexHTML))
	t.Execute(w, activeTodos)
}

func showCompleted(w http.ResponseWriter, r *http.Request) {
	completedTodos := make([]todo, 0)
	for _, todo := range entries.todos {
		if todo.Completed {
			completedTodos = append(completedTodos, *todo)
		}
	}

	t := template.Must(template.New("index").Parse(indexHTML))
	t.Execute(w, completedTodos)
}

func clearTodos(w http.ResponseWriter, r *http.Request) {
	entries.Lock()
	defer entries.Unlock()
	entries.todos = make(map[int]*todo)

	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	entries.Lock()
	defer entries.Unlock()
	index++
	entries.todos[index] = &todo{ID: index, Title: r.FormValue("text")}

	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	entries.Lock()
	defer entries.Unlock()

	if id, err := strconv.Atoi(mux.Vars(r)["id"]); err == nil {
		entries.todos[id].Completed = !entries.todos[id].Completed
	}

	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	entries.Lock()
	defer entries.Unlock()

	if id, err := strconv.Atoi(mux.Vars(r)["id"]); err == nil {
		delete(entries.todos, id)
	}

	w.Header().Add("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/add", addTodo).Methods("POST")
	mux.HandleFunc("/active", showActive).Methods("GET")
	mux.HandleFunc("/completed", showCompleted).Methods("GET")
	mux.HandleFunc("/clear", clearTodos).Methods("GET")
	mux.HandleFunc("/update/{id}", updateTodo).Methods("GET")
	mux.HandleFunc("/delete/{id}", deleteTodo).Methods("GET")
	mux.HandleFunc("/", serveIndexHTML)

	http.Handle("/", mux)
	http.ListenAndServe(":8080", nil)
}
