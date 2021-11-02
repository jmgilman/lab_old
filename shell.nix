{
  pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/f151bae7ab9bfdf6f125a2a03aad0ba6bf0589dc.tar.gz") {}
}:

let
  boots = import cli/boots.nix {
    lib = pkgs.lib;
    buildGoModule = pkgs.buildGoModule;
  };
in
pkgs.mkShell {
  buildInputs = [
    boots
    pkgs.consul
    pkgs.nomad
    pkgs.vault
    pkgs.vagrantS
  ];
}