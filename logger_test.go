package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLogDebug(t *testing.T) {
	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Test with debug mode off
	InitLogger(false)
	LogDebug("This should not appear")
	
	// Test with debug mode on
	InitLogger(true)
	LogDebug("Test message %s", "with args")

	_ = w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if strings.Contains(output, "This should not appear") {
		t.Error("Debug message appeared when debug mode was off")
	}

	if !strings.Contains(output, "[DEBUG] Test message with args") {
		t.Error("Debug message did not appear when debug mode was on")
	}
}

func TestLogError(t *testing.T) {
	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	LogError("Error message %d", 42)

	_ = w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "[ERROR] Error message 42") {
		t.Error("Error message did not appear correctly")
	}
}