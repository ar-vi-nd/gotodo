package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ToDo struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Description string             `json:"description" bson:"description"`
	Title       string             `json:"title" bson:"title"`
	Completed   bool               `json:"completed" bson:"completed"`
}

var client *mongo.Client
var todoCollection *mongo.Collection

func InitializeMongoDB() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://igaming:hCDmJeZUHD4mcm1y@cluster0.g5mbz7u.mongodb.net/igaming"))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	todoCollection = client.Database("todoApp").Collection("todos")
	fmt.Println("Connected to MongoDB!")
}

func ListToDos(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := todoCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding todos: %v", err)
		http.Error(w, "Could not fetch todos", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var todos []ToDo
	if err = cursor.All(ctx, &todos); err != nil {
		log.Printf("Error decoding todos: %v", err)
		http.Error(w, "Error reading todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func CreateToDo(w http.ResponseWriter, r *http.Request) {
	var newToDo ToDo
	if err := json.NewDecoder(r.Body).Decode(&newToDo); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newToDo.ID = primitive.NewObjectID()
	_, err := todoCollection.InsertOne(ctx, newToDo)
	if err != nil {
		log.Printf("Error inserting todo: %v", err)
		http.Error(w, "Could not create todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newToDo)
}

func GetToDo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var toDo ToDo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = todoCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&toDo)
	if err != nil {
		log.Printf("Error finding todo: %v", err)
		http.Error(w, "To-do not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toDo)
}

func UpdateToDo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedToDo map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updatedToDo); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{"$set": updatedToDo}

	_, err = todoCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		log.Printf("Error updating todo: %v", err)
		http.Error(w, "Could not update todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedToDo)
}

func DeleteToDo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := todoCollection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil || result.DeletedCount == 0 {
		log.Printf("Error deleting todo or todo not found: %v", err)
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

func InitializeRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/todos", func(r chi.Router) {
		r.Get("/", ListToDos)
		r.Post("/", CreateToDo)
		r.Get("/{id}", GetToDo)
		r.Patch("/{id}", UpdateToDo)
		r.Delete("/{id}", DeleteToDo)

	})

	return r

}

func main() {
	InitializeMongoDB()
	r := InitializeRouter()
	log.Fatal(http.ListenAndServe(":3000", r))
}
