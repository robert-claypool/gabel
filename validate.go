package main

import (
	"fmt"
	"regexp"
	"strings"
)

var colorRegex = regexp.MustCompile(`^#?[0-9a-fA-F]{6}$`)

// Validates and normalizes a hex color
func validateColor(color string) (string, error) {
	if !colorRegex.MatchString(color) {
		return "", fmt.Errorf("invalid color format: %s (expected hex like #d73a4a)", color)
	}
	
	// Normalize to 6 chars without #
	color = strings.TrimPrefix(color, "#")
	return color, nil
}

// Validates a label before creation
func validateLabel(label Label) error {
	if strings.TrimSpace(label.Name) == "" {
		return fmt.Errorf("label name cannot be empty")
	}
	
	if _, err := validateColor(label.Color); err != nil {
		return err
	}
	
	// GitHub limits description to 100 chars
	if len(label.Description) > 100 {
		return fmt.Errorf("label description too long (max 100 chars)")
	}
	
	return nil
}