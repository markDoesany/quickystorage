package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/markDoesany/quickymessenger/database"
	"github.com/markDoesany/quickymessenger/handlers"
	"github.com/markDoesany/quickymessenger/services"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	database.InitDB()

	// Set up the persistent menu
	if err := services.SetupPersistentMenu(); err != nil {
		log.Printf("Warning: Could not set up persistent menu: %v", err)
	} else {
		log.Println("Persistent menu set up successfully!")
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/", handlers.Webhook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	srv := &http.Server{
		Handler: handler,
		Addr:    "localhost:" + port,
	}

	log.Printf("HTTP server listening at %v", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
