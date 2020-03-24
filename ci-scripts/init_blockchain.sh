#!/bin/bash
echo "9n-yg8beVQzuh7Th" | deccli keys add val
decd init mynode --chain-id decimal-testnet
decd add-genesis-account dx1fyqf7gp0gzmpzwxfrah9veaxt8ysl68khxwfjm 1000000000000000000000000000tdcl
decd add-genesis-account dx1dvwgj5hc3uqjemk22v9gq7um30v0l6wpcwnwnm 1000000000000000000000000000tdcl
decd add-genesis-account $(deccli keys show val -a) 100000000000000000stake
echo "9n-yg8beVQzuh7Th" | decd gentx --name val --website decimalchain.com
deccli config chain-id decimal-testnet
deccli config trust-node true
decd collect-gentxs