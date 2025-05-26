package main

// Label represents a GitHub label
type Label struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// PickerItem represents a label in the picker with selection state
type PickerItem struct {
	Label      Label
	Selected   bool
	IsDestOnly bool
}

// ActionSummary describes what will happen to labels
type ActionSummary struct {
	ToCreate []Label
	ToDelete []Label
	ToKeep   []Label
}