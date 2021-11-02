#!/bin/bash

# Check for Nix
if ! command -v nix &> /dev/null
then
    echo "Installing Nix..."
    curl -o install-nix-2.3.16 https://releases.nixos.org/nix/nix-2.3.16/install
    curl -o install-nix-2.3.16.asc https://releases.nixos.org/nix/nix-2.3.16/install.asc
    gpg2 --recv-keys B541D55301270E0BCF15CA5D8170B4726D7198DE
    gpg2 --verify ./install-nix-2.3.16.asc
    bash ./install-nix-2.3.16

    echo "Cleaning up..."
    rm install-nix-2.3.16
    rm install-nix-2.3.16.asc
fi

docker build -t glab_dev -f dev/Dockerfile .
docker run -it -v $(pwd):/lab glab_dev /bin/bash