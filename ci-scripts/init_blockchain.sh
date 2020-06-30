#!/bin/bash

nodes=(
    "test-node-fra1-01"
    "test-node-fra1-02"
    "test-node-nyc3-01"
    "test-node-sgp1-01"
)

# Prepare all nodes
for node in "${nodes[@]}"
do
    # Remove previous run data
    rm -rf ~/.decimal-$node/daemon

    # Initialize new blockchain
    decd init $node --home ~/.decimal-$node/daemon --chain-id decimal-testnet-06-30-15-00

    # Add initial funds to the genesis file
    decd add-genesis-account dx16rr3cvdgj8jsywhx8lfteunn9uz0xg2c7ua9nl 40000011000000000000000000tdel --home ~/.decimal-$node/daemon # validator on test-node-fra1-01 (40,000,011 tDEL)
    decd add-genesis-account dx1ajytg8jg8ypx0rj9p792x32fuxyezga43jd3ry 40000011000000000000000000tdel --home ~/.decimal-$node/daemon # validator on test-node-fra1-02 (40,000,011 tDEL)
    decd add-genesis-account dx1azre0dtclv5y05ufynkhswzh0cwh4ktzlas3mp 40000011000000000000000000tdel --home ~/.decimal-$node/daemon # validator on test-node-nyc3-01 (40,000,011 tDEL)
    decd add-genesis-account dx1j3j2mwxnvlmsu2tkwm4z5390vq8v337wd6hmg2 40000011000000000000000000tdel --home ~/.decimal-$node/daemon # validator on test-node-sgp1-01 (40,000,011 tDEL)
    decd add-genesis-account dx12k95ukkqzjhkm9d94866r4d9fwx7tsd82r8pjd 160000000000000000000000000tdel --home ~/.decimal-$node/daemon # faucet (160,000,000 tDEL)
    decd add-genesis-account dx1esffyu0wxk6eez77fhzdxfgvjp4646hqm9sx6c 19999956000000000000000000tdel --home ~/.decimal-$node/daemon # tanker (19,999,956 tDEL)
done

# Add initial signed transactions to the genesis file
decd gentx \
    --name test-node-fra1-01 \
    --sequence 0 \
    --amount 40000000000000000000000000tdel \
    --fees 11000000000000000000tdel \
    --pubkey dxvalconspub1zcjduepquc5nas24rhqm0l8lyte0dfx2k3uda56wdn998lyrs6mpvsk9xmks6xa0ly \
    --details "Declaring validator on test-node-fra1-01" \
    --website decimalchain.com \
    --node-id bf7a6b366e3c451a3c12b3a6c01af7230fb92fc7 \
    --ip 139.59.133.148 \
    --keyring-backend test \
    --home ~/.decimal-test-node-fra1-01/daemon

# Add initial signed transactions to the genesis file
decd gentx \
    --name test-node-fra1-02 \
    --sequence 0 \
    --amount 40000000000000000000000000tdel \
    --fees 11000000000000000000tdel \
    --pubkey dxvalconspub1zcjduepq5hj3p750mves8wpmh4ywy6yjkz72sppr2kmzk7lzeedyelauwamsl3c57q \
    --details "Declaring validator on test-node-fra1-02" \
    --website decimalchain.com \
    --node-id c0b9b6c9a0f95e3d2f4aed890806739fc77faefd \
    --ip 64.225.110.228 \
    --keyring-backend test \
    --home ~/.decimal-test-node-fra1-02/daemon

# Add initial signed transactions to the genesis file
decd gentx \
    --name test-node-nyc3-01 \
    --sequence 0 \
    --amount 40000000000000000000000000tdel \
    --fees 11000000000000000000tdel \
    --pubkey dxvalconspub1zcjduepqjwlm5xcsp60v6fgwt95zq624yjhjnpkrzm209c5f8ajz8rdvq6gsdx5y2l \
    --details "Declaring validator on test-node-nyc3-01" \
    --website decimalchain.com \
    --node-id 76b81a4b817b39d63a3afe1f3a294f2a8f5c55b0 \
    --ip 64.225.56.107 \
    --keyring-backend test \
    --home ~/.decimal-test-node-nyc3-01/daemon

# Add initial signed transactions to the genesis file
decd gentx \
    --name test-node-sgp1-01 \
    --sequence 0 \
    --amount 40000000000000000000000000tdel \
    --fees 11000000000000000000tdel \
    --pubkey dxvalconspub1zcjduepq73se7rmlycftjta3ydjvrjmn28rrweyxyg42tzfr5lcw6lx8zl7qc6dpss \
    --details "Declaring validator on test-node-sgp1-01" \
    --website decimalchain.com \
    --node-id 29e566c41d51be90fa53340ba4edccefbebe8cb2 \
    --ip 139.59.192.48 \
    --keyring-backend test \
    --home ~/.decimal-test-node-sgp1-01/daemon

# Finish all nodes
for node in "${nodes[@]}"
do
    # Finish new blockchain initialization
    decd collect-gentxs --home ~/.decimal-$node/daemon
done
