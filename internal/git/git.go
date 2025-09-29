package git

import (
	"fmt"
	"os"
	"os/exec"
)

func CreateRepo(repoName string) error {
	// Create directory with proper permissions
	err := os.MkdirAll(repoName, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", repoName, err)
	}

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = repoName // Execute git init INSIDE the new directory
	err = cmd.Run()
	if err != nil {
		// Cleanup on failure - atomic operation principle
		os.RemoveAll(repoName)
		return fmt.Errorf("failed to initialize git repo in %s: %w", repoName, err)
	}

	fmt.Printf("Repository '%s' created and initialized\n", repoName)
	return nil
}

func CreateRemote(visibility string, repoName string) error {
	// Build proper GitHub CLI command
	var cmd *exec.Cmd

	if visibility == "private" {
		cmd = exec.Command("gh", "repo", "create", repoName, "--private")
	} else {
		cmd = exec.Command("gh", "repo", "create", repoName, "--public")
	}

	// Execute and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create remote repository: %w\nOutput: %s", err, output)
	}

	fmt.Printf("Remote repository '%s' created on GitHub (%s)\n", repoName, visibility)
	return nil
}

func AddRemoteOrigin(repoName string, username string) error {
	repoURL := fmt.Sprintf("git@github.com:%s/%s.git", username, repoName)

	cmd := exec.Command("git", "remote", "add", "origin", repoURL)
	cmd.Dir = repoName

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to add remote origin: %w", err)
	}

	fmt.Printf("Remote origin added: %s\n", repoURL)
	return nil
}
