package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/huh"
	"github.com/jblais493/go-repo/internal/devenv"
	"github.com/jblais493/go-repo/internal/git"
	"github.com/jblais493/go-repo/internal/golang"
	"github.com/jblais493/go-repo/internal/secrets"
)

func CreateRepoInteractive() error {
	var repoName string
	var visibility string
	var username string

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

	fmt.Println("Configuring Github repository...")
	if err := git.CreateRepo(repoName); err != nil {
		return fmt.Errorf("local repository creation failed: %w", err)
	}
	if err := git.CreateRemote(visibility, repoName); err != nil {
		return fmt.Errorf("remote repository creation failed: %w", err)
	}
	if err := git.AddRemoteOrigin(repoName, username); err != nil {
		return fmt.Errorf("failed to add remote origin: %w", err)
	}

	fmt.Println("Setting up development environment with devenv...")
	if err := devenv.CreateDevenv(repoName); err != nil {
		return fmt.Errorf("devenv setup failed: %w", err)
	}

	fmt.Println("Initializing Go project structure...")
	modulePath := fmt.Sprintf("github.com/%s/%s", username, repoName)
	if err := golang.InitGoProject(repoName, modulePath); err != nil {
		return fmt.Errorf("Go project initialization failed: %w", err)
	}

	fmt.Println("Initializing secrets management...")
	if err := secrets.SecretsGen(repoName); err != nil {
		// Non-fatal? Or fatal?
		return fmt.Errorf("secrets initialization failed: %w", err)
	}

	fmt.Println("Initializing nix flake...")
	if err := flakeInit(repoName); err != nil {
		// Non-fatal? Or fatal?
		return fmt.Errorf("secrets initialization failed: %w", err)
	}

	fmt.Println("Creating first commit...")
	if err := git.FirstCommit(repoName, username); err != nil {
		return fmt.Errorf("local repository creation failed: %w", err)
	}

	fmt.Printf("✨ Repository '%s' created successfully!\n", repoName)
	fmt.Printf("📁 Local: ./%s\n", repoName)
	fmt.Printf("Secrets gerated at %s/secrets\n", repoName)
	fmt.Printf("🌐 Remote: https://github.com/%s/%s\n", username, repoName)
	fmt.Printf("🛠️  Development environment configured with devenv\n")
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  cd %s\n", repoName)
	fmt.Printf("  direnv allow\n")

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

func flakeInit(projectPath string) error {
	cmd := exec.Command("nix", "flake", "init")
	cmd.Dir = projectPath // Run in the new repo
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
