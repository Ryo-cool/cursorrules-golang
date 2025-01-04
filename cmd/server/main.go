package main

import (
	"log"
	"net/http"
	"os"

	"example.com/cursorrules-golang/internal/database"
	"example.com/cursorrules-golang/internal/handlers"
	"example.com/cursorrules-golang/internal/middleware"
)

func main() {
	db := database.InitDB()
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/users", handlers.UsersHandler(db))
	mux.HandleFunc("/users/", handlers.UserHandler(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	err := http.ListenAndServe(":"+port, middleware.Logging(mux))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
