package main

import (
	"testing"
)

func TestBuildPickerItems(t *testing.T) {
	sourceLabels := []Label{
		{Name: "bug", Color: "#d73a4a", Description: "Something isn't working"},
		{Name: "feature", Color: "#a2eeef", Description: "New feature"},
		{Name: "duplicate", Color: "#cfd3d7", Description: "This issue already exists"},
	}

	destLabels := []Label{
		{Name: "bug", Color: "#ff0000", Description: "A bug"},
		{Name: "wontfix", Color: "#ffffff", Description: "This will not be worked on"},
	}

	items := buildPickerItems(sourceLabels, destLabels)

	// Should have 4 items total: 2 dest-only + 2 source-only
	if len(items) != 4 {
		t.Errorf("Expected 4 items, got %d", len(items))
	}

	// First items should be dest-only
	if !items[0].IsDestOnly || items[0].Label.Name != "bug" {
		t.Error("First item should be dest-only 'bug'")
	}
	if !items[1].IsDestOnly || items[1].Label.Name != "wontfix" {
		t.Error("Second item should be dest-only 'wontfix'")
	}

	// Remaining items should be source labels (excluding duplicate 'bug')
	if items[2].IsDestOnly || items[2].Label.Name != "feature" {
		t.Error("Third item should be source 'feature'")
	}
	if items[3].IsDestOnly || items[3].Label.Name != "duplicate" {
		t.Error("Fourth item should be source 'duplicate'")
	}

	// All should be selected by default
	for i, item := range items {
		if !item.Selected {
			t.Errorf("Item %d should be selected by default", i)
		}
	}
}

func TestCalculateActions(t *testing.T) {
	selectedLabels := []Label{
		{Name: "bug", Color: "#d73a4a"},
		{Name: "feature", Color: "#a2eeef"},
	}

	destLabels := []Label{
		{Name: "bug", Color: "#ff0000"},
		{Name: "wontfix", Color: "#ffffff"},
		{Name: "documentation", Color: "#0075ca"},
	}

	summary := calculateActions(selectedLabels, destLabels)

	// Should create 1 label (feature)
	if len(summary.ToCreate) != 1 || summary.ToCreate[0].Name != "feature" {
		t.Errorf("Expected to create 'feature', got %v", summary.ToCreate)
	}

	// Should delete 2 labels (wontfix, documentation)
	if len(summary.ToDelete) != 2 {
		t.Errorf("Expected to delete 2 labels, got %d", len(summary.ToDelete))
	}

	// Should keep 1 label (bug)
	if len(summary.ToKeep) != 1 || summary.ToKeep[0].Name != "bug" {
		t.Errorf("Expected to keep 'bug', got %v", summary.ToKeep)
	}
}

func TestCalculateActionsCaseInsensitive(t *testing.T) {
	selectedLabels := []Label{
		{Name: "Bug", Color: "#d73a4a"},
		{Name: "FEATURE", Color: "#a2eeef"},
	}

	destLabels := []Label{
		{Name: "bug", Color: "#ff0000"},
		{Name: "feature", Color: "#00ff00"},
	}

	summary := calculateActions(selectedLabels, destLabels)

	// Should recognize case-insensitive matches
	if len(summary.ToCreate) != 0 {
		t.Errorf("Expected no labels to create, got %d", len(summary.ToCreate))
	}

	if len(summary.ToDelete) != 0 {
		t.Errorf("Expected no labels to delete, got %d", len(summary.ToDelete))
	}

	if len(summary.ToKeep) != 2 {
		t.Errorf("Expected to keep 2 labels, got %d", len(summary.ToKeep))
	}
}