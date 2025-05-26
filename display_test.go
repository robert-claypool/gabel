package main

import (
	"strings"
	"testing"
)

func TestHexToRGB(t *testing.T) {
	tests := []struct {
		hex string
		r   uint8
		g   uint8
		b   uint8
	}{
		{"#d73a4a", 215, 58, 74},
		{"d73a4a", 215, 58, 74},
		{"#000000", 0, 0, 0},
		{"#ffffff", 255, 255, 255},
		{"invalid", 128, 128, 128}, // Default gray
		{"", 128, 128, 128},         // Default gray
	}

	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			r, g, b := hexToRGB(tt.hex)
			if r != tt.r || g != tt.g || b != tt.b {
				t.Errorf("hexToRGB(%q) = (%d, %d, %d), want (%d, %d, %d)",
					tt.hex, r, g, b, tt.r, tt.g, tt.b)
			}
		})
	}
}

func TestTruncateDescription(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Short description", "Short description"},
		{"This is exactly sixty characters long when we count them all!", "This is exactly sixty characters long when we count them ..."},
		{"This is a very long description that exceeds sixty characters and should be truncated", "This is a very long description that exceeds sixty charac..."},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := truncateDescription(tt.input)
			if result != tt.expected {
				t.Errorf("truncateDescription(%q) = %q, want %q", tt.input, result, tt.expected)
			}
			if len(tt.input) > 60 && !strings.HasSuffix(result, "...") {
				t.Error("Truncated description should end with ...")
			}
		})
	}
}

func TestFormatLabel(t *testing.T) {
	label := Label{
		Name:        "bug",
		Color:       "d73a4a",
		Description: "Something isn't working",
	}

	// Test without description
	result := FormatLabel(label, false)
	if !strings.Contains(result, "bug") {
		t.Error("Formatted label should contain name")
	}
	if !strings.Contains(result, "#d73a4a") {
		t.Error("Formatted label should contain hex color with #")
	}
	if strings.Contains(result, "Something isn't working") {
		t.Error("Formatted label without verbose should not contain description")
	}

	// Test with description
	result = FormatLabel(label, true)
	if !strings.Contains(result, "Something isn't working") {
		t.Error("Formatted label with verbose should contain description")
	}
}