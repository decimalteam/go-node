#!/bin/bash

nodes=(
    "dev-node-fra1-01"
    "dev-node-fra1-02"
    "dev-node-tor1-01"
)

# Prepare all nodes
for node in "${nodes[@]}"
do
    # Remove previous run data
    rm -rf ~/.decimal-$node/daemon

    # Initialize new blockchain
    decd init $node --home ~/.decimal-$node/daemon --chain-id decimal-devnet-07-30-13-40

    # Add initial funds to the genesis file
    decd add-genesis-account dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v  40000000000000000000000000del --home ~/.decimal-$node/daemon # validator on fra1-01  ( 40,000,000 DEL)
    decd add-genesis-account dx1mvqrrrlcd0gdt256jxg7n68e4neppu5t24e8h6  40000000000000000000000000del --home ~/.decimal-$node/daemon # validator on fra1-02  ( 40,000,000 DEL)
    decd add-genesis-account dx1nrr6er27mmcufmaqm4dyu6c5r6489cfm35m4ft  40000000000000000000000000del --home ~/.decimal-$node/daemon # validator on tor1-01  ( 40,000,000 DEL)
    decd add-genesis-account dx1tvqxh4x7pedyqpzqp9tdf068k4q9j2hm3lmghl 200000000000000000000000000del --home ~/.decimal-$node/daemon # faucet                (200,000,000 DEL)
    decd add-genesis-account dx1mtlnpmwf8zr6pek6gq25nv45x2890sne2ap0cc  20000000000000000000000000del --home ~/.decimal-$node/daemon # tanker                ( 20,000,000 DEL)
done

# Add initial signed transactions to the genesis file
decd gentx \
    --name dev-node-fra1-01 \
    --sequence 0 \
    --amount 40000000000000000000000000del \
    --pubkey dxvalconspub1zcjduepqwwrzqe9dg6tq2s0ps9sl0pqu42x3ypqjyykfaems3stnhpg30agqjczlcj \
    --details "Declaring validator on dev-node-fra1-01" \
    --website decimalchain.com \
    --node-id 8a2cc38f5264e9699abb8db91c9b4a4a061f000d \
    --ip 46.101.127.241 \
    --keyring-backend test \
    --home ~/.decimal-dev-node-fra1-01/daemon

# Add initial signed transactions to the genesis file
decd gentx \
    --name dev-node-fra1-02 \
    --sequence 0 \
    --amount 40000000000000000000000000del \
    --pubkey dxvalconspub1zcjduepquhdcn6578xh37gpmwn89vlq8cu402gm5nvnkels3kpnz9a3gcyxqgexz4q \
    --details "Declaring validator on dev-node-fra1-02" \
    --website decimalchain.com \
    --node-id e0e7a88de0b39bd2adceb3516d353582ff94ec15 \
    --ip 164.90.211.234 \
    --keyring-backend test \
    --home ~/.decimal-dev-node-fra1-02/daemon

# Add initial signed transactions to the genesis file
decd gentx \
    --name dev-node-tor1-01 \
    --sequence 0 \
    --amount 40000000000000000000000000del \
    --pubkey dxvalconspub1zcjduepqcl2c373fljlpm0lfut5c4adq4x08fsd9ufqm6fsq3e6p4fqg7xeqkx0m98 \
    --details "Declaring validator on dev-node-tor1-01" \
    --website decimalchain.com \
    --node-id 27fcfef145b3717c5d639ec72fb12f9c43da98f0 \
    --ip 167.99.182.218 \
    --keyring-backend test \
    --home ~/.decimal-dev-node-tor1-01/daemon

# Finish all nodes
for node in "${nodes[@]}"
do
    # Finish new blockchain initialization
    decd collect-gentxs --home ~/.decimal-$node/daemon
done
