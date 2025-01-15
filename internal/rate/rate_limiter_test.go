package rate

import (
	"testing"
	"time"
)

func TestNewLimiter(t *testing.T) {
	cfg := LimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         3,
		QueueSize:         5,
		Enabled:           true,
	}

	limiter := New(cfg)

	if limiter == nil {
		t.Fatal("Expected non-nil limiter")
	}

	if cap(limiter.requests) != cfg.QueueSize {
		t.Errorf("Expected queue size %d, got %d", cfg.QueueSize, cap(limiter.requests))
	}

	if cap(limiter.burstyLimiter) != cfg.BurstSize {
		t.Errorf("Expected burst size %d, got %d", cfg.BurstSize, cap(limiter.burstyLimiter))
	}
}

func TestLimiterAllow(t *testing.T) {
	tests := []struct {
		name          string
		config        LimiterConfig
		requests      int
		expectedAllow bool
		sleepDuration time.Duration
	}{
		{
			name: "Allow burst requests",
			config: LimiterConfig{
				RequestsPerSecond: 10,
				BurstSize:         3,
				QueueSize:         5,
				Enabled:           true,
			},
			requests:      3,
			expectedAllow: true,
		},
		{
			name: "Exceed burst limit",
			config: LimiterConfig{
				RequestsPerSecond: 10,
				BurstSize:         2,
				QueueSize:         5,
				Enabled:           true,
			},
			requests:      3,
			expectedAllow: false,
		},
		{
			name: "Test replenishment",
			config: LimiterConfig{
				RequestsPerSecond: 10,
				BurstSize:         1,
				QueueSize:         5,
				Enabled:           true,
			},
			requests:      2,
			sleepDuration: time.Millisecond * 200, // Wait for token replenishment
			expectedAllow: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := New(tt.config)

			// Consume initial tokens
			for i := 0; i < tt.requests-1; i++ {
				limiter.Allow()
			}

			if tt.sleepDuration > 0 {
				time.Sleep(tt.sleepDuration)
			}

			// Test final request
			result := limiter.Allow()
			if result != tt.expectedAllow {
				t.Errorf("Expected Allow() to return %v, got %v", tt.expectedAllow, result)
			}
		})
	}
}

func TestLimiterConcurrency(t *testing.T) {
	cfg := LimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         5,
		QueueSize:         10,
		Enabled:           true,
	}

	limiter := New(cfg)
	concurrentRequests := 10
	results := make(chan bool, concurrentRequests)

	// Send concurrent requests
	for i := 0; i < concurrentRequests; i++ {
		go func() {
			results <- limiter.Allow()
		}()
	}

	// Collect results
	allowed := 0
	for i := 0; i < concurrentRequests; i++ {
		if <-results {
			allowed++
		}
	}

	// We expect only BurstSize number of requests to be allowed
	if allowed > cfg.BurstSize {
		t.Errorf("Expected maximum %d requests to be allowed, got %d", cfg.BurstSize, allowed)
	}
}

func TestLimiterReplenishment(t *testing.T) {
	cfg := LimiterConfig{
		RequestsPerSecond: 10,
		BurstSize:         2,
		QueueSize:         5,
		Enabled:           true,
	}

	limiter := New(cfg)

	// Consume all initial tokens
	for i := 0; i < cfg.BurstSize; i++ {
		if !limiter.Allow() {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}

	// Next request should be denied
	if limiter.Allow() {
		t.Error("Expected request to be denied after consuming all tokens")
	}

	// Wait for token replenishment
	time.Sleep(time.Second / time.Duration(cfg.RequestsPerSecond))

	// Should be allowed after replenishment
	if !limiter.Allow() {
		t.Error("Expected request to be allowed after token replenishment")
	}
}
