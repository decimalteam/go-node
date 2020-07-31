package cli

const mainNetGenesis = `
{
  "genesis_time": "2020-07-31T11:00:00.000000Z",
  "chain_id": "decimal-mainnet-07-30",
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
                    "website": "about:blank",
                    "security_contact": "email: empty",
                    "details": "The test validator ever"
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
                "signature": "ZxQ8Tbf+ZH+5j8c7NNYGTAK4x0ytlpnPqIlSV6CDLjMNYFjRwD7kDUiv72XTx2S6BwVxCm+mc8JAH2aKIQ+vbA=="
              }
            ],
            "memo": "e221e380920b9b42bba268b6b644eb6ae81587af@81.89.56.52:26656"
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
                    "moniker": "Tank",
                    "identity": "",
                    "website": "about:blank",
                    "security_contact": "email: empty",
                    "details": "The test validator ever"
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
                "signature": "1Emi6b3DjZ55C8S3yVVkVsqw2PMrq6EGVzKG29GgcugBFFGHxC6b6GDYoHPVKI5twQUzxZe0qN1W3wI50miXPw=="
              }
            ],
            "memo": "f67f2ae52e8687eb9816a1dd246261118e869336@171.25.221.204:26656"
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
                    "moniker": "Main",
                    "identity": "",
                    "website": "about:blank",
                    "security_contact": "email: empty",
                    "details": "The test validator ever"
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
                "signature": "v+69t3fFTrfe7B0b3C5ucy3uY3Kau9AZGPlxtorNdaFA2gLcmlCy6Hzky40YwGEs/GkRIzBQXCQuuw3tQljNag=="
              }
            ],
            "memo": "c09aee9eb6e1ac84e22cab9b53d57b7898755980@171.25.221.205:26656"
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
                    "moniker": "crypton",
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
                "signature": "YtIiTcPhBL5EZTxCoIUZBtFsukIHPw1P/uNlsFuG3wYLycYP++1O5vIh/A8DxvZOlkSNdfJMVTkCFaJCoKb30A=="
              }
            ],
            "memo": "a1ca88f09330204fd8d96e7a55da0964399cf6de@135.181.5.158:26656"
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
            "address": "dx1hh2mqw4xan5n52df0lprcj5nw3hwr95aqvw3mp",
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
    "check": {}
  }
}
`
