package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock dependencies
type mockUserService struct {
	mock.Mock
}

func TestUserSearchHandler(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		query          string
		setupMock      func(*mockUserService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:  "successful search",
			query: "test_user",
			setupMock: func(m *mockUserService) {
				m.On("Search", mock.Anything, "test_user").Return([]User{
					{ID: 1, Name: "test_user"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []User{
				{ID: 1, Name: "test_user"},
			},
		},
		{
			name:  "empty query",
			query: "",
			setupMock: func(m *mockUserService) {
				m.On("Search", mock.Anything, "").Return(nil, ErrInvalidQuery)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: ErrorResponse{
				Error: "invalid search query",
			},
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSvc := new(mockUserService)
			tt.setupMock(mockSvc)
			handler := NewUserSearchHandler(mockSvc)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/users/search?q="+tt.query, nil)
			w := httptest.NewRecorder()

			// Execute request
			handler.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockSvc.AssertExpectations(t)
		})
	}
}

// Performance test
func BenchmarkUserSearch(b *testing.B) {
	mockSvc := new(mockUserService)
	mockSvc.On("Search", mock.Anything, "test_user").Return([]User{
		{ID: 1, Name: "test_user"},
	}, nil)

	handler := NewUserSearchHandler(mockSvc)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/users/search?q=test_user", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}
