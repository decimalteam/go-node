#!/bin/bash

# Add initial signed transactions to the genesis file
decd gen-declare-candidate-tx \
    --name test-node-fra1-01 \
    --sequence 0 \
    --amount 40000000000000000000000000tdel \
    --pubkey dxvalconspub1zcjduepquc5nas24rhqm0l8lyte0dfx2k3uda56wdn998lyrs6mpvsk9xmks6xa0ly \
    --details "Declaring validator on test-node-fra1-01" \
    --website decimalchain.com \
    --node-id bf7a6b366e3c451a3c12b3a6c01af7230fb92fc7 \
    --ip 139.59.133.148 \
    --chain-id decimal-testnet-06-09-13-00 \
    --keyring-backend test \
    --home ~/.decimal-test-node-fra1-01/daemon | jq '.value.signatures[0].signature'

# Add initial signed transactions to the genesis file
decd gen-declare-candidate-tx \
    --name test-node-fra1-02 \
    --sequence 0 \
    --amount 40000000000000000000000000tdel \
    --pubkey dxvalconspub1zcjduepq5hj3p750mves8wpmh4ywy6yjkz72sppr2kmzk7lzeedyelauwamsl3c57q \
    --details "Declaring validator on test-node-fra1-02" \
    --website decimalchain.com \
    --node-id c0b9b6c9a0f95e3d2f4aed890806739fc77faefd \
    --ip 64.225.110.228 \
    --chain-id decimal-testnet-06-09-13-00 \
    --keyring-backend test \
    --home ~/.decimal-test-node-fra1-02/daemon | jq '.value.signatures[0].signature'

# Add initial signed transactions to the genesis file
decd gen-declare-candidate-tx \
    --name test-node-nyc3-01 \
    --sequence 0 \
    --amount 40000000000000000000000000tdel \
    --pubkey dxvalconspub1zcjduepqjwlm5xcsp60v6fgwt95zq624yjhjnpkrzm209c5f8ajz8rdvq6gsdx5y2l \
    --details "Declaring validator on test-node-nyc3-01" \
    --website decimalchain.com \
    --node-id 76b81a4b817b39d63a3afe1f3a294f2a8f5c55b0 \
    --ip 64.225.56.107 \
    --chain-id decimal-testnet-06-09-13-00 \
    --keyring-backend test \
    --home ~/.decimal-test-node-nyc3-01/daemon | jq '.value.signatures[0].signature'

# Add initial signed transactions to the genesis file
decd gen-declare-candidate-tx \
    --name test-node-sgp1-01 \
    --sequence 0 \
    --amount 40000000000000000000000000tdel \
    --pubkey dxvalconspub1zcjduepq73se7rmlycftjta3ydjvrjmn28rrweyxyg42tzfr5lcw6lx8zl7qc6dpss \
    --details "Declaring validator on test-node-sgp1-01" \
    --website decimalchain.com \
    --node-id 29e566c41d51be90fa53340ba4edccefbebe8cb2 \
    --ip 167.99.182.218 \
    --chain-id decimal-testnet-06-09-13-00 \
    --keyring-backend test \
    --home ~/.decimal-test-node-sgp1-01/daemon | jq '.value.signatures[0].signature'

# # Finish all nodes
# for node in "${nodes[@]}"
# do
#     # Finish new blockchain initialization
#     decd collect-gentxs --home ~/.decimal-$node/daemon
# done
