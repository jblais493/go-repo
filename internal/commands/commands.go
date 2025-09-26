// go-repo/internal/commands.go

package commands

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"os"
	"os/exec"
)

func CreateRepoInteractive() error {
	var repoName string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Repository Name").
				Value(&repoName),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	return CreateRepo(repoName)
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
