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

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(100, 1000) // 100 requests per second, bucket size 1000

	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/users", handlers.UsersHandler(db))
	mux.HandleFunc("/users/", handlers.UserHandler(db))
	mux.HandleFunc("/health", handlers.HealthCheckHandler(db))

	// Apply middleware chain
	handler := middleware.Logging(
		rateLimiter.RateLimit(
			middleware.AuthMiddleware(mux),
		),
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
