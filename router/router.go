package routes

import (
	"gotodo/modules/todo/controllers"
	"gotodo/modules/todo/repositories"
	"gotodo/modules/todo/services"

	"github.com/go-chi/chi/v5"
)

func InitializeRouter() *chi.Mux {
	r := chi.NewRouter()

	todoRepo := &repositories.TodoRepository{}
	todoService := services.NewTodoService(todoRepo)
	todoController := controllers.NewTodoController(todoService)

	r.Get("/todos", todoController.ListToDos)
	r.Post("/todos", todoController.CreateToDo)
	r.Get("/todos/{id}", todoController.GetToDo)
	r.Put("/todos/{id}", todoController.UpdateToDo)
	r.Delete("/todos/{id}", todoController.DeleteToDo)

	return r
}
