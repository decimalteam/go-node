#!/bin/bash
rm -r ~/.decimal/daemon
mkdir ~/.decimal/daemon/config/gentx
decd init mynode --chain-id decimal
deccli keys add val
deccli keys add test1
echo "12345678" | decd add-genesis-account $(deccli keys show val -a) 1000000000000000000000000del 
echo "12345678" | decd add-genesis-account $(deccli keys show test1 -a) 1000000000000000000000000del 
deccli config chain-id decimal
deccli config output json
deccli config indent true
deccli config trust-node true
decd gentx --name val
decd collect-gentxs