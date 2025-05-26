package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"golang.org/x/term"
)

// Shows interactive picker and returns selected labels
func ShowPicker(sourceLabels, destLabels []Label, destRepo string, verbose bool) ([]Label, error) {
	// Check if we're in an interactive terminal
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, fmt.Errorf("interactive picker requires a terminal")
	}
	
	items := buildPickerItems(sourceLabels, destLabels)
	
	fmt.Printf("\nCurrent state → Desired state for %s:\n\n", destRepo)
	
	// Use a simple select/deselect loop instead of promptui's Select
	// because we need multi-selection with toggling
	selectedMap := make(map[string]bool)
	for i, item := range items {
		selectedMap[fmt.Sprintf("%d", i)] = item.Selected
	}
	
	currentIndex := 0
	
	for {
		// Clear screen and redraw (more compatible)
		fmt.Print("\033[2J\033[H")
		fmt.Printf("Current state → Desired state for %s:\n\n", destRepo)
		
		// Show separator after dest-only labels
		lastDestOnly := -1
		for i, item := range items {
			if item.IsDestOnly {
				lastDestOnly = i
			}
		}
		
		// Display all items
		for i, item := range items {
			if lastDestOnly >= 0 && i == lastDestOnly+1 {
				fmt.Println("  ────────────────────────────────────────────────")
			}
			
			selected := selectedMap[fmt.Sprintf("%d", i)]
			checkbox := "[ ]"
			if selected {
				checkbox = "[✓]"
			}
			
			cursor := "  "
			if i == currentIndex {
				cursor = "> "
			}
			
			label := FormatLabel(item.Label, verbose)
			if item.IsDestOnly {
				label += " (dest only)"
				if !selected {
					label += " [WARN] will be deleted"
				}
			}
			
			fmt.Printf("%s%s %s\n", cursor, checkbox, label)
		}
		
		selectedCount := 0
		for _, sel := range selectedMap {
			if sel {
				selectedCount++
			}
		}
		
		fmt.Printf("\n  Space: toggle  a: toggle all  ↑/↓: navigate  Enter: confirm  q: quit\n")
		fmt.Printf("\n  %d selected\n", selectedCount)
		
		// Get single keypress
		key := getKeypress()
		
		switch key {
		case 'q', 'Q':
			return nil, fmt.Errorf("cancelled")
		case '\n', '\r': // Enter
			var selected []Label
			for i, item := range items {
				if selectedMap[fmt.Sprintf("%d", i)] {
					selected = append(selected, item.Label)
				}
			}
			return selected, nil
		case ' ': // Space
			key := fmt.Sprintf("%d", currentIndex)
			selectedMap[key] = !selectedMap[key]
		case 'a', 'A': // Toggle all
			// Count how many are currently selected
			allSelected := true
			for _, v := range selectedMap {
				if !v {
					allSelected = false
					break
				}
			}
			// Toggle all to opposite state
			for k := range selectedMap {
				selectedMap[k] = !allSelected
			}
		case '\x1b': // Escape sequence
			// Read the rest of the escape sequence
			getKeypress() // [
			switch getKeypress() {
			case 'A': // Up arrow
				if currentIndex > 0 {
					currentIndex--
				}
			case 'B': // Down arrow
				if currentIndex < len(items)-1 {
					currentIndex++
				}
			}
		}
	}
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
		if len(summary.ToCreate) == 1 {
			fmt.Printf("  • Create 1 label\n")
		} else {
			fmt.Printf("  • Create %d labels\n", len(summary.ToCreate))
		}
	}
	if len(summary.ToDelete) > 0 {
		if len(summary.ToDelete) == 1 {
			fmt.Printf("  • Delete 1 label (%s)\n", summary.ToDelete[0].Name)
		} else {
			fmt.Printf("  • Delete %d labels\n", len(summary.ToDelete))
		}
	}
	if len(summary.ToKeep) > 0 {
		if len(summary.ToKeep) == 1 {
			fmt.Printf("  • Keep 1 existing label\n")
		} else {
			fmt.Printf("  • Keep %d existing labels\n", len(summary.ToKeep))
		}
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

// Gets a single keypress from the terminal
func getKeypress() byte {
	// Put terminal in raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return 0
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	b := make([]byte, 1)
	_, _ = os.Stdin.Read(b)
	return b[0]
}

// Calculates what actions need to be taken
func calculateActions(selectedLabels, destLabels []Label) ActionSummary {
	selectedMap := make(map[string]bool)
	destMap := make(map[string]Label)
	
	// Index selected labels by lowercase name
	for _, label := range selectedLabels {
		selectedMap[strings.ToLower(label.Name)] = true
	}
	
	// Index destination labels by lowercase name
	for _, label := range destLabels {
		destMap[strings.ToLower(label.Name)] = label
	}
	
	summary := ActionSummary{
		ToCreate: []Label{},
		ToDelete: []Label{},
		ToKeep:   []Label{},
	}
	
	// Find labels to create (in selected but not in dest)
	for _, label := range selectedLabels {
		if _, exists := destMap[strings.ToLower(label.Name)]; !exists {
			summary.ToCreate = append(summary.ToCreate, label)
		} else {
			summary.ToKeep = append(summary.ToKeep, label)
		}
	}
	
	// Find labels to delete (in dest but not selected)
	for _, label := range destLabels {
		if !selectedMap[strings.ToLower(label.Name)] {
			summary.ToDelete = append(summary.ToDelete, label)
		}
	}
	
	return summary
}

// Applies the changes to the destination repository
func applyChanges(summary ActionSummary, destRepo string) error {
	totalOps := len(summary.ToDelete) + len(summary.ToCreate)
	currentOp := 0
	
	// Delete labels
	for _, label := range summary.ToDelete {
		currentOp++
		fmt.Printf("[%d/%d] Deleting %s...\n", currentOp, totalOps, label.Name)
		if err := DeleteLabel(destRepo, label.Name); err != nil {
			return fmt.Errorf("failed to delete label %s: %v", label.Name, err)
		}
	}
	
	// Create labels
	for _, label := range summary.ToCreate {
		currentOp++
		fmt.Printf("[%d/%d] Creating %s...\n", currentOp, totalOps, label.Name)
		if err := CreateLabel(destRepo, label); err != nil {
			return fmt.Errorf("failed to create label %s: %v", label.Name, err)
		}
	}
	
	fmt.Printf("\nDone! ")
	if len(summary.ToCreate) > 0 {
		if len(summary.ToCreate) == 1 {
			fmt.Printf("Created 1 label. ")
		} else {
			fmt.Printf("Created %d labels. ", len(summary.ToCreate))
		}
	}
	if len(summary.ToDelete) > 0 {
		if len(summary.ToDelete) == 1 {
			fmt.Printf("Deleted 1 label. ")
		} else {
			fmt.Printf("Deleted %d labels. ", len(summary.ToDelete))
		}
	}
	if len(summary.ToKeep) > 0 {
		if len(summary.ToKeep) == 1 {
			fmt.Printf("Kept 1 existing label.")
		} else {
			fmt.Printf("Kept %d existing labels.", len(summary.ToKeep))
		}
	}
	fmt.Println()
	
	return nil
}