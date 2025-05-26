package main

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// Shows interactive picker and returns selected labels
func ShowPicker(sourceLabels, destLabels []Label, destRepo string, verbose bool) ([]Label, error) {
	// Build picker items
	items := buildPickerItems(sourceLabels, destLabels)
	
	// TODO: Implement interactive selection with promptui
	// For now, return all source labels as placeholder
	
	var selected []Label
	for _, item := range items {
		if item.Selected {
			selected = append(selected, item.Label)
		}
	}
	
	return selected, nil
}

// Builds unified list of picker items
func buildPickerItems(sourceLabels, destLabels []Label) []PickerItem {
	items := []PickerItem{}
	destMap := make(map[string]Label)
	
	// Index destination labels by lowercase name
	for _, label := range destLabels {
		destMap[strings.ToLower(label.Name)] = label
	}
	
	// Add destination-only labels first
	for _, label := range destLabels {
		items = append(items, PickerItem{
			Label:      label,
			Selected:   true,
			IsDestOnly: true,
		})
	}
	
	// Add source labels
	for _, label := range sourceLabels {
		// Skip if already exists in destination
		if _, exists := destMap[strings.ToLower(label.Name)]; exists {
			continue
		}
		
		items = append(items, PickerItem{
			Label:      label,
			Selected:   true,
			IsDestOnly: false,
		})
	}
	
	return items
}

// Shows final confirmation and applies changes
func ConfirmAndApply(selectedLabels, destLabels []Label, destRepo string) error {
	summary := calculateActions(selectedLabels, destLabels)
	
	// Show final state
	fmt.Printf("\nFinal state for %s:\n", destRepo)
	for _, label := range selectedLabels {
		fmt.Printf("  ✓ %s\n", FormatLabel(label, false))
	}
	
	// Show actions
	fmt.Printf("\nActions:\n")
	if len(summary.ToCreate) > 0 {
		fmt.Printf("  • Create %d labels\n", len(summary.ToCreate))
	}
	if len(summary.ToDelete) > 0 {
		fmt.Printf("  • Delete %d labels", len(summary.ToDelete))
		if len(summary.ToDelete) == 1 {
			fmt.Printf(" (%s)", summary.ToDelete[0].Name)
		}
		fmt.Println()
	}
	if len(summary.ToKeep) > 0 {
		fmt.Printf("  • Keep %d existing labels\n", len(summary.ToKeep))
	}
	
	// Confirm
	prompt := promptui.Prompt{
		Label:     "Proceed",
		IsConfirm: true,
		Default:   "n",
	}
	
	_, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("cancelled")
	}
	
	// Apply changes
	return applyChanges(summary, destRepo)
}

// Calculates what actions need to be taken
func calculateActions(selectedLabels, destLabels []Label) ActionSummary {
	// TODO: Implement action calculation
	return ActionSummary{}
}

// Applies the changes to the destination repository
func applyChanges(summary ActionSummary, destRepo string) error {
	// TODO: Implement applying changes
	return nil
}