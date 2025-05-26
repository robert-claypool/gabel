package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Formats a label for display
func FormatLabel(label Label, showDescription bool) string {
	name := label.Name
	hex := label.Color
	if !strings.HasPrefix(hex, "#") {
		hex = "#" + hex
	}

	// Try to show a colored block, fallback to plain text
	block := getColorBlock(hex)
	result := fmt.Sprintf("%s %s %s", block, name, hex)

	if showDescription && label.Description != "" {
		desc := truncateDescription(label.Description)
		result += fmt.Sprintf("  %s", desc)
	}

	return result
}

// Returns a colored block using fatih/color package
func getColorBlock(hex string) string {
	r, g, b := hexToRGB(hex)
	
	// Let fatih/color handle terminal detection
	c := color.New(color.Attribute(0))
	c = c.Add(color.Attribute(38)) // Set foreground
	c = c.Add(color.Attribute(2))  // RGB mode
	
	// Try 24-bit color, package will fallback if not supported
	return c.Sprintf("\033[38;2;%d;%d;%dm█\033[0m", r, g, b)
}

// Converts hex to RGB
func hexToRGB(hex string) (r, g, b uint8) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 128, 128, 128
	}

	_, _ = fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return
}

// Truncates description to 60 chars
func truncateDescription(desc string) string {
	if len(desc) <= 60 {
		return desc
	}
	return desc[:57] + "..."
}