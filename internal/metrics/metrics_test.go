package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {
	metrics := GetMetrics()
	assert.NotNil(t, metrics)
}

func TestRecordRequest(t *testing.T) {
	metrics := GetMetrics()

	testCases := []struct {
		name     string
		duration time.Duration
		success  bool
	}{
		{
			name:     "successful request",
			duration: 100 * time.Millisecond,
			success:  true,
		},
		{
			name:     "failed request",
			duration: 1 * time.Second,
			success:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			metrics.RecordRequest(tc.duration, tc.success)

			// Get snapshot and verify
			snapshot := metrics.GetSnapshot()
			assert.Greater(t, snapshot.TotalRequests, uint64(0))

			if tc.success {
				assert.Greater(t, snapshot.SuccessfulRequests, uint64(0))
			} else {
				assert.Greater(t, snapshot.FailedRequests, uint64(0))
			}

			// Verify response time metrics
			assert.GreaterOrEqual(t, snapshot.MaxResponseTime, snapshot.MinResponseTime)
			assert.GreaterOrEqual(t, snapshot.AverageResponseTime, snapshot.MinResponseTime)
			assert.LessOrEqual(t, snapshot.AverageResponseTime, snapshot.MaxResponseTime)
		})
	}
}

func TestMetricsSnapshot(t *testing.T) {
	metrics := GetMetrics()

	// Record some test data
	metrics.RecordRequest(100*time.Millisecond, true)
	metrics.RecordRequest(200*time.Millisecond, true)
	metrics.RecordRequest(300*time.Millisecond, false)

	snapshot := metrics.GetSnapshot()
	assert.Equal(t, uint64(3), snapshot.TotalRequests)
	assert.Equal(t, uint64(2), snapshot.SuccessfulRequests)
	assert.Equal(t, uint64(1), snapshot.FailedRequests)
	assert.Greater(t, snapshot.AverageResponseTime, float64(0))
	assert.NotEqual(t, time.Time{}, snapshot.LastUpdated)
}

func TestRateLimitAndAuthFailures(t *testing.T) {
	metrics := GetMetrics()

	// Record rate limit and auth failures
	metrics.RecordRateLimit()
	metrics.RecordRateLimit()
	metrics.RecordAuthFailure()

	snapshot := metrics.GetSnapshot()
	assert.Equal(t, uint64(2), snapshot.RateLimitExceeded)
	assert.Equal(t, uint64(1), snapshot.AuthFailures)
	assert.NotEqual(t, time.Time{}, snapshot.LastUpdated)
}

func TestConcurrentAccess(t *testing.T) {
	metrics := GetMetrics()
	const numGoroutines = 10
	const numOperations = 100

	// Start multiple goroutines to test concurrent access
	done := make(chan bool)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < numOperations; j++ {
				metrics.RecordRequest(time.Duration(j)*time.Millisecond, j%2 == 0)
				metrics.GetSnapshot()
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify final state
	snapshot := metrics.GetSnapshot()
	assert.Equal(t, uint64(numGoroutines*numOperations), snapshot.TotalRequests)
	assert.Equal(t, uint64(numGoroutines*numOperations/2), snapshot.SuccessfulRequests)
}
