package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestCheckGitHubCLI(t *testing.T) {
	// This test requires gh to be installed
	err := CheckGitHubCLI()
	
	// Check if gh is in PATH
	_, ghErr := exec.LookPath("gh")
	
	if ghErr != nil && err == nil {
		t.Error("Expected error when gh is not found, but got nil")
	}
	
	if ghErr == nil && err != nil {
		// gh exists but auth might fail - that's expected in test environment
		expectedMsg := "Not authenticated with GitHub"
		if err.Error() != expectedMsg+"\nRun: gh auth login" {
			t.Errorf("Unexpected error: %v", err)
		}
	}
}


func TestIsValidRepoFormat(t *testing.T) {
	// Additional edge cases
	tests := []struct {
		repo  string
		valid bool
	}{
		{"", false},
		{"owner", false},
		{"owner/", false},
		{"/repo", false},
		{"owner/repo/extra", false},
		{"owner/repo", true},
		{"my-org/my-repo", true},
		{"123/456", true},
		{"owner_name/repo_name", true},
	}
	
	for _, tt := range tests {
		result := isValidRepo(tt.repo)
		if result != tt.valid {
			t.Errorf("isValidRepo(%q) = %v, want %v", tt.repo, result, tt.valid)
		}
	}
}

// Test error message formats
func TestErrorMessages(t *testing.T) {
	// Save original PATH
	oldPath := os.Getenv("PATH")
	defer func() { _ = os.Setenv("PATH", oldPath) }()
	
	// Test with empty PATH (gh not found)
	_ = os.Setenv("PATH", "")
	err := CheckGitHubCLI()
	
	if err == nil {
		t.Error("Expected error when gh is not in PATH")
	} else if err.Error() != "GitHub CLI (gh) is required but not found.\nInstall it from: https://cli.github.com" {
		t.Errorf("Unexpected error message: %v", err)
	}
}