package repositories

import (
	"context"
	"gotodo/config"
	"gotodo/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoRepository struct{}

func (r *TodoRepository) List() ([]models.ToDo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := config.TodoCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding todos: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var todos []models.ToDo
	if err := cursor.All(ctx, &todos); err != nil {
		log.Printf("Error decoding todos: %v", err)
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepository) Create(newToDo models.ToDo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.TodoCollection.InsertOne(ctx, newToDo)
	return err
}

func (r *TodoRepository) Get(id primitive.ObjectID) (models.ToDo, error) {
	var toDo models.ToDo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := config.TodoCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&toDo)
	return toDo, err
}

func (r *TodoRepository) Update(id primitive.ObjectID, updatedData map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.TodoCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updatedData})
	return err
}

func (r *TodoRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := config.TodoCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
