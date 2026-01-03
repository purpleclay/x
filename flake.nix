{
  description = "Experimental Purple Clay libraries";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";

    git-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
      };
    };

    go-overlay = {
      url = "github:purpleclay/go-overlay";
      inputs = {
        nixpkgs.follows = "nixpkgs";
      };
    };
  };

  outputs = {
    nixpkgs,
    flake-utils,
    git-hooks,
    go-overlay,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [go-overlay.overlays.default];
        };

        buildInputs = with pkgs; [
          alejandra
          go-bin.versions."1.24.11"
          gofumpt
          golangci-lint
          go-overlay.packages.${system}.govendor
          nil
          typos
        ];

        pre-commit-check = git-hooks.lib.${system}.run {
          src = ./.;
          package = pkgs.prek;
          hooks = {
            alejandra = {
              enable = true;
              settings = {
                check = true;
              };
            };

            typos = {
              enable = true;
              entry = "${pkgs.typos}/bin/typos";
            };
          };
        };
      in
        with pkgs; {
          checks = {
            inherit pre-commit-check;
          };

          packages = {
            theme = callPackage ./default.nix {
              inherit buildGoWorkspace;
              go = go-bin.fromGoMod ./theme/go.mod;
            };
          };

          devShells.default = mkShell {
            inherit (pre-commit-check) shellHook;
            buildInputs = buildInputs ++ pre-commit-check.enabledPackages;
          };
        }
    );
}
