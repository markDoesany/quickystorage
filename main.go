package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/markDoesany/quickymessenger/database"
	"github.com/markDoesany/quickymessenger/handlers"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	database.InitDB()
	handler := http.NewServeMux()
	handler.HandleFunc("/", handlers.Webhook)

	srv := &http.Server{
		Handler: handler,
		Addr:    "localhost:5000",
	}

	log.Printf("HTTP server listening at %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
