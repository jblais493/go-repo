package main

import (
	cmd "github.com/jblais493/go-repo/internal/commands"
	"log"
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
		log.Fatal(err)
	}
}
