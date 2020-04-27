#!/bin/bash

# WARNING: Do not use "test" keyring backend in a production!
# It is used not to radically simplify testnet CI.

# Create validator key pair
deccli keys add val --keyring-backend test

# Initialize new blockchain
decd init mynode --chain-id decimal-testnet

# Add initial funds to the genesis file
decd add-genesis-account $(deccli keys show val -a --keyring-backend test) 100000000000000000tdcl
decd add-genesis-account dx1fyqf7gp0gzmpzwxfrah9veaxt8ysl68khxwfjm 1000000000000000000000000000tdcl
decd add-genesis-account dx1dvwgj5hc3uqjemk22v9gq7um30v0l6wpcwnwnm 1000000000000000000000000000tdcl
decd add-genesis-account dx12k95ukkqzjhkm9d94866r4d9fwx7tsd82r8pjd 1000000000000000000000000000tdcl
decd add-genesis-account dx1cjkq662yycfy03euktm4p5rem0q53f9x89crap 1000000000000000000000000000tdcl
decd add-genesis-account dx178h6nvsqg5vr2gq8xv8f9jlmkpg3mjvzktsmj2 1000000000000000000000000000tdcl

# Add initial signed transactions to the genesis file
decd gentx --name val --website decimalchain.com --keyring-backend test

# Configure created blockchain
deccli config chain-id decimal-testnet
deccli config trust-node true

# Finish new blockchain initialization
decd collect-gentxs
