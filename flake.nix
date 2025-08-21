{
  description = "patch2pr is a tool to create github prs without cloning the repos";

  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";

  outputs =
    {
      self,
      nixpkgs,
    }:
    let
      systems = nixpkgs.legacyPackages.x86_64-linux.go.meta.platforms;
      lib = nixpkgs.lib;
    in
    {
      packages = lib.genAttrs systems (
        sys:
        let
          pkgs = nixpkgs.legacyPackages.${sys};
        in
        rec {
          patch2pr = pkgs.callPackage ./. { };
          default = patch2pr;
        }
      );
    };
}
