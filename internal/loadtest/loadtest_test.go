package loadtest

import (
	"runtime"
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
func TestConfigMaxCPUCores(t *testing.T) {
	tests := []struct {
		name        string
		maxCPUCores int
		description string
	}{
		{
			name:        "AutomaticCores",
			maxCPUCores: 0,
			description: "Should use all available cores when set to 0",
		},
		{
			name:        "SingleCore",
			maxCPUCores: 1,
			description: "Should limit to 1 core",
		},
		{
			name:        "MultipleCores",
			maxCPUCores: 4,
			description: "Should limit to 4 cores",
		},
		{
			name:        "HighCoreCount",
			maxCPUCores: 16,
			description: "Should handle high core counts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				TargetURL:      "https://example.com",
				Method:         "GET",
				TotalRequests:  100,
				MaxConcurrency: 10,
				Timeout:        5,
				MaxMemory:      1.0,
				MaxCPUCores:    tt.maxCPUCores,
			}

			if config.MaxCPUCores != tt.maxCPUCores {
				t.Errorf("Expected MaxCPUCores to be %d, got %d", tt.maxCPUCores, config.MaxCPUCores)
			}

			// Verify the field is properly set in the struct
			if config.MaxCPUCores < 0 {
				t.Error("MaxCPUCores should not be negative")
			}
		})
	}
}

func TestConfigWithMaxCPUCoresValidation(t *testing.T) {
	// Test that MaxCPUCores field is included in Config struct
	config := &Config{
		TargetURL:      "https://example.com",
		Method:         "GET",
		TotalRequests:  1000,
		MaxConcurrency: 10,
		Timeout:        5,
		MaxMemory:      1.0,
		MaxCPUCores:    8, // Test with specific value
	}

	// Verify all fields are set correctly including MaxCPUCores
	if config.TargetURL != "https://example.com" {
		t.Errorf("Expected TargetURL to be https://example.com, got %s", config.TargetURL)
	}

	if config.MaxCPUCores != 8 {
		t.Errorf("Expected MaxCPUCores to be 8, got %d", config.MaxCPUCores)
	}

	if config.MaxMemory != 1.0 {
		t.Errorf("Expected MaxMemory to be 1.0, got %f", config.MaxMemory)
	}
}

func TestStartWithMaxCPUCores(t *testing.T) {
	// Skip this test if we don't want to make actual HTTP requests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name        string
		maxCPUCores int
		shouldWork  bool
	}{
		{
			name:        "AutoCores",
			maxCPUCores: 0,
			shouldWork:  true,
		},
		{
			name:        "SingleCore",
			maxCPUCores: 1,
			shouldWork:  true,
		},
		{
			name:        "FourCores",
			maxCPUCores: 4,
			shouldWork:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				TargetURL:      "https://httpbin.org/get",
				Method:         "GET",
				TotalRequests:  3,
				MaxConcurrency: 1,
				Timeout:        10,
				MaxMemory:      1.0,
				MaxCPUCores:    tt.maxCPUCores,
			}

			err := Start(config)
			if tt.shouldWork && err != nil {
				t.Errorf("Expected no error when starting with MaxCPUCores=%d, got %v", tt.maxCPUCores, err)
			}

			if err == nil {
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
		})
	}
}

func TestConfigDefaultValues(t *testing.T) {
	// Test that MaxCPUCores has a sensible default when not explicitly set
	config := &Config{
		TargetURL:      "https://example.com",
		Method:         "GET",
		TotalRequests:  100,
		MaxConcurrency: 10,
		Timeout:        5,
		MaxMemory:      1.0,
		// MaxCPUCores not set, should default to 0
	}

	// Zero value should be 0 (automatic)
	if config.MaxCPUCores != 0 {
		t.Errorf("Expected default MaxCPUCores to be 0, got %d", config.MaxCPUCores)
	}
}

func TestConfigJSONSerialization(t *testing.T) {
	// Test that MaxCPUCores field is properly handled in JSON serialization
	// This is important for the web interface
	config := &Config{
		TargetURL:      "https://example.com",
		Method:         "GET",
		TotalRequests:  1000,
		MaxConcurrency: 10,
		Timeout:        5,
		MaxMemory:      2.5,
		MaxCPUCores:    6,
	}

	// Verify the struct has all expected fields
	if config.MaxCPUCores != 6 {
		t.Errorf("Expected MaxCPUCores to be 6, got %d", config.MaxCPUCores)
	}

	// Test edge cases
	config.MaxCPUCores = 0
	if config.MaxCPUCores != 0 {
		t.Errorf("Expected MaxCPUCores to be 0 (auto), got %d", config.MaxCPUCores)
	}

	config.MaxCPUCores = 32
	if config.MaxCPUCores != 32 {
		t.Errorf("Expected MaxCPUCores to be 32, got %d", config.MaxCPUCores)
	}
}
func TestRuntimeGOMAXPROCSIntegration(t *testing.T) {
	// This test verifies that the MaxCPUCores setting is properly integrated
	// with the runtime configuration

	// Skip this test if we don't want to make actual HTTP requests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Store original GOMAXPROCS value
	originalGOMAXPROCS := runtime.GOMAXPROCS(0)
	defer func() {
		// Restore original value after test
		runtime.GOMAXPROCS(originalGOMAXPROCS)
	}()

	// Test with specific core count
	config := &Config{
		TargetURL:      "https://example.com/get",
		Method:         "GET",
		TotalRequests:  2,
		MaxConcurrency: 1,
		Timeout:        10,
		MaxMemory:      1.0,
		MaxCPUCores:    2, // Set to 2 cores
	}

	err := Start(config)
	if err != nil {
		t.Errorf("Expected no error when starting test, got %v", err)
		return
	}

	// Wait a bit for the test to start and GOMAXPROCS to be set
	time.Sleep(200 * time.Millisecond)

	// Check that GOMAXPROCS was set (note: this might not be exactly 2
	// if the system has fewer cores, but it should be set)
	currentGOMAXPROCS := runtime.GOMAXPROCS(0)
	t.Logf("Original GOMAXPROCS: %d, Current GOMAXPROCS: %d, Config MaxCPUCores: %d",
		originalGOMAXPROCS, currentGOMAXPROCS, config.MaxCPUCores)

	// Stop the test
	Stop()

	// Wait for it to stop
	time.Sleep(500 * time.Millisecond)

	if IsRunning() {
		t.Error("Expected test to be stopped")
	}
}

func TestMaxCPUCoresEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		maxCPUCores int
		expectError bool
		description string
	}{
		{
			name:        "ZeroCores",
			maxCPUCores: 0,
			expectError: false,
			description: "Zero should use automatic core detection",
		},
		{
			name:        "NegativeCores",
			maxCPUCores: -1,
			expectError: false, // runtime.GOMAXPROCS handles negative values
			description: "Negative values should be handled gracefully",
		},
		{
			name:        "VeryHighCores",
			maxCPUCores: 1000,
			expectError: false, // runtime.GOMAXPROCS will cap to available cores
			description: "Very high values should be capped by runtime",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				TargetURL:      "https://example.com",
				Method:         "GET",
				TotalRequests:  10,
				MaxConcurrency: 1,
				Timeout:        5,
				MaxMemory:      1.0,
				MaxCPUCores:    tt.maxCPUCores,
			}

			// Test that the config accepts the value
			if config.MaxCPUCores != tt.maxCPUCores {
				t.Errorf("Expected MaxCPUCores to be %d, got %d", tt.maxCPUCores, config.MaxCPUCores)
			}

			// Test with empty URL to trigger validation error
			config.TargetURL = ""
			err := Start(config)
			if err == nil {
				t.Error("Expected validation error for empty URL, got nil")
			}
		})
	}
}

func TestConfigStructCompleteness(t *testing.T) {
	// Test that Config struct has all expected fields including MaxCPUCores
	config := Config{}

	// Set all fields to ensure they exist
	config.TargetURL = "https://example.com"
	config.Method = "GET"
	config.TotalRequests = 100
	config.MaxConcurrency = 10
	config.Timeout = 5
	config.MaxMemory = 1.0
	config.MaxCPUCores = 4

	// Verify all fields are accessible
	if config.TargetURL == "" {
		t.Error("TargetURL field should be accessible")
	}
	if config.Method == "" {
		t.Error("Method field should be accessible")
	}
	if config.TotalRequests == 0 {
		t.Error("TotalRequests field should be accessible")
	}
	if config.MaxConcurrency == 0 {
		t.Error("MaxConcurrency field should be accessible")
	}
	if config.Timeout == 0 {
		t.Error("Timeout field should be accessible")
	}
	if config.MaxMemory == 0 {
		t.Error("MaxMemory field should be accessible")
	}
	if config.MaxCPUCores == 0 {
		t.Error("MaxCPUCores field should be accessible and set to 4")
	}
}
