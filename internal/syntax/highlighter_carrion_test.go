package syntax

import (
	"testing"

	_ "github.com/javanhut/vem/internal/syntax/lexers"
)

func TestCarrionHighlighting(t *testing.T) {
	// Create a highlighter for a .crl file
	h := NewHighlighter("test.crl")
	if h == nil {
		t.Fatal("Failed to create highlighter")
	}

	// Check that the language is detected as Carrion
	lang := h.GetLanguage()
	if lang != "Carrion" {
		t.Errorf("Expected language 'Carrion', got '%s'", lang)
	}

	// Test highlighting some Carrion code
	testCode := `spell greet(name):
    message = f"Hello, {name}!"
    print(message)
    return True`

	tokens := h.HighlightLine(0, testCode)
	if len(tokens) == 0 {
		t.Error("No tokens generated for Carrion code")
	}

	// Verify highlighting is enabled
	if !h.IsEnabled() {
		t.Error("Syntax highlighting should be enabled")
	}
}

func TestCarrionKeywordHighlighting(t *testing.T) {
	h := NewHighlighter("program.crl")

	testCases := []struct {
		name string
		code string
	}{
		{"Function definition", "spell calculate(x, y):"},
		{"Class definition", "grim Person:"},
		{"Error handling", "attempt:"},
		{"Concurrency", "diverge worker:"},
		{"Main entry", "main:"},
		{"Convergence", "converge worker1, worker2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokens := h.HighlightLine(0, tc.code)
			if len(tokens) == 0 {
				t.Errorf("No tokens generated for: %s", tc.code)
			}
		})
	}
}

func TestCarrionShouldHighlight(t *testing.T) {
	testCases := []struct {
		path     string
		expected bool
	}{
		{"test.crl", true},
		{"program.crl", true},
		{"script.py", true},
		{"main.go", true},
		{"binary.exe", false},
		{"image.png", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := ShouldHighlight(tc.path)
			if result != tc.expected {
				t.Errorf("ShouldHighlight(%s) = %v, expected %v", tc.path, result, tc.expected)
			}
		})
	}
}
