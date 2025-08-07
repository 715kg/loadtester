package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/715kg/loadtester/internal/loadtest"
)

func TestServeHTML(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serveHTML)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Load Tester") {
		t.Error("Response should contain 'Load Tester'")
	}

	if !strings.Contains(body, "<!DOCTYPE html>") {
		t.Error("Response should contain HTML doctype")
	}
}

func TestHandleStartInvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/start", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleStart)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestHandleStartInvalidJSON(t *testing.T) {
	invalidJSON := `{"invalid": json}`
	req, err := http.NewRequest("POST", "/start", strings.NewReader(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleStart)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response["success"] != false {
		t.Error("Expected success to be false for invalid JSON")
	}
}

func TestHandleStartEmptyURL(t *testing.T) {
	config := loadtest.Config{
		TargetURL:      "",
		Method:         "GET",
		TotalRequests:  100,
		MaxConcurrency: 10,
		Timeout:        5,
		MaxMemory:      1.0,
	}

	jsonData, _ := json.Marshal(config)
	req, err := http.NewRequest("POST", "/start", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleStart)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response["success"] != false {
		t.Error("Expected success to be false for empty URL")
	}
}

func TestHandleStop(t *testing.T) {
	req, err := http.NewRequest("POST", "/stop", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleStop)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response["success"] != true {
		t.Error("Expected success to be true for stop request")
	}
}

func TestHandleStopInvalidMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/stop", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleStop)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func TestHandleStats(t *testing.T) {
	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleStats)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type to be application/json, got %s", contentType)
	}

	var stats loadtest.Stats
	err = json.Unmarshal(rr.Body.Bytes(), &stats)
	if err != nil {
		t.Fatal(err)
	}

	// Basic validation of stats structure
	if stats.Total < 0 {
		t.Error("Stats Total should not be negative")
	}

	if stats.Progress < 0 || stats.Progress > 1 {
		t.Error("Stats Progress should be between 0 and 1")
	}
}

func TestStopCurrentTest(t *testing.T) {
	// This function should not panic when called
	stopCurrentTest()

	// Wait a bit to ensure any goroutines finish
	time.Sleep(100 * time.Millisecond)
}

// Benchmark tests
func BenchmarkServeHTML(b *testing.B) {
	req, _ := http.NewRequest("GET", "/", nil)

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(serveHTML)
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkHandleStats(b *testing.B) {
	req, _ := http.NewRequest("GET", "/stats", nil)

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handleStats)
		handler.ServeHTTP(rr, req)
	}
}
