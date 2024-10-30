package controllers

import (
	"encoding/json"
	"gotodo/models"
	"gotodo/modules/todo/services"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoController struct {
	service *services.TodoService
}

func NewTodoController(service *services.TodoService) *TodoController {
	return &TodoController{service: service}
}

func (c *TodoController) ListToDos(w http.ResponseWriter, r *http.Request) {
	todos, err := c.service.ListToDos()
	if err != nil {
		log.Printf("Error fetching todos: %v", err)
		http.Error(w, "Could not fetch todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (c *TodoController) CreateToDo(w http.ResponseWriter, r *http.Request) {
	var newToDo models.ToDo
	if err := json.NewDecoder(r.Body).Decode(&newToDo); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newToDo.ID = primitive.NewObjectID()
	if err := c.service.CreateToDo(newToDo); err != nil {
		log.Printf("Error inserting todo: %v", err)
		http.Error(w, "Could not create todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newToDo)
}

func (c *TodoController) GetToDo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	toDo, err := c.service.GetToDoByID(id)
	if err != nil {
		log.Printf("Error finding todo: %v", err)
		http.Error(w, "To-do not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toDo)
}

func (c *TodoController) UpdateToDo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var updatedToDo map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updatedToDo); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := c.service.UpdateToDoByID(id, updatedToDo); err != nil {
		log.Printf("Error updating todo: %v", err)
		http.Error(w, "Could not update todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedToDo)
}

func (c *TodoController) DeleteToDo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := c.service.DeleteToDoByID(id); err != nil {
		log.Printf("Error deleting todo: %v", err)
		http.Error(w, "To-do not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"message": "To-do deleted successfully",
		"status":  "success",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
