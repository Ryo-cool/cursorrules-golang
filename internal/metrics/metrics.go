package metrics

import (
	"sync"
	"time"
)

// Metrics represents application metrics
type Metrics struct {
	mu sync.RWMutex

	// Request metrics
	TotalRequests      uint64
	SuccessfulRequests uint64
	FailedRequests     uint64

	// Response time metrics (in milliseconds)
	AverageResponseTime float64
	MinResponseTime     float64
	MaxResponseTime     float64

	// Rate limiting metrics
	RateLimitExceeded uint64

	// Authentication metrics
	AuthFailures uint64

	// Last update timestamp
	LastUpdated time.Time
}

// MetricsSnapshot represents a snapshot of metrics without mutex
type MetricsSnapshot struct {
	TotalRequests       uint64
	SuccessfulRequests  uint64
	FailedRequests      uint64
	AverageResponseTime float64
	MinResponseTime     float64
	MaxResponseTime     float64
	RateLimitExceeded   uint64
	AuthFailures        uint64
	LastUpdated         time.Time
}

var (
	defaultMetrics *Metrics
	once           sync.Once
)

// GetMetrics returns the singleton metrics instance
func GetMetrics() *Metrics {
	once.Do(func() {
		defaultMetrics = &Metrics{
			MinResponseTime: float64(^uint64(0) >> 1), // Initialize with max value
			LastUpdated:     time.Now(),
		}
	})
	return defaultMetrics
}

// RecordRequest records request metrics
func (m *Metrics) RecordRequest(duration time.Duration, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.TotalRequests++
	if success {
		m.SuccessfulRequests++
	} else {
		m.FailedRequests++
	}

	// Update response time metrics
	durationMs := float64(duration.Milliseconds())
	m.AverageResponseTime = (m.AverageResponseTime*float64(m.TotalRequests-1) + durationMs) / float64(m.TotalRequests)

	if durationMs < m.MinResponseTime {
		m.MinResponseTime = durationMs
	}
	if durationMs > m.MaxResponseTime {
		m.MaxResponseTime = durationMs
	}

	m.LastUpdated = time.Now()
}

// RecordRateLimit records a rate limit exceeded event
func (m *Metrics) RecordRateLimit() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RateLimitExceeded++
	m.LastUpdated = time.Now()
}

// RecordAuthFailure records an authentication failure
func (m *Metrics) RecordAuthFailure() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.AuthFailures++
	m.LastUpdated = time.Now()
}

// GetSnapshot returns a copy of the current metrics without the mutex
func (m *Metrics) GetSnapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return MetricsSnapshot{
		TotalRequests:       m.TotalRequests,
		SuccessfulRequests:  m.SuccessfulRequests,
		FailedRequests:      m.FailedRequests,
		AverageResponseTime: m.AverageResponseTime,
		MinResponseTime:     m.MinResponseTime,
		MaxResponseTime:     m.MaxResponseTime,
		RateLimitExceeded:   m.RateLimitExceeded,
		AuthFailures:        m.AuthFailures,
		LastUpdated:         m.LastUpdated,
	}
}
