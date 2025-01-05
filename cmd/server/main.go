package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/cursorrules-golang/internal/cache"
	"example.com/cursorrules-golang/internal/database"
	"example.com/cursorrules-golang/internal/handlers"
	"example.com/cursorrules-golang/internal/metrics"
	"example.com/cursorrules-golang/internal/middleware"
)

func main() {
	// Initialize components
	db := database.InitDB()
	defer db.Close()

	// Initialize cache with configuration
	cacheConfig := cache.Config{
		MaxItems:        10000,
		CleanupInterval: 5 * time.Minute,
	}
	cache := cache.New(cacheConfig)

	// Initialize metrics
	metrics := metrics.GetMetrics()

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(100, 1000) // 100 requests per second, bucket size 1000

	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/users", handlers.UsersHandler(db))
	mux.HandleFunc("/users/", handlers.UserHandler(db))
	mux.HandleFunc("/users/search", handlers.SearchUsersHandler(db, cache))
	mux.HandleFunc("/health", handlers.HealthCheckHandler(db))
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// Combine application metrics with cache stats
		stats := metrics.GetSnapshot()
		cacheStats := cache.GetStats()

		// Merge metrics
		response := map[string]interface{}{
			"app_metrics": stats,
			"cache_stats": cacheStats,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

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
