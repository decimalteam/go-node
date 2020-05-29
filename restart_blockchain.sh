#!/bin/bash

# Rebuild decd/deccli binaries
make all

# Remove $HOME/.decimal directory
rm -rf "$HOME/.decimal"

# Copy prepared $HOME/.decimal directory
cp -r "$HOME/.decimal-test" "$HOME/.decimal"

# Start decd
decd start