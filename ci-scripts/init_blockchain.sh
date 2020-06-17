#!/bin/bash
export PATH=$PATH:$HOME/decimal
rm -r ~/.decimal
# WARNING: Do not use "test" keyring backend in a production!
# It is used not to radically simplify testnet CI.

# Create validator key pair
deccli keys add val --keyring-backend test
deccli keys add spammer --keyring-backend test

# rm -rf ~/.decimal/daemon

# Initialize new blockchain
decd init test-node-fra1-01 --chain-id decimal-testnet-06-17-14-00
decd init test-node-fra1-02 --chain-id decimal-testnet-06-17-14-00
decd init test-node-nyc3-01 --chain-id decimal-testnet-06-17-14-00
decd init test-node-sgp1-01 --chain-id decimal-testnet-06-17-14-00

# Add initial funds to the genesis file
# decd add-genesis-account $(deccli keys show val -a --keyring-backend test) 100000000000000000tdel
# decd add-genesis-account $(deccli keys show spammer -a --keyring-backend test) 1000000000000000000000000000tdel

# Add initial funds to the genesis file
decd add-genesis-account dx16rr3cvdgj8jsywhx8lfteunn9uz0xg2c7ua9nl 40000011000000000000000000tdel # validator on test-node-fra1-01 (40,000,011 tDEL)
decd add-genesis-account dx1ajytg8jg8ypx0rj9p792x32fuxyezga43jd3ry 40000011000000000000000000tdel # validator on test-node-fra1-02 (40,000,011 tDEL)
decd add-genesis-account dx1azre0dtclv5y05ufynkhswzh0cwh4ktzlas3mp 40000011000000000000000000tdel # validator on test-node-nyc3-01 (40,000,011 tDEL)
decd add-genesis-account dx1j3j2mwxnvlmsu2tkwm4z5390vq8v337wd6hmg2 40000011000000000000000000tdel # validator on test-node-sgp1-01 (40,000,011 tDEL)
decd add-genesis-account dx12k95ukkqzjhkm9d94866r4d9fwx7tsd82r8pjd 160000000000000000000000000tdel # faucet (160,000,000 tDEL)
decd add-genesis-account dx1esffyu0wxk6eez77fhzdxfgvjp4646hqm9sx6c 19999956000000000000000000tdel # tanker (19,999,956 tDEL)

# decd add-genesis-account dxd0c71c31a891e5023ae63fd2bcf2732f04f32158 10000000000000000000000000tdel # validator on test-node-fra1-01 (10,000,000 tDEL)
# decd add-genesis-account dxec88b41e483902678e450f8aa34549e1899123b5 10000000000000000000000000tdel # validator on test-node-fra1-02 (10,000,000 tDEL)
# decd add-genesis-account dxe88797b578fb2847d38924ed7838577e1d7ad962 10000000000000000000000000tdel # validator on test-node-nyc3-01 (10,000,000 tDEL)
# decd add-genesis-account dx9464adb8d367f70e297676ea2a44af600ec8c7ce 10000000000000000000000000tdel # validator on test-node-sgp1-01 (10,000,000 tDEL)
# decd add-genesis-account dx558b4e5ac014af6d95a5a9f5a1d5a54b8de5c1a7 1000000000000000000000000000000tdel # faucet (1,000,000,000,000 tDEL)
# decd add-genesis-account dxcc129271ee35b59c8bde4dc4d3250c906baaeae0 1000000000000000000000000000000tdel # tanker (1,000,000,000,000 tDEL)

# decd add-genesis-account dx1fyqf7gp0gzmpzwxfrah9veaxt8ysl68khxwfjm 1000000000000000000000000tdel # test account (1,000,000 tDEL)
# decd add-genesis-account dx1dvwgj5hc3uqjemk22v9gq7um30v0l6wpcwnwnm 1000000000000000000000000tdel # test account (1,000,000 tDEL)
# decd add-genesis-account dx1cjkq662yycfy03euktm4p5rem0q53f9x89crap 1000000000000000000000000tdel # test account (1,000,000 tDEL)
# decd add-genesis-account dx178h6nvsqg5vr2gq8xv8f9jlmkpg3mjvzktsmj2 1000000000000000000000000tdel # test account (1,000,000 tDEL)

# Add initial signed transactions to the genesis file
    # --commission-rate 0.1 \
decd gentx \
    --name test-node-fra1-01 \
    --sequence 0 \
    --amount 40000000000000000000000000tdel \
    --pubkey dxvalconspub1zcjduepquc5nas24rhqm0l8lyte0dfx2k3uda56wdn998lyrs6mpvsk9xmks6xa0ly \
    --details "Declaring validator on test-node-fra1-01" \
    --website decimalchain.com \
    --node-id bf7a6b366e3c451a3c12b3a6c01af7230fb92fc7 \
    --ip 139.59.133.148 \
    --keyring-backend test
decd collect-gentxs

    # --commission-rate 0.2 \
decd gentx \
    --name test-node-fra1-02 \
    --sequence 0 \
    --amount 40000000000000000000000000tdel \
    --pubkey dxvalconspub1zcjduepq5hj3p750mves8wpmh4ywy6yjkz72sppr2kmzk7lzeedyelauwamsl3c57q \
    --details "Declaring validator on test-node-fra1-02" \
    --website decimalchain.com \
    --node-id c0b9b6c9a0f95e3d2f4aed890806739fc77faefd \
    --ip 64.225.110.228 \
    --keyring-backend test
decd collect-gentxs

    # --commission-rate 0.3 \
decd gentx \
    --name test-node-nyc3-01 \
    --sequence 0 \
    --amount 40000000000000000000000000tdel \
    --pubkey dxvalconspub1zcjduepqjwlm5xcsp60v6fgwt95zq624yjhjnpkrzm209c5f8ajz8rdvq6gsdx5y2l \
    --details "Declaring validator on test-node-nyc3-01" \
    --website decimalchain.com \
    --node-id 76b81a4b817b39d63a3afe1f3a294f2a8f5c55b0 \
    --ip 64.225.56.107 \
    --keyring-backend test
decd collect-gentxs

    # --commission-rate 0.4 \
decd gentx \
    --name test-node-sgp1-01 \
    --sequence 0 \
    --amount 40000000000000000000000000tdel \
    --pubkey dxvalconspub1zcjduepq73se7rmlycftjta3ydjvrjmn28rrweyxyg42tzfr5lcw6lx8zl7qc6dpss \
    --details "Declaring validator on test-node-sgp1-01" \
    --website decimalchain.com \
    --node-id 29e566c41d51be90fa53340ba4edccefbebe8cb2 \
    --ip 139.59.192.48 \
    --keyring-backend test
decd collect-gentxs

# Configure created blockchain
deccli config chain-id decimal-testnet-06-17-14-00
deccli config trust-node true

# Finish new blockchain initialization
decd collect-gentxs
