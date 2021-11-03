#!/bin/bash

# Generate credentials
boots secrets set vix-username admin
boots secrets generate -l 12 -n 1 -s 1 vix-password

USER=$(boots secrets get vix-username)
PASS=$(boots secrets get vix-password)

# Setup vmrest
sudo ./expect.sh "$USER" "$PASS"