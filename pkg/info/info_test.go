package info

import (
	"strings"
	"testing"
)

func TestProgramConstants(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"ProgramName", ProgramName, "Load Tester"},
		{"Author", Author, "Alexey Ulyanov"},
		{"License", License, "MIT"},
		{"Category", Category, "Development Tools"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.value, tt.expected)
			}
		})
	}
}

func TestWebsiteURL(t *testing.T) {
	if !strings.HasPrefix(Website, "https://github.com/715kg") {
		t.Errorf("Website should start with https://github.com/715kg, got %q", Website)
	}
}

func TestProgramDescription(t *testing.T) {
	if len(ProgramDescription) == 0 {
		t.Error("ProgramDescription should not be empty")
	}

	if !strings.Contains(ProgramDescription, "HTTP") {
		t.Error("ProgramDescription should mention HTTP")
	}
}

func TestPurpose(t *testing.T) {
	if len(Purpose) == 0 {
		t.Error("Purpose should not be empty")
	}

	expectedKeywords := []string{"testing", "performance"}
	for _, keyword := range expectedKeywords {
		if !strings.Contains(strings.ToLower(Purpose), keyword) {
			t.Errorf("Purpose should contain keyword %q", keyword)
		}
	}
}

func TestProgramSignature(t *testing.T) {
	if !strings.Contains(ProgramSignature, "Load Tester") {
		t.Error("ProgramSignature should contain program name")
	}

	if !strings.Contains(ProgramSignature, "BEGIN PROGRAM INFO") {
		t.Error("ProgramSignature should have proper format")
	}

	if !strings.Contains(ProgramSignature, "END PROGRAM INFO") {
		t.Error("ProgramSignature should have proper format")
	}
}
