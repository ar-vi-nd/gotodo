package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ToDo struct {
	ID          int    `json:"id"`
	DESCRIPTION string `json:"description"`
	Title       string `json:"title"`
	Completed   bool   `json:"completed"`
}

var todos = []ToDo{}
var nextID = 1

func FindToDo(id int) (*ToDo, error) {
	for i, t := range todos {
		if t.ID == id {
			return &todos[i], nil
		}
	}
	return nil, errors.New("to-do not found")
}

func ListToDos(w http.ResponseWriter, r *http.Request) {
	println("in this function")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func CreateToDo(w http.ResponseWriter, r *http.Request) {
	var newToDo ToDo

	// var printTodo ToDo
	// something := json.NewDecoder(r.Body)
	// something.Decode(&printTodo)
	// fmt.Printf(err.Error())
	// fmt.Printf(printTodo.Title)
	// fmt.Printf(something)
	// fmt.Printf(something.Decode(&printTodo))

	if err := json.NewDecoder(r.Body).Decode(&newToDo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newToDo.ID = nextID
	nextID++
	todos = append(todos, newToDo)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newToDo)
}

func GetToDo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	toDo, err := FindToDo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toDo)
}

func UpdateToDo(w http.ResponseWriter, r *http.Request) {
	params := chi.URLParam(r, "id")
	id, err := strconv.Atoi(params)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	toDo, err := FindToDo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(toDo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toDo)
}

func DeleteToDo(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			response := map[string]interface{}{
				"message": "To-do deleted successfully",
				"status":  "success",
				"toDo":    t, // optional: include the deleted to-do item
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK) // Set the status code to 200 OK
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	http.Error(w, "To-do not found", http.StatusNotFound)
}

func InitializeRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/todos", func(r chi.Router) {
		r.Get("/", ListToDos)         // GET /todos
		r.Post("/", CreateToDo)       // POST /todos
		r.Get("/{id}", GetToDo)       // GET /todos/{id}
		r.Put("/{id}", UpdateToDo)    // PUT /todos/{id}
		r.Delete("/{id}", DeleteToDo) // DELETE /todos/{id}
	})

	return r
}

func main() {
	r := InitializeRouter()
	http.ListenAndServe(":3000", r)
}
