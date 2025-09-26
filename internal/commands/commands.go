// go-repo/internal/commands.go

package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
)

func CreateRepoInteractive() error {
	var repoName string
	var visibility string
	var username string // Added missing variable declaration

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("GitHub Username").
				Description("Your GitHub username").
				Value(&username).
				Validate(validateNotEmpty),
			huh.NewInput().
				Title("Repository Name").
				Description("What should we call your new repository?").
				Value(&repoName).
				Validate(validateRepoName),
			huh.NewSelect[string]().
				Title("Visibility").
				Description("Should this be public or private?").
				Options(
					huh.NewOption("Public", "public"),
					huh.NewOption("Private", "private"),
				).
				Value(&visibility),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("form input failed: %w", err)
	}

	// Execute operations sequentially - only one return at the end
	fmt.Println("Creating local repository...")
	if err := CreateRepo(repoName); err != nil {
		return fmt.Errorf("local repository creation failed: %w", err)
	}

	fmt.Println("Creating remote repository on GitHub...")
	if err := CreateRemote(visibility, repoName); err != nil {
		return fmt.Errorf("remote repository creation failed: %w", err)
	}

	fmt.Println("Linking local and remote repositories...")
	if err := AddRemoteOrigin(repoName, username); err != nil {
		return fmt.Errorf("failed to add remote origin: %w", err)
	}

	fmt.Printf("✨ Repository '%s' created successfully!\n", repoName)
	fmt.Printf("📁 Local: ./%s\n", repoName)
	fmt.Printf("🌐 Remote: https://github.com/%s/%s\n", username, repoName)

	return nil
}

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

// Validation functions
func validateNotEmpty(s string) error {
	if s == "" {
		return fmt.Errorf("this field cannot be empty")
	}
	return nil
}

func validateRepoName(s string) error {
	if s == "" {
		return fmt.Errorf("repository name cannot be empty")
	}
	if len(s) > 100 {
		return fmt.Errorf("repository name too long")
	}
	return nil
}
