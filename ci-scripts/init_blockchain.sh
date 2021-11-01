#!/bin/bash

# Add initial signed transactions to the genesis file
decd gentx dev-node-01 40000000000000000000000000del \
    --chain-id decimal-devnet-28-10-16-00 \
    --pubkey dxvalconspub1zcjduepqcl2c373fljlpm0lfut5c4adq4x08fsd9ufqm6fsq3e6p4fqg7xeqkx0m98 \
    --moniker dev-node-01 \
    --details "Declaring validator on dev-node-01" \
    --website decimalchain.com \
    --node-id 27fcfef145b3717c5d639ec72fb12f9c43da98f0 \
    --ip "185.242.122.122" \
    --log_level info

# Add initial signed transactions to the genesis file
decd gentx dev-node-02 40000000000000000000000000del \
    --chain-id decimal-devnet-28-10-16-00 \
    --pubkey dxvalconspub1zcjduepqwwrzqe9dg6tq2s0ps9sl0pqu42x3ypqjyykfaems3stnhpg30agqjczlcj \
    --moniker dev-node-02 \
    --details "Declaring validator on dev-node-02" \
    --website decimalchain.com \
    --node-id 8a2cc38f5264e9699abb8db91c9b4a4a061f000d \
    --ip "185.242.122.123" \
    --log_level info
# Add initial signed transactions to the genesis file
decd gentx dev-node-03 40000000000000000000000000del \
    --chain-id decimal-devnet-28-10-16-00 \
    --pubkey dxvalconspub1zcjduepquhdcn6578xh37gpmwn89vlq8cu402gm5nvnkels3kpnz9a3gcyxqgexz4q \
    --moniker dev-node-03 \
    --details "Declaring validator on dev-node-03" \
    --website decimalchain.com \
    --node-id e0e7a88de0b39bd2adceb3516d353582ff94ec15 \
    --ip "185.242.122.124" \
    --log_level info