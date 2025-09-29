package secrets

import (
	"os"
	"os/exec"
)

func secretsGen(projectPath string) error {
	cmd := exec.Command("nix", "run", "github:jblais493/go-secrets", "--", "generate")
	cmd.Dir = projectPath // Run in the new repo
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
