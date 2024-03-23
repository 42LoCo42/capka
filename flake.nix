{
  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        inherit (pkgs.lib.fileset) toSource unions;
      in
      rec {
        packages.capka-lib = pkgs.buildGoModule {
          pname = "capka-lib";
          version = "0.1.0";
          src = toSource {
            root = ./.;
            fileset = unions [
              ./go.mod
              ./go.sum
              ./lib
            ];
          };

          vendorHash = "sha256-eo9820j+WJehClFpw+iEOZkS4BpDmNd73DlkeoiTEcY=";
          nativeBuildInputs = with pkgs; [ pkg-config ];
          buildInputs = with pkgs; [ libsodium ];
        };

        devShells.default = pkgs.mkShell {
          inputsFrom = builtins.attrValues packages;

          packages = with pkgs; [
            gopls
          ];
        };
      }
    );
}
