package cli

const testNetGenesis = `
{
  "genesis_time": "2020-07-31T15:30:00.000000Z",
  "chain_id": "decimal-testnet-07-31-22-30",
  "consensus_params": {
    "block": {
      "max_bytes": "10000000",
      "max_gas": "-1",
      "time_iota_ms": "1000"
    },
    "evidence": {
      "max_age_num_blocks": "100000",
      "max_age_duration": "86400000000000"
    },
    "validator": {
      "pub_key_types": ["ed25519"]
    }
  },
  "app_hash": "",
  "app_state": {
    "validator": {
      "params": {
        "unbonding_time": "3600000000000",
        "max_validators": 256,
        "max_entries": 7,
        "bond_denom": "tdel",
        "historical_entries": 0,
        "max_delegations": 1000
      },
      "last_total_power": "0",
      "last_validator_powers": null,
      "validators": null,
      "delegations": null,
      "unbonding_delegations": null,
      "exported": false
    },
    "genutil": {
      "gentxs": [
        {
          "type": "cosmos-sdk/StdTx",
          "value": {
            "msg": [
              {
                "type": "validator/declare_candidate",
                "value": {
                  "commission": "0.100000000000000000",
                  "validator_addr": "dxvaloper16rr3cvdgj8jsywhx8lfteunn9uz0xg2czw6gx5",
                  "reward_addr": "dx16rr3cvdgj8jsywhx8lfteunn9uz0xg2c7ua9nl",
                  "pub_key": {
                    "type": "tendermint/PubKeyEd25519",
                    "value": "5ik+wVUdwbf8/yLy9qTKtHje005sylP8g4a2FkLFNu0="
                  },
                  "stake": {
                    "denom": "tdel",
                    "amount": "40000000000000000000000000"
                  },
                  "description": {
                    "moniker": "test-node-fra1-01",
                    "identity": "",
                    "website": "decimalchain.com",
                    "security_contact": "",
                    "details": "Declaring validator on test-node-fra1-01"
                  }
                }
              }
            ],
            "fee": {
              "amount": [],
              "gas": "200000"
            },
            "signatures": [
              {
                "pub_key": {
                  "type": "tendermint/PubKeySecp256k1",
                  "value": "Ann34HTQiiPi/Ht/2eSaDWVwoov2ycuYjpL2eMMpNQl0"
                },
                "signature": "Kc9WsuDL9j4SRC0EyIXITnZXRixYyh7f9E0ocA3giDcURsKS/0PwcYm6wJiiJkSUJIPohDyi3oiOofTqEX3t/g=="
              }
            ],
            "memo": ""
          }
        },
        {
          "type": "cosmos-sdk/StdTx",
          "value": {
            "msg": [
              {
                "type": "validator/declare_candidate",
                "value": {
                  "commission": "0.100000000000000000",
                  "validator_addr": "dxvaloper1ajytg8jg8ypx0rj9p792x32fuxyezga4dq2uk0",
                  "reward_addr": "dx1ajytg8jg8ypx0rj9p792x32fuxyezga43jd3ry",
                  "pub_key": {
                    "type": "tendermint/PubKeyEd25519",
                    "value": "peUQ+o/bMwO4O71I4miSsLyoBCNVtit74s5aTP+8d3c="
                  },
                  "stake": {
                    "denom": "tdel",
                    "amount": "40000000000000000000000000"
                  },
                  "description": {
                    "moniker": "test-node-fra1-02",
                    "identity": "",
                    "website": "decimalchain.com",
                    "security_contact": "",
                    "details": "Declaring validator on test-node-fra1-02"
                  }
                }
              }
            ],
            "fee": {
              "amount": [],
              "gas": "200000"
            },
            "signatures": [
              {
                "pub_key": {
                  "type": "tendermint/PubKeySecp256k1",
                  "value": "A+0Sm90CYEkcknXu/cvYx9eKpau17Yyd54mzDdtSG9bZ"
                },
                "signature": "4FSytBvKZcnUF9M4uCZICtmdPdX6Rr9UTXU0lCg82Jt6SvDwV1bNMo+sXCh79+uicJ2GBUxgAMPlHm7J/x1bkg=="
              }
            ],
            "memo": ""
          }
        },
        {
          "type": "cosmos-sdk/StdTx",
          "value": {
            "msg": [
              {
                "type": "validator/declare_candidate",
                "value": {
                  "commission": "0.100000000000000000",
                  "validator_addr": "dxvaloper1azre0dtclv5y05ufynkhswzh0cwh4ktzr0huw2",
                  "reward_addr": "dx1azre0dtclv5y05ufynkhswzh0cwh4ktzlas3mp",
                  "pub_key": {
                    "type": "tendermint/PubKeyEd25519",
                    "value": "k7+6GxAOns0lDlloIGlVJK8phsMW1PLiiT9kI42sBpE="
                  },
                  "stake": {
                    "denom": "tdel",
                    "amount": "40000000000000000000000000"
                  },
                  "description": {
                    "moniker": "test-node-nyc3-01",
                    "identity": "",
                    "website": "decimalchain.com",
                    "security_contact": "",
                    "details": "Declaring validator on test-node-nyc3-01"
                  }
                }
              }
            ],
            "fee": {
              "amount": [],
              "gas": "200000"
            },
            "signatures": [
              {
                "pub_key": {
                  "type": "tendermint/PubKeySecp256k1",
                  "value": "Ax72b3e3Tz8Wx7Iq9FXaM1sbTit+1AujjKHRsTafhrtE"
                },
                "signature": "xAQxy6Bm1CCOdQzg49HA7xSpohweytWDwDrHFLhmbUMns2dRqiIbp99G0Wy1TCtiVjg50HQkRynAyGB73Y3+Ng=="
              }
            ],
            "memo": ""
          }
        },
        {
          "type": "cosmos-sdk/StdTx",
          "value": {
            "msg": [
              {
                "type": "validator/declare_candidate",
                "value": {
                  "commission": "0.100000000000000000",
                  "validator_addr": "dxvaloper1j3j2mwxnvlmsu2tkwm4z5390vq8v337w3gskap",
                  "reward_addr": "dx1j3j2mwxnvlmsu2tkwm4z5390vq8v337wd6hmg2",
                  "pub_key": {
                    "type": "tendermint/PubKeyEd25519",
                    "value": "9GGfD38mErkvsSNkwctzUcY3ZIYiKqWJI6fw7XzHF/w="
                  },
                  "stake": {
                    "denom": "tdel",
                    "amount": "40000000000000000000000000"
                  },
                  "description": {
                    "moniker": "test-node-sgp1-01",
                    "identity": "",
                    "website": "decimalchain.com",
                    "security_contact": "",
                    "details": "Declaring validator on test-node-sgp1-01"
                  }
                }
              }
            ],
            "fee": {
              "amount": [],
              "gas": "200000"
            },
            "signatures": [
              {
                "pub_key": {
                  "type": "tendermint/PubKeySecp256k1",
                  "value": "AtKn3ANbRqsIg8zpF0/03t9kvEuUtd7ZS9VnPu/8zF5z"
                },
                "signature": "aMSlXt5csFVIIFUBUyX2oeXyAZiTtetpx46t+AI7pgRxAwEoXWFKlcCxmkUaNRq62mgTvL1/lDko7w3c42p6UQ=="
              }
            ],
            "memo": ""
          }
        }
      ]
    },
    "auth": {
      "params": {
        "max_memo_characters": "256",
        "tx_sig_limit": "7",
        "tx_size_cost_per_byte": "10",
        "sig_verify_cost_ed25519": "590",
        "sig_verify_cost_secp256k1": "1000"
      },
      "accounts": [
        {
          "type": "cosmos-sdk/Account",
          "value": {
            "address": "dx16rr3cvdgj8jsywhx8lfteunn9uz0xg2c7ua9nl",
            "coins": [
              {
                "denom": "tdel",
                "amount": "40000000000000000000000000"
              }
            ],
            "public_key": "",
            "account_number": 0,
            "sequence": 0
          }
        },
        {
          "type": "cosmos-sdk/Account",
          "value": {
            "address": "dx1ajytg8jg8ypx0rj9p792x32fuxyezga43jd3ry",
            "coins": [
              {
                "denom": "tdel",
                "amount": "40000000000000000000000000"
              }
            ],
            "public_key": "",
            "account_number": 0,
            "sequence": 0
          }
        },
        {
          "type": "cosmos-sdk/Account",
          "value": {
            "address": "dx1azre0dtclv5y05ufynkhswzh0cwh4ktzlas3mp",
            "coins": [
              {
                "denom": "tdel",
                "amount": "40000000000000000000000000"
              }
            ],
            "public_key": "",
            "account_number": 0,
            "sequence": 0
          }
        },
        {
          "type": "cosmos-sdk/Account",
          "value": {
            "address": "dx1j3j2mwxnvlmsu2tkwm4z5390vq8v337wd6hmg2",
            "coins": [
              {
                "denom": "tdel",
                "amount": "40000000000000000000000000"
              }
            ],
            "public_key": "",
            "account_number": 0,
            "sequence": 0
          }
        },
        {
          "type": "cosmos-sdk/Account",
          "value": {
            "address": "dx1twjeqjvqcznu2uqagms55e4a8rtakapddkgcsm",
            "coins": [
              {
                "denom": "tdel",
                "amount": "160000000000000000000000000"
              }
            ],
            "public_key": "",
            "account_number": 0,
            "sequence": 0
          }
        },
        {
          "type": "cosmos-sdk/Account",
          "value": {
            "address": "dx1a0329z2gdh98mn4m6uzssan82t7vf03nv7t38g",
            "coins": [
              {
                "denom": "tdel",
                "amount": "20000000000000000000000000"
              }
            ],
            "public_key": "",
            "account_number": 0,
            "sequence": 0
          }
        }
      ]
    },
    "bank": {
      "send_enabled": true
    },
    "params": null,
    "supply": {
      "supply": []
    },
    "coin": {
      "title": "Test decimal coin",
      "symbol": "tdel",
      "initial_volume": "340000000000000000000000000"
    },
    "check": {}
  }
}
`
