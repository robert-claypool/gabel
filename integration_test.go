// +build integration

package main

import (
	"testing"
)

// Run with: go test -tags=integration -v ./...
// These tests require real GitHub access

func TestFetchLabelsIntegration(t *testing.T) {
	// Test with a well-known public repo
	labels, err := FetchLabels("golang/go")
	if err != nil {
		t.Skipf("Skipping integration test: %v", err)
	}
	
	if len(labels) == 0 {
		t.Error("Expected to find labels in golang/go repo")
	}
	
	// Check that we got valid label data
	for _, label := range labels {
		if label.Name == "" {
			t.Error("Found label with empty name")
		}
		if label.Color == "" {
			t.Error("Found label with empty color")
		}
	}
}

func TestFetchLabelsErrors(t *testing.T) {
	// Test non-existent repo
	_, err := FetchLabels("this-owner-does-not-exist-12345/repo")
	if err == nil {
		t.Error("Expected error for non-existent repo")
	}
	
	// Test invalid format
	_, err = FetchLabels("not-a-valid-format")
	if err == nil {
		t.Error("Expected error for invalid repo format")
	}
}