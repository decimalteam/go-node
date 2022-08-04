package cli

const devNetGenesis = `
{
  "genesis_time": "2022-08-04T00:00:00.000000Z",
  "chain_id": "decimal-devnet-08-04-01-00",
  "consensus_params": {
    "block": {
      "max_bytes": "10000000",
      "max_gas": "100000",
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
        "unbonding_time": "1800000000000",
        "max_validators": 256,
        "max_entries": 7,
        "bond_denom": "del",
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
                  "validator_addr": "dxvaloper1lx4lvt8sjuxj8vw5dcf6knnq0pacre4wx926l8",
                  "reward_addr": "dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v",
                  "pub_key": {
                    "type": "tendermint/PubKeyEd25519",
                    "value": "c4YgZK1GlgVB4YFh94Qcqo0SBBIhLJ7ncIwXO4URf1A="
                  },
                  "stake": {
                    "denom": "del",
                    "amount": "40000000000000000000000000"
                  },
                  "description": {
                    "moniker": "dev-node-validator-01",
                    "identity": "",
                    "website": "decimalchain.com",
                    "security_contact": "",
                    "details": "Declaring validator on dev-node-validator-01"
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
                  "value": "AyXteacATBJVsNGhTGUwgHfLO4mJTXhvK/H/2MLOTylo"
                },
                "signature": "r/ddEVIU89OcvB+UAfZ1yWpriBY5IojH0RziGNRbltVD6i9odImTxfUqg0iYFEx7PcbU3FHjUisYHtWndDyfwQ=="
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
                  "validator_addr": "dxvaloper1mvqrrrlcd0gdt256jxg7n68e4neppu5tk872z3",
                  "reward_addr": "dx1mvqrrrlcd0gdt256jxg7n68e4neppu5t24e8h6",
                  "pub_key": {
                    "type": "tendermint/PubKeyEd25519",
                    "value": "5duJ6p45rx8gO3TOVnwHxyr1I3SbJ2z+EbBmIvYowQw="
                  },
                  "stake": {
                    "denom": "del",
                    "amount": "40000000000000000000000000"
                  },
                  "description": {
                    "moniker": "dev-node-validator-02",
                    "identity": "",
                    "website": "decimalchain.com",
                    "security_contact": "",
                    "details": "Declaring validator on dev-node-validator-02"
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
                  "value": "A8s6NPtmz3ywH2LGxNnRjjsEIdt53ZuRQy7q6Mof0iA9"
                },
                "signature": "4OPTyRqX3007I1OwjP7BuSm/PI777/k5yVsx4TFKO4JYY6Kfc6DCtTIaUFKo+KeS/QHQ6OgYt2ZmKSHRY2nTXQ=="
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
                  "validator_addr": "dxvaloper1nrr6er27mmcufmaqm4dyu6c5r6489cfmdxucuq",
                  "reward_addr": "dx1nrr6er27mmcufmaqm4dyu6c5r6489cfm35m4ft",
                  "pub_key": {
                    "type": "tendermint/PubKeyEd25519",
                    "value": "x9WI+in8vh2/6eLpivWgqZ50waXiQb0mAI50GqQI8bI="
                  },
                  "stake": {
                    "denom": "del",
                    "amount": "40000000000000000000000000"
                  },
                  "description": {
                    "moniker": "dev-node-validator-03",
                    "identity": "",
                    "website": "decimalchain.com",
                    "security_contact": "",
                    "details": "Declaring validator on dev-node-validator-03"
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
                  "value": "A/jGezwfhOzAyTaRbg3off9HYzvtUO4CxaRrBJzdHGlw"
                },
                "signature": "h4t40XlgArqrKoYGyX02Y+rIKFVYKkUC97YicXUhzu8J1z2vlf2dm6xURKOR1l0EH0N2QupEK0aXimXWFUcQDQ=="
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
            "address": "dx1lx4lvt8sjuxj8vw5dcf6knnq0pacre4w6hdh2v",
            "coins": [
              {
                "denom": "del",
                "amount": "40000000000000000000000000"
              }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
          }
        },
        {
          "type": "cosmos-sdk/Account",
          "value": {
            "address": "dx1mvqrrrlcd0gdt256jxg7n68e4neppu5t24e8h6",
            "coins": [
              {
                "denom": "del",
                "amount": "40000000000000000000000000"
              }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
          }
        },
        {
          "type": "cosmos-sdk/Account",
          "value": {
            "address": "dx1nrr6er27mmcufmaqm4dyu6c5r6489cfm35m4ft",
            "coins": [
              {
                "denom": "del",
                "amount": "40000000000000000000000000"
              }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
          }
        },
        {
          "type": "cosmos-sdk/Account",
          "value": {
            "address": "dx1tvqxh4x7pedyqpzqp9tdf068k4q9j2hm3lmghl",
            "coins": [
              {
                "denom": "del",
                "amount": "200000000000000000000000000"
              }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
          }
        },
        {
          "type": "cosmos-sdk/Account",
          "value": {
            "address": "dx1mtlnpmwf8zr6pek6gq25nv45x2890sne2ap0cc",
            "coins": [
              {
                "denom": "del",
                "amount": "20000000000000000000000000"
              }
            ],
            "public_key": null,
            "account_number": "0",
            "sequence": "0"
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
      "title": "Decimal coin",
      "symbol": "del",
      "initial_volume": "340000000000000000000000000"
    },
    "check": {},
    "multisig": {},
    "gov": {
      "starting_proposal_id": "1",
      "votes": null,
      "proposals": null,
      "tally_params": {
        "quorum": "0.334000000000000000",
        "threshold": "0.500000000000000000"
      }
    },
    "swap": {
      "params":
        {
          "locked_time_in": "43200000000000",
          "locked_time_out": "86400000000000"
        },
      "swaps": null
    },
    "nft": {
      "owners": [],
      "collections": {}
    }
  }
}
`
