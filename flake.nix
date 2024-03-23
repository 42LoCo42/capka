{
  outputs = { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        inherit (pkgs.lib.fileset) toSource unions;
      in
      rec {
        packages = rec {
          default = capka;
          capka = pkgs.buildGoModule {
            pname = "capka";
            version = "0.1.0";
            src = toSource {
              root = ./.;
              fileset = unions [
                ./go.mod
                ./go.sum

                ./capka.c
                ./capka.h

                ./capka.go
                ./capka_test.go

                ./cmd
              ];
            };

            vendorHash = "sha256-51wP3ov6rOYnzQmhGr5Hj8j8RQRnB5GB7YoYtcEratk=";
            nativeBuildInputs = with pkgs; [ pkg-config ];
            buildInputs = with pkgs; [ libsodium ];
          };
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
