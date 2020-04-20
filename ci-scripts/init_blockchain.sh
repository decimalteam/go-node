#!/bin/bash

# Set validator password to environment variables to avoid copy-paste
export DECIMAL_VALIDATOR_PASSWORD=9n-yg8beVQzuh7Th

# Create validator key pair
deccli keys add val <<EOF
$DECIMAL_VALIDATOR_PASSWORD
$DECIMAL_VALIDATOR_PASSWORD
EOF

# Initialize new blockchain
decd init mynode --chain-id decimal-testnet
decd add-genesis-account dx1fyqf7gp0gzmpzwxfrah9veaxt8ysl68khxwfjm 1000000000000000000000000000tdcl
decd add-genesis-account dx1dvwgj5hc3uqjemk22v9gq7um30v0l6wpcwnwnm 1000000000000000000000000000tdcl
decd add-genesis-account $(echo "$DECIMAL_VALIDATOR_PASSWORD" | deccli keys show val -a) 100000000000000000tdcl

# Add initial signed transactions to the genesis file
# TODO: It does not work for now!
decd gentx --name val --website decimalchain.com <<EOF
$DECIMAL_VALIDATOR_PASSWORD
$DECIMAL_VALIDATOR_PASSWORD
$DECIMAL_VALIDATOR_PASSWORD
EOF

# Configure created blockchain
deccli config chain-id decimal-testnet
deccli config trust-node true

# Finish new blockchain initialization
decd collect-gentxs

# Unset validator password from environment variables
unset DECIMAL_VALIDATOR_PASSWORD