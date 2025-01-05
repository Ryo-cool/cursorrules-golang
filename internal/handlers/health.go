package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

// HealthStatus represents the system health information
type HealthStatus struct {
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	DBStatus    string    `json:"db_status"`
	SystemStats struct {
		GoRoutines  int     `json:"goroutines"`
		MemoryUsage float64 `json:"memory_usage_mb"`
		CPUUsage    float64 `json:"cpu_usage_percent"`
	} `json:"system_stats"`
}

// HealthCheckHandler handles health check requests
func HealthCheckHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := HealthStatus{
			Status:    "healthy",
			Timestamp: time.Now(),
		}

		// Check database connection
		if err := db.Ping(); err != nil {
			status.Status = "degraded"
			status.DBStatus = "unavailable"
		} else {
			status.DBStatus = "connected"
		}

		// Collect system metrics
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		status.SystemStats.GoRoutines = runtime.NumGoroutine()
		status.SystemStats.MemoryUsage = float64(memStats.Alloc) / 1024 / 1024
		// Note: CPU usage calculation would require additional implementation

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}
