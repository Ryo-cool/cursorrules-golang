package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMetricsCollector(t *testing.T) {
	collector := NewMetricsCollector()
	assert.NotNil(t, collector)
}

func TestRecordRequestDuration(t *testing.T) {
	collector := NewMetricsCollector()

	testCases := []struct {
		name     string
		path     string
		duration time.Duration
	}{
		{
			name:     "normal request",
			path:     "/api/v1/users",
			duration: 100 * time.Millisecond,
		},
		{
			name:     "slow request",
			path:     "/api/v1/posts",
			duration: 1 * time.Second,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			collector.RecordRequestDuration(tc.path, tc.duration)
			// Verify metrics were recorded correctly
			stats := collector.GetStats()
			assert.Contains(t, stats, tc.path)
		})
	}
}

func TestGetStats(t *testing.T) {
	collector := NewMetricsCollector()

	// Record some test data
	collector.RecordRequestDuration("/test", 100*time.Millisecond)
	collector.RecordRequestDuration("/test", 200*time.Millisecond)

	stats := collector.GetStats()
	assert.NotEmpty(t, stats)
	assert.Contains(t, stats, "/test")
}

func TestMetricsReset(t *testing.T) {
	collector := NewMetricsCollector()

	// Add some metrics
	collector.RecordRequestDuration("/test", 100*time.Millisecond)

	// Reset metrics
	collector.Reset()

	// Verify metrics are cleared
	stats := collector.GetStats()
	assert.Empty(t, stats)
}
