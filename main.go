package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	debug   bool
)

var rootCmd = &cobra.Command{
	Use:     "gabel source-repo dest-repo",
	Short:   "Safely copy GitHub labels between repositories",
	Long:    "Gabel helps you copy GitHub labels from one repo to another with an interactive picker.",
	Version: Version,
	Args:    cobra.ExactArgs(2),
	Run:     run,
}

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show label descriptions")
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Show debug logs")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	sourceRepo := args[0]
	destRepo := args[1]

	InitLogger(debug)

	if !isValidRepo(sourceRepo) || !isValidRepo(destRepo) {
		fmt.Fprintf(os.Stderr, "Error: Invalid repo format. Use 'owner/repo' format.\n")
		os.Exit(1)
	}

	if err := CheckGitHubCLI(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	LogDebug("Source repo: %s", sourceRepo)
	LogDebug("Destination repo: %s", destRepo)

	fmt.Printf("Fetching labels from %s...\n", sourceRepo)
	sourceLabels, err := FetchLabels(sourceRepo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching labels from %s: %v\n", sourceRepo, err)
		os.Exit(1)
	}

	if len(sourceLabels) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No labels found in %s\n", sourceRepo)
		os.Exit(1)
	}

	fmt.Printf("Fetching labels from %s...\n", destRepo)
	destLabels, err := FetchLabels(destRepo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching labels from %s: %v\n", destRepo, err)
		os.Exit(1)
	}

	LogDebug("Found %d labels in source, %d labels in destination", len(sourceLabels), len(destLabels))

	selectedLabels, err := ShowPicker(sourceLabels, destLabels, destRepo, verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(selectedLabels) == 0 {
		fmt.Println("No labels selected. Nothing to do.")
		return
	}

	if err := ConfirmAndApply(selectedLabels, destLabels, destRepo); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func isValidRepo(repo string) bool {
	parts := strings.Split(repo, "/")
	return len(parts) == 2 && parts[0] != "" && parts[1] != ""
}