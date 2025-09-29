// go-repo/internal/golang/golang.go

package golang

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const mainGoTemplate = `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`

const airTomlTemplate = `root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ."
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
`

// InitGoProject orchestrates all Go-specific initialization
func InitGoProject(projectPath, modulePath string) error {
	if err := initModule(projectPath, modulePath); err != nil {
		return fmt.Errorf("failed to initialize Go module: %w", err)
	}

	if err := createMainFile(projectPath); err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}

	if err := initAir(projectPath); err != nil {
		return fmt.Errorf("failed to initialize Air: %w", err)
	}

	return nil
}

func initModule(projectPath, modulePath string) error {
	cmd := exec.Command("go", "mod", "init", modulePath)
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Printf("Go module initialized: %s\n", modulePath)
	return nil
}

func createMainFile(projectPath string) error {
	mainPath := filepath.Join(projectPath, "main.go")

	if err := os.WriteFile(mainPath, []byte(mainGoTemplate), 0644); err != nil {
		return err
	}

	fmt.Printf("Created main.go with Hello World boilerplate\n")
	return nil
}

func initAir(projectPath string) error {
	airConfigPath := filepath.Join(projectPath, ".air.toml")

	// Write our custom air configuration
	if err := os.WriteFile(airConfigPath, []byte(airTomlTemplate), 0644); err != nil {
		return err
	}

	fmt.Printf("Air hot-reload configured (.air.toml created)\n")
	return nil
}
