package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var TodoCollection *mongo.Collection

func InitializeMongoDB() {
	var err error
	Client, err = mongo.NewClient(options.Client().ApplyURI("mongodb+srv://igaming:hCDmJeZUHD4mcm1y@cluster0.g5mbz7u.mongodb.net/igaming"))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = Client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	TodoCollection = Client.Database("todoApp").Collection("todos")
	log.Println("Connected to MongoDB!")
}
