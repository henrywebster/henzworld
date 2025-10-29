{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    nativeBuildInputs = with pkgs.buildPackages; [
      go
      air
      flyctl
      sqlite
      dbmate
      djlint
    ];

    shellHook = ''
      if [ -f .env ]; then
        set -a
        source .env
        set +a
      fi

      export PATH=$PATH:$(go env GOPATH)/bin
    '';
}
