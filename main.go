// main.go

package main

import (
	"fmt"
	cmd "github.com/jblais493/go-repo/internal/commands"
	"os"
)

type Init struct {
	Username string
}

type RepoConfig struct {
	RepoName    string
	Path        string
	Visibility  string
	Description string
	License     string
}

func main() {
	if err := cmd.CreateRepoInteractive(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
