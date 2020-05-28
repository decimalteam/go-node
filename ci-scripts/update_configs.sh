#!/bin/sh

VALIDATOR_IP="139.59.133.148"
VALIDATOR_RPC="http://$VALIDATOR_IP/rpc"

export PATH=$PATH:$HOME/decimal
echo "Updating configs..."
echo "Step 1. Clear current state."
rm $HOME/.decimal -r

echo "Create basic configs."
decd init $(hostname) --chain-id=decimal-testnet-05-28-17-00
rm -r $HOME/.decimal/daemon/data
shopt -s extglob
cd $HOME/.decimal/daemon/config/
rm -v !("config.toml")
cd $HOME
shopt -u extglob

echo "Step 2. Downloading genesis."
curl "$VALIDATOR_RPC/genesis" | jq ".result.genesis" >$HOME/.decimal/daemon/config/genesis.json

echo "Step 3. Fetching node id."
curl "$VALIDATOR_RPC/status" | jq ".result.node_info.id" | awk '{gsub("\"", ""); print}' >$HOME/VALIDATOR_NODE_ID
echo "Got $(cat $HOME/VALIDATOR_NODE_ID) node id."
ORIGINAL_STR='persistent_peers = "*"'
REPLACE_STR='persistent_peers = "'$(cat $HOME/VALIDATOR_NODE_ID)@$VALIDATOR_IP':26656"'
REPLACE_RULE='s/'$ORIGINAL_STR'/'$REPLACE_STR'/g'

sed -i "$REPLACE_RULE" $HOME/.decimal/daemon/config/config.toml

echo "Enabling prometheus"
sed -i 's/prometheus = false/prometheus = true/g' $HOME/.decimal/daemon/config/config.toml
echo "Done!"
