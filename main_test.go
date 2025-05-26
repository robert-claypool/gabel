package main

import (
	"testing"
)

func TestIsValidRepo(t *testing.T) {
	tests := []struct {
		name     string
		repo     string
		expected bool
	}{
		{"valid repo", "owner/repo", true},
		{"valid repo with dash", "owner-name/repo-name", true},
		{"valid repo with numbers", "owner123/repo456", true},
		{"missing owner", "/repo", false},
		{"missing repo", "owner/", false},
		{"no slash", "ownerrepo", false},
		{"too many slashes", "owner/repo/extra", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidRepo(tt.repo)
			if result != tt.expected {
				t.Errorf("isValidRepo(%q) = %v, want %v", tt.repo, result, tt.expected)
			}
		})
	}
}