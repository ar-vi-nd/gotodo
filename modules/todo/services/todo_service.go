package services

import (
	"gotodo/models"
	"gotodo/modules/todo/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoService struct {
	repo *repositories.TodoRepository
}

func NewTodoService(repo *repositories.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) ListToDos() ([]models.ToDo, error) {
	return s.repo.List()
}

func (s *TodoService) CreateToDo(newToDo models.ToDo) error {
	return s.repo.Create(newToDo)
}

func (s *TodoService) GetToDoByID(id string) (models.ToDo, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.ToDo{}, err
	}
	return s.repo.Get(objID)
}

func (s *TodoService) UpdateToDoByID(id string, updatedData map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repo.Update(objID, updatedData)
}

func (s *TodoService) DeleteToDoByID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(objID)
}
