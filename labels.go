package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

// Verifies gh CLI is installed and authenticated
func CheckGitHubCLI() error {
	_, err := exec.LookPath("gh")
	if err != nil {
		return fmt.Errorf("GitHub CLI (gh) is required but not found.\nInstall it from: https://cli.github.com")
	}

	cmd := exec.Command("gh", "auth", "status")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not authenticated with GitHub.\nRun: gh auth login")
	}

	return nil
}

// Retrieves all labels from a repository
func FetchLabels(repo string) ([]Label, error) {
	LogDebug("Fetching labels from %s", repo)

	cmd := exec.Command("gh", "api", fmt.Sprintf("repos/%s/labels", repo), "--paginate")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr := string(exitErr.Stderr)
			if strings.Contains(stderr, "HTTP 404") {
				return nil, fmt.Errorf("repository not found: %s", repo)
			}
			if strings.Contains(stderr, "HTTP 403") {
				if strings.Contains(stderr, "rate limit") {
					return nil, fmt.Errorf("GitHub API rate limit exceeded. Please wait and try again")
				}
				return nil, fmt.Errorf("access denied. You may not have permission to view %s", repo)
			}
			return nil, fmt.Errorf("GitHub API error: %s", stderr)
		}
		return nil, err
	}

	var labels []Label
	if err := json.Unmarshal(output, &labels); err != nil {
		return nil, fmt.Errorf("failed to parse labels: %v", err)
	}

	LogDebug("Fetched %d labels from %s", len(labels), repo)
	return labels, nil
}

func CreateLabel(repo string, label Label) error {
	LogDebug("Creating label '%s' in %s", label.Name, repo)
	
	// Validate before creation
	if err := validateLabel(label); err != nil {
		return err
	}
	
	color, _ := validateColor(label.Color) // Already validated above

	args := []string{
		"api",
		fmt.Sprintf("repos/%s/labels", repo),
		"-f", fmt.Sprintf("name=%s", label.Name),
		"-f", fmt.Sprintf("color=%s", color),
		"-f", fmt.Sprintf("description=%s", label.Description),
	}

	cmd := exec.Command("gh", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func UpdateLabel(repo string, label Label) error {
	LogDebug("Updating label '%s' in %s", label.Name, repo)

	args := []string{
		"api",
		fmt.Sprintf("repos/%s/labels/%s", repo, url.PathEscape(label.Name)),
		"--method", "PATCH",
		"-f", fmt.Sprintf("color=%s", strings.TrimPrefix(label.Color, "#")),
		"-f", fmt.Sprintf("description=%s", label.Description),
	}

	cmd := exec.Command("gh", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DeleteLabel(repo string, labelName string) error {
	LogDebug("Deleting label '%s' from %s", labelName, repo)

	cmd := exec.Command("gh", "api",
		fmt.Sprintf("repos/%s/labels/%s", repo, url.PathEscape(labelName)),
		"--method", "DELETE")
	cmd.Stderr = os.Stderr
	return cmd.Run()
}