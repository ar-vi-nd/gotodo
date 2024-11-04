package main

import (
	"gotodo/config"
	routes "gotodo/router"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize MongoDB connection
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	config.InitializeMongoDB()

	// Initialize router and start server
	r := routes.InitializeRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // default port if not specified
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}
