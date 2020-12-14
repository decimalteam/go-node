package cli

const MainNetGenesis = `
{
  "genesis_time": "2020-08-01T15:00:00.000000Z",
  "chain_id": "decimal-mainnet-08-01",
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
      "pub_key_types": [
        "ed25519"
      ]
    }
  },
  "app_hash": "",
  "app_state": {
    "validator": {
      "params": {
        "unbonding_time": "2592000000000000",
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
                  "validator_addr": "dxvaloper1q7uy7u0ewgsgwgzvvmmkqddy6nujj60rm6dr3y",
                  "reward_addr": "dx1q7uy7u0ewgsgwgzvvmmkqddy6nujj60r8g2wy0",
                  "pub_key": {
                    "type": "tendermint/PubKeyEd25519",
                    "value": "HpWT0IlxRZ6r7NwqS2X36EwTrS5x/iLOCmaxqbpPpIs="
                  },
                  "stake": {
                    "denom": "del",
                    "amount": "100000000000000000000000"
                  },
                  "description": {
                    "moniker": "BitTeam",
                    "identity": "",
                    "website": "https://bit.team/",
                    "security_contact": "email: support@bit.team",
                    "details": ""
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
                  "value": "A9STcqOqpRAtz68XADSe3tdFc064LHeruMcNwJ4oizc1"
                },
                "signature": "RyfC9trGP54bUNDsgfYK47Wd0QjCe/z7IRLevpIgEd1qvUooMgCFCbhufF6iGzno26f8rAa1M/he0LWqKBX/fQ=="
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
                  "commission": "0.050000000000000000",
                  "validator_addr": "dxvaloper1eu9juhlsa4svhdhh4w2mknqtarnemvf3x93qll",
                  "reward_addr": "dx1eu9juhlsa4svhdhh4w2mknqtarnemvf36hkd25",
                  "pub_key": {
                    "type": "tendermint/PubKeyEd25519",
                    "value": "T7qeKz/Sl60oLt0h0OC8890tPONZD3dCnpPWjuqIPK0="
                  },
                  "stake": {
                    "denom": "del",
                    "amount": "100000000000000000000000"
                  },
                  "description": {
                    "moniker": "Main Node",
                    "identity": "",
                    "website": "about:blank",
                    "security_contact": "",
                    "details": ""
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
                  "value": "A8DlwOt+vGis0Asv2BygA+RVpNfBTHVsF2kQSjE9/7dg"
                },
                "signature": "K4t8T2w4fDuZe88nO3mo35Np4l/1xgB5bawqpH+3p8Mr926lxB29Igt3VKbl8NapsEVvdgiMelj1lrnaCTcoTQ=="
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
                  "commission": "0.080000000000000000",
                  "validator_addr": "dxvaloper1r9hdv2p89yyekxcmvgpud6wgg2xpln0nz3rf0q",
                  "reward_addr": "dx1r9hdv2p89yyekxcmvgpud6wgg2xpln0n7ryy6t",
                  "pub_key": {
                    "type": "tendermint/PubKeyEd25519",
                    "value": "Eg73El48hqIlLe0mYnqpig5Ua6enMPkgeyALhe/Ixzg="
                  },
                  "stake": {
                    "denom": "del",
                    "amount": "100000000000000000000000"
                  },
                  "description": {
                    "moniker": "Turing",
                    "identity": "",
                    "website": "about:blank",
                    "security_contact": "",
                    "details": ""
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
                  "value": "A3CznhbVm2bV7buF8gbBIr5UMa2S9DZTUwP5mYg1IrsA"
                },
                "signature": "at79Fj7ah12tLxAEKZd9MUwyx3EvHW/uniBnCayruVondcy3RhT2JD/PjZdsCzwYj74iT/fV6kGnO1C1uqdRbQ=="
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
                  "validator_addr": "dxvaloper1naeup2d7gc30tw0wqxgt4c5yau0atmxvhdjpkr",
                  "reward_addr": "dx1naeup2d7gc30tw0wqxgt4c5yau0atmxvtl4vrg",
                  "pub_key": {
                    "type": "tendermint/PubKeyEd25519",
                    "value": "W7m1XQHmVGGSx/H+eYVbJk0WbtPL1h53CndXk32dhA8="
                  },
                  "stake": {
                    "denom": "del",
                    "amount": "100000000000000000000000"
                  },
                  "description": {
                    "moniker": "Crypton",
                    "identity": "",
                    "website": "https://crypton.studio",
                    "security_contact": "email: security@crypton.studio",
                    "details": "Crypton Validator Node - 01"
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
                  "value": "AkwYGJxv6TvP/VQchZCLopLy5xN6/oLgyQI8ZPkBsA6B"
                },
                "signature": "ocohyP9fPK2FlGdS6f/vKdC59WErcsqA62iHlcLgy3JbncID7GUD1qAp97/9VOUihMAgyAkdoI7G4XtA2L8rzA=="
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
            "address": "dx15ljaa3yyqg5mhlf4q8vsdamudfthanc9gj7dr6",
            "coins": [
              {
                "denom": "del",
                "amount": "39900000000000000000000000"
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
            "address": "dx1wmhaeu3wtzlqj4lk79caceacqmz30evq88dw69",
            "coins": [
              {
                "denom": "del",
                "amount": "39900000000000000000000000"
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
            "address": "dx1edcq8td9xc95eaa5l00k0yh03dwjdn92vh399r",
            "coins": [
              {
                "denom": "del",
                "amount": "39900000000000000000000000"
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
            "address": "dx1l5e2y798txvp8z4f0g8etp6rpyrzttew3es5g8",
            "coins": [
              {
                "denom": "del",
                "amount": "39900000000000000000000000"
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
            "address": "dx1ttk9a4messsw66g0hsm97xztyhl60vu0n8t6ft",
            "coins": [
              {
                "denom": "del",
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
            "address": "dx1q7uy7u0ewgsgwgzvvmmkqddy6nujj60r8g2wy0",
            "coins": [
              {
                "denom": "del",
                "amount": "100000000000000000000000"
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
            "address": "dx1eu9juhlsa4svhdhh4w2mknqtarnemvf36hkd25",
            "coins": [
              {
                "denom": "del",
                "amount": "100000000000000000000000"
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
            "address": "dx1r9hdv2p89yyekxcmvgpud6wgg2xpln0n7ryy6t",
            "coins": [
              {
                "denom": "del",
                "amount": "100000000000000000000000"
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
            "address": "dx1naeup2d7gc30tw0wqxgt4c5yau0atmxvtl4vrg",
            "coins": [
              {
                "denom": "del",
                "amount": "100000000000000000000000"
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
      "title": "Decimal coin",
      "symbol": "del",
      "initial_volume": "200000000000000000000000000"
    },
    "check": {},
	"gov": {
      "starting_proposal_id": "1",
      "votes": null,
      "proposals": null,
      "tally_params": {
        "quorum": "0.667000000000000000",
        "threshold": "0"
      }
    },
    "swap": {
      "params":
      {
        "locked_time_in": "43200000000000",
        "locked_time_out": "86400000000000"
      },
      "swaps": null
    }
  }
}
`
