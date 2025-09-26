package main

import (
	"context"
	"fmt"
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
	fmt.Printf("hello")
}
