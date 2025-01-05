package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/cursorrules-golang/internal/cache"
	"example.com/cursorrules-golang/internal/models"
	"github.com/stretchr/testify/assert"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func TestSearchUsersHandler(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		searchBy       string
		page           string
		pageSize       string
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "successful search by name",
			query:          "test_user",
			searchBy:       "name",
			page:           "1",
			pageSize:       "10",
			expectedStatus: http.StatusOK,
			expectedBody: models.PaginatedResponse{
				Data: []models.User{
					{ID: 1, Name: "test_user", Email: "test@example.com", Age: 25},
				},
				Pagination: struct {
					CurrentPage int  `json:"current_page"`
					PageSize    int  `json:"page_size"`
					TotalItems  int  `json:"total_items"`
					TotalPages  int  `json:"total_pages"`
					HasNext     bool `json:"has_next"`
					HasPrevious bool `json:"has_previous"`
				}{
					CurrentPage: 1,
					PageSize:    10,
					TotalItems:  1,
					TotalPages:  1,
					HasNext:     false,
					HasPrevious: false,
				},
			},
		},
		{
			name:           "invalid search parameters",
			query:          "",
			searchBy:       "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBody: ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid search parameters",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test DB and cache
			db := createTestDB(t)
			cache := createTestCache(t)

			handler := SearchUsersHandler(db, cache)

			// Create request with query parameters
			url := "/users/search?search=" + tt.query
			if tt.searchBy != "" {
				url += "&search_by=" + tt.searchBy
			}
			if tt.page != "" {
				url += "&page=" + tt.page
			}
			if tt.pageSize != "" {
				url += "&page_size=" + tt.pageSize
			}

			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			// Execute request
			handler.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			// Compare response with expected body
			if tt.expectedStatus == http.StatusOK {
				// For successful responses, compare specific fields
				respMap := response.(map[string]interface{})
				assert.NotNil(t, respMap["data"])
				assert.NotNil(t, respMap["pagination"])
			} else {
				// For error responses, compare the entire structure
				assert.Equal(t, tt.expectedBody, response)
			}
		})
	}
}

// BenchmarkSearchUsers runs performance tests for the search handler
func BenchmarkSearchUsers(b *testing.B) {
	// Setup test environment
	db := createTestDB(b)
	cache := createTestCache(b)
	handler := SearchUsersHandler(db, cache)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/users/search?search=test&search_by=name", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}

// Helper functions for testing
func createTestDB(t testing.TB) *sql.DB {
	// Create an in-memory SQLite database for testing
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Create test tables and data
	if err := setupTestDB(db); err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}

	return db
}

func createTestCache(_ testing.TB) *cache.Cache {
	return cache.New(cache.Config{
		MaxItems:        100,
		CleanupInterval: time.Minute,
	})
}

func setupTestDB(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			age INTEGER NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	// Insert test data
	_, err = db.Exec(`
		INSERT INTO users (name, email, age) VALUES
		('test_user', 'test@example.com', 25)
	`)
	return err
}
