{
  description = "A quick scaffolding system for go projects";
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      packages.${system}.default = pkgs.buildGoModule {
        pname = "scaffold";
        version = "1.0.0";
        src = ./.;

        vendorHash = "sha256-kuheizBlwHUpozXPmFnoYOJ4zzyUDu3U5k8j0N1bhGc=";

        # Critical: Wrap binary with runtime dependencies
        nativeBuildInputs = [ pkgs.makeWrapper ];

        postInstall = ''
          # Rename binary
          mv $out/bin/go-repo $out/bin/scaffold

          # Wrap with required tools in PATH
          wrapProgram $out/bin/scaffold \
            --prefix PATH : ${pkgs.lib.makeBinPath [
              pkgs.git
              pkgs.gh
              pkgs.devenv
              pkgs.nix
              pkgs.direnv
              pkgs.age
            ]}
        '';
      };

      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          git
          gh
          devenv
          age
          direnv
        ];
      };
    };
}
