echo "rm -rf ~/.decimal/daemon"
sudo rm -rf ~/.decimal/daemon

echo "decd init mynode --chain-id decimal"
decd init mynode --chain-id decimal

echo "deccli keys add test1"
deccli keys add test1

echo "deccli keys add test2"
deccli keys add test2

echo "echo \"12345678\" | decd add-genesis-account $(deccli keys show test1 -a) 1000000000000000000000000tdel"
echo "12345678" | decd add-genesis-account $(deccli keys show test1 -a) 1000000000000000000000000del

echo "echo \"12345678\" | decd add-genesis-account $(deccli keys show test2 -a) 1000000000000000000000000tdel"
echo "12345678" | decd add-genesis-account $(deccli keys show test2 -a) 1000000000000000000000000del

echo "deccli config chain-id decimal"
deccli config chain-id decimal

echo "deccli config output json"
deccli config output json

echo "deccli config indent true"
deccli config indent true

echo "deccli config trust-node true"
deccli config trust-node true

echo "decd gentx --name test1"
decd gentx --name test1

echo "decd collect-gentxs"
decd collect-gentxs

echo "decd updater url bech32"
decd updater http://185.242.122.122/file dx1s53f6fnhz0vd4ah7639jzfq2nvfmzq9xqqvp05
