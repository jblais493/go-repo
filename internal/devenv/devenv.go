// go-repo/internal/devenv.go

package devenv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const devenvNixTemplate = `{ pkgs, lib, config, ... }:
let
  # Add secrets here like so:
  # analytics = { envVar = "ANALYTICS_URL"; };
  secretsConfig = {};

  # Function to load a secret with custom config
  loadSecretWithConfig = name: cfg: ''
    if [[ -f ./secrets/${name}.age ]]; then
      ${cfg.envVar}=$(age -d -i ~/.config/age/keys.txt ./secrets/${name}.age 2>/dev/null)
      if [[ $? -eq 0 ]]; then
        export ${cfg.envVar}
        echo "✓ Loaded ${name} → ${cfg.envVar}"
      else
        echo "❌ Failed to decrypt ${name}"
      fi
    else
      echo "⚠️  Secret file ./secrets/${name}.age not found"
    fi
  '';

  # Generate all secret loading commands
  loadAllSecrets = lib.concatStringsSep "\n"
    (lib.mapAttrsToList loadSecretWithConfig secretsConfig);
in
{
  # Add packages here:
  languages.go.enable = true;
  packages = [
    pkgs.air
    pkgs.age
    pkgs.just
  ];

  # Add runtime scripts here:
  scripts = {};

  # Additional processes can be added here:
  process.managers.process-compose.enable = true;
  processes = {
    go-dev.exec = "air";
  };

  # Show loaded development environment:
  enterShell = ''
    echo "🔐 Loading development secrets..."
    ${loadAllSecrets}
    echo "✅ Development environment ready"
    echo "📊 Available environment variables:"
    ${lib.concatStringsSep "\n" (lib.mapAttrsToList (name: cfg:
      ''echo "  ${cfg.envVar} (from ${name}.age)"''
    ) secretsConfig)}
  '';
}`

func CreateDevenv(projectPath string) error {
	// Step 1: Run 'devenv init' to initialize devenv structure
	if err := initDevenv(projectPath); err != nil {
		return fmt.Errorf("failed to initialize devenv: %w", err)
	}

	// Step 2: Write our custom devenv.nix configuration
	if err := writeDevenvConfig(projectPath); err != nil {
		return fmt.Errorf("failed to write devenv.nix: %w", err)
	}

	return nil
}

func initDevenv(projectPath string) error {
	cmd := exec.Command("devenv", "init")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func writeDevenvConfig(projectPath string) error {
	devenvPath := filepath.Join(projectPath, "devenv.nix")

	return os.WriteFile(devenvPath, []byte(devenvNixTemplate), 0644)
}
