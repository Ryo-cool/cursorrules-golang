package models

// QueryParams represents common query parameters for API endpoints
type QueryParams struct {
	// Search parameters
	Search   string `json:"search"`
	SearchBy string `json:"search_by"`
	MinAge   *int   `json:"min_age,omitempty"`
	MaxAge   *int   `json:"max_age,omitempty"`

	// Pagination parameters
	Page     int `json:"page"`
	PageSize int `json:"page_size"`

	// Sorting parameters
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"` // asc or desc
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination struct {
		CurrentPage int  `json:"current_page"`
		PageSize    int  `json:"page_size"`
		TotalItems  int  `json:"total_items"`
		TotalPages  int  `json:"total_pages"`
		HasNext     bool `json:"has_next"`
		HasPrevious bool `json:"has_previous"`
	} `json:"pagination"`
}

// NewQueryParams creates a new QueryParams with default values
func NewQueryParams() QueryParams {
	return QueryParams{
		Page:      1,
		PageSize:  10,
		SortBy:    "id",
		SortOrder: "asc",
	}
}
