package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"example.com/cursorrules-golang/internal/cache"
	"example.com/cursorrules-golang/internal/metrics"
	"example.com/cursorrules-golang/internal/models"
)

// SearchUsersHandler handles user search requests with pagination
func SearchUsersHandler(db *sql.DB, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		metrics := metrics.GetMetrics()

		// Parse query parameters
		query := r.URL.Query()
		params := models.NewQueryParams()

		if page := query.Get("page"); page != "" {
			fmt.Sscanf(page, "%d", &params.Page)
		}
		if pageSize := query.Get("page_size"); pageSize != "" {
			fmt.Sscanf(pageSize, "%d", &params.PageSize)
		}

		params.Search = query.Get("search")
		params.SearchBy = query.Get("search_by")
		params.SortBy = query.Get("sort_by")
		params.SortOrder = strings.ToLower(query.Get("sort_order"))

		// Try to get from cache first
		cacheKey := fmt.Sprintf("users:search:%v", params)
		if cached, found := cache.Get(cacheKey); found {
			if response, ok := cached.(models.PaginatedResponse); ok {
				metrics.RecordRequest(time.Since(start), true)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
				return
			}
		}

		// Build the SQL query
		baseQuery := "SELECT id, name, email, age FROM users WHERE 1=1"
		countQuery := "SELECT COUNT(*) FROM users WHERE 1=1"
		var conditions []string
		var args []interface{}

		if params.Search != "" {
			switch params.SearchBy {
			case "name":
				conditions = append(conditions, "name LIKE ?")
				args = append(args, "%"+params.Search+"%")
			case "email":
				conditions = append(conditions, "email LIKE ?")
				args = append(args, "%"+params.Search+"%")
			}
		}

		if len(conditions) > 0 {
			whereClause := " AND " + strings.Join(conditions, " AND ")
			baseQuery += whereClause
			countQuery += whereClause
		}

		// Add sorting
		if params.SortBy != "" {
			baseQuery += fmt.Sprintf(" ORDER BY %s %s", params.SortBy,
				strings.ToUpper(params.SortOrder))
		}

		// Add pagination
		offset := (params.Page - 1) * params.PageSize
		baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", params.PageSize, offset)

		// Get total count
		var totalItems int
		err := db.QueryRow(countQuery, args...).Scan(&totalItems)
		if err != nil {
			metrics.RecordRequest(time.Since(start), false)
			http.Error(w, "Failed to count users", http.StatusInternalServerError)
			return
		}

		// Execute the main query
		rows, err := db.Query(baseQuery, args...)
		if err != nil {
			metrics.RecordRequest(time.Since(start), false)
			http.Error(w, "Failed to query users", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age); err != nil {
				metrics.RecordRequest(time.Since(start), false)
				http.Error(w, "Failed to scan user", http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		// Prepare paginated response
		totalPages := (totalItems + params.PageSize - 1) / params.PageSize
		response := models.PaginatedResponse{
			Data: users,
		}
		response.Pagination.CurrentPage = params.Page
		response.Pagination.PageSize = params.PageSize
		response.Pagination.TotalItems = totalItems
		response.Pagination.TotalPages = totalPages
		response.Pagination.HasNext = params.Page < totalPages
		response.Pagination.HasPrevious = params.Page > 1

		// Cache the response
		cache.Set(cacheKey, response, 5*time.Minute)

		metrics.RecordRequest(time.Since(start), true)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
