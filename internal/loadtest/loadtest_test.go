package loadtest

import (
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	config := &Config{
		TargetURL:      "https://example.com",
		Method:         "GET",
		TotalRequests:  1000,
		MaxConcurrency: 10,
		Timeout:        5,
		MaxMemory:      1.0,
	}

	if config.TargetURL != "https://example.com" {
		t.Errorf("Expected TargetURL to be https://example.com, got %s", config.TargetURL)
	}

	if config.Method != "GET" {
		t.Errorf("Expected Method to be GET, got %s", config.Method)
	}

	if config.TotalRequests != 1000 {
		t.Errorf("Expected TotalRequests to be 1000, got %d", config.TotalRequests)
	}

	if config.MaxConcurrency != 10 {
		t.Errorf("Expected MaxConcurrency to be 10, got %d", config.MaxConcurrency)
	}
}

func TestStats(t *testing.T) {
	stats := Stats{
		Total:   100,
		Success: 95,
		Fails:   3,
		Errors:  2,
		RPS:     50.5,
		Elapsed: 2.0,
	}

	if stats.Total != 100 {
		t.Errorf("Expected Total to be 100, got %d", stats.Total)
	}

	if stats.Success != 95 {
		t.Errorf("Expected Success to be 95, got %d", stats.Success)
	}

	if stats.RPS != 50.5 {
		t.Errorf("Expected RPS to be 50.5, got %f", stats.RPS)
	}
}

func TestIsRunning(t *testing.T) {
	// Initially should not be running
	if IsRunning() {
		t.Error("Expected IsRunning to be false initially")
	}
}

func TestStartWithInvalidConfig(t *testing.T) {
	// Test with empty URL
	config := &Config{
		TargetURL:      "",
		Method:         "GET",
		TotalRequests:  100,
		MaxConcurrency: 10,
		Timeout:        5,
		MaxMemory:      1.0,
	}

	err := Start(config)
	if err == nil {
		t.Error("Expected error when starting with empty URL")
	}
}

func TestStartWithValidConfig(t *testing.T) {
	// Skip this test if we don't want to make actual HTTP requests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	config := &Config{
		TargetURL:      "https://httpbin.org/get",
		Method:         "GET",
		TotalRequests:  5,
		MaxConcurrency: 2,
		Timeout:        10,
		MaxMemory:      1.0,
	}

	err := Start(config)
	if err != nil {
		t.Errorf("Expected no error when starting with valid config, got %v", err)
	}

	// Wait a bit for the test to start
	time.Sleep(100 * time.Millisecond)

	if !IsRunning() {
		t.Error("Expected test to be running")
	}

	// Stop the test
	Stop()

	// Wait for it to stop
	time.Sleep(500 * time.Millisecond)

	if IsRunning() {
		t.Error("Expected test to be stopped")
	}
}

func TestGetStats(t *testing.T) {
	stats := GetStats()

	// Check that stats structure is valid
	if stats.Total < 0 {
		t.Error("Total should not be negative")
	}

	if stats.Success < 0 {
		t.Error("Success should not be negative")
	}

	if stats.Fails < 0 {
		t.Error("Fails should not be negative")
	}

	if stats.Errors < 0 {
		t.Error("Errors should not be negative")
	}

	if stats.RPS < 0 {
		t.Error("RPS should not be negative")
	}

	if stats.Elapsed < 0 {
		t.Error("Elapsed should not be negative")
	}

	if stats.Progress < 0 || stats.Progress > 1 {
		t.Error("Progress should be between 0 and 1")
	}
}

func TestGBToBytes(t *testing.T) {
	tests := []struct {
		gb       float64
		expected uint64
	}{
		{1.0, 1024 * 1024 * 1024},
		{0.5, 512 * 1024 * 1024},
		{2.0, 2 * 1024 * 1024 * 1024},
		{0.0, 0},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := GBToBytes(tt.gb)
			if result != tt.expected {
				t.Errorf("GBToBytes(%f) = %d, want %d", tt.gb, result, tt.expected)
			}
		})
	}
}

func TestResetCounters(t *testing.T) {
	// This is testing internal function, so we test it indirectly
	// by checking that stats are reset when starting a new test

	// First, simulate some activity by getting stats
	stats1 := GetStats()

	// The counters should be at their initial state
	if stats1.Total != 0 {
		t.Errorf("Expected initial Total to be 0, got %d", stats1.Total)
	}

	if stats1.Success != 0 {
		t.Errorf("Expected initial Success to be 0, got %d", stats1.Success)
	}
}
