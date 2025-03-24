{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go_1_24
    pkgs.nodejs
    pkgs.yarn
  ];

  shellHook = ''
      yarn set version 4.5.3
      export PATH=$PATH:~/go/bin
  '';
}

