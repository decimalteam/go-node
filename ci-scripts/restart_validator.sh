echo "Current node is a validator."
rm -r ~/.decimal/daemon/data/*.db
rm -r ~/.decimal/daemon/data/*.wal
rm ~/.decimal/daemon/config/write-file-atomic-*
rm -r ~/.decimal/daemon/config/gentx
echo "Wipe priv_validator_state."
echo $(cat ~/.decimal/daemon/data/priv_validator_state.json | jq '.height="0"' | jq '.step=0') >~/.decimal/daemon/data/priv_validator_state.json
echo "Restating daemon."
sudo systemctl restart decd
