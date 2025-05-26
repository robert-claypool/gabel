package main

import (
	"strings"
	"testing"
)

func TestValidateColor(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		wantErr  bool
	}{
		{"#d73a4a", "d73a4a", false},
		{"d73a4a", "d73a4a", false},
		{"#FFFFFF", "FFFFFF", false},
		{"ffffff", "ffffff", false},
		{"#fff", "", true},          // Too short
		{"#gggggg", "", true},       // Invalid hex
		{"12345", "", true},         // Too short
		{"#1234567", "", true},      // Too long
		{"", "", true},              // Empty
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := validateColor(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if result != tt.expected {
				t.Errorf("validateColor(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestValidateLabel(t *testing.T) {
	tests := []struct {
		name    string
		label   Label
		wantErr bool
		errMsg  string
	}{
		{
			"valid label",
			Label{Name: "bug", Color: "#d73a4a", Description: "Something isn't working"},
			false,
			"",
		},
		{
			"empty name",
			Label{Name: "", Color: "#d73a4a", Description: "desc"},
			true,
			"empty",
		},
		{
			"whitespace name",
			Label{Name: "   ", Color: "#d73a4a", Description: "desc"},
			true,
			"empty",
		},
		{
			"invalid color",
			Label{Name: "bug", Color: "invalid", Description: "desc"},
			true,
			"color",
		},
		{
			"description too long",
			Label{Name: "bug", Color: "#d73a4a", Description: strings.Repeat("a", 101)},
			true,
			"too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLabel(tt.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateLabel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("validateLabel() error = %v, should contain %q", err, tt.errMsg)
			}
		})
	}
}