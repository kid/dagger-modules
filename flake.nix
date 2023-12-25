{
  inputs = {
    nixpkgs.url = "github:NixOs/nixpkgs/nixos-unstable";
    devenv.url = "github:cachix/devenv";
    dagger.url = "github:dagger/nix";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.devenv.flakeModule
      ];

      systems = [ "x86_64-linux" ];

      perSystem = { inputs', pkgs, ... }: {
        devenv.shells = {
          nix = {
            packages = with pkgs; [
              nil
              rnix-lsp
            ];
          };

          dagger = {
            packages = [ inputs'.dagger.packages.dagger ];
          };
        };
      };
    };
}
