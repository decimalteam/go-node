# Decimal Go Node

## Requirements

- [`git`](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [`golang` 1.14+](https://golang.org/doc/install)
- shell tools [`curl`](https://curl.haxx.se/download.html) and [`jq`](https://stedolan.github.io/jq/download/)
- building essentials

To install building essentials (which include [`make`](https://www.gnu.org/software/make/)) use following commands:

```bash
# Ubuntu:
sudo apt-get install build-essential

# macOS:
brew install coreutils
```

## Installing

Clone repository

```bash
git clone https://bitbucket.org/decimalteam/go-node.git
cd go-node
```

Build and install Decimal Go Node from source code

```bash
make all
```

Confirm `decd` and `deccli` are built and installed properly. For that use help command to retrieve `decd` and `deccli` usage information

```bash
decd --help
deccli --help
```

## Configuring

First of all, make sure directory at path `"$HOME/.decimal/daemon"` does not exist.

***WARNING*** *It will remove your current Decimal blockchain state if exists!*

```bash
rm -rf "$HOME/.decimal/daemon"
```

Time to determine proper chain ID currently used in the network and initialize new Decimal network node

```bash
NODE_MONIKER="$USER-node" # You are free to choose other name for your node
CHAIN_ID="$(curl -s 'https://testnet-gate.decimalchain.com/api/rpc/genesis/chain')"
decd init "$NODE_MONIKER" --chain-id "$CHAIN_ID"
```

Download proper `genesis.json` from master node

```bash
curl -s 'https://testnet-gate.decimalchain.com/api/rpc/genesis' | jq '.result.genesis' >> "$HOME/.decimal/daemon/config/genesis.json"
```

Add proper `persistent_peers` to `config.toml` file

```toml
# Comma separated list of nodes to keep persistent connections to
persistent_peers = "bf7a6b366e3c451a3c12b3a6c01af7230fb92fc7@139.59.133.148:26656"
```

## Running

To run Decimal node it is enough to exec command

```bash
decd start
```

Decimal node required some time to sync blockchain on new deployed node so it is time to take a breath. Enjoy!

## Validating

Once your Decimal node is synced and in actual state, it becomes possible to participate in block generating process and earn some coins.

TODO: To be continued...
