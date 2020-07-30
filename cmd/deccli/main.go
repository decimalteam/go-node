package main

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"

	"bitbucket.org/decimalteam/go-node/app"
	"bitbucket.org/decimalteam/go-node/config"
	coinCmd "bitbucket.org/decimalteam/go-node/x/coin/client/cli"
	multisigCmd "bitbucket.org/decimalteam/go-node/x/multisig/client/cli"
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := app.MakeCodec()

	// Read in the configuration file for the sdk
	_config := sdk.GetConfig()
	_config.SetCoinType(60)
	_config.SetFullFundraiserPath("44'/60'/0'/0/0")
	_config.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	_config.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	_config.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)
	_config.Seal()

	rootCmd := &cobra.Command{
		Use:   "deccli",
		Short: "Decimal Client Console",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of decimal node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		flags.LineBreak,
		lcd.ServeCommand(cdc, registerRoutes),
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
		version.Cmd,
		flags.NewCompletionCmd(rootCmd, true),
		getGenesis(),
	)

	executor := cli.PrepareMainCmd(rootCmd, "AU", app.DefaultCLIHome)
	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	app.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		flags.LineBreak,
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		flags.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(coinCmd.GetQueryCmd("coin", cdc), cdc)
	app.ModuleBasics.AddQueryCommands(multisigCmd.GetQueryCmd("multisig", cdc), cdc)
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		bankcmd.SendTxCmd(cdc),
		flags.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		flags.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(coinCmd.GetTxCmd(cdc), cdc)
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	return txCmd
}

func initConfig(cmd *cobra.Command) error {
	home, err := cmd.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	cfgFile := path.Join(home, "config", "config.toml")
	if _, err := os.Stat(cfgFile); err == nil {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(flags.FlagChainID, cmd.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, cmd.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, cmd.PersistentFlags().Lookup(cli.OutputFlag))
}

func getGenesis() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-genesis [mainnet|testnet]",
		Short: "Get genesis",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			network := args[0]

			switch network {
			case "mainnet":
				fmt.Print(mainNetGenesis)
			case "testnet":
				fmt.Print(testNetGenesis)
			case "devnet":
				fmt.Print(devNetGenesis)
			default:
				return fmt.Errorf("invalid network")
			}
			return nil
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

const devNetGenesis = `
{
  "genesis_time": "2020-07-22T11:50:00.000000Z",
  "chain_id": "decimal-devnet-07-23-20-55",
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
                    "moniker": "dev-node-fra1-01",
                    "identity": "",
                    "website": "decimalchain.com",
                    "security_contact": "",
                    "details": "Declaring validator on dev-node-fra1-01"
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
                "signature": "Bp7Ergky/TonbvNyJo26Tfp+VhZLu93ZCI3UrXPI5HBbkzM6eYTxAm7/89pT8BigO72bduE9Skl+xRCe7kuOuA=="
              }
            ],
            "memo": "8a2cc38f5264e9699abb8db91c9b4a4a061f000d@46.101.127.241:26656"
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
                    "moniker": "dev-node-fra1-02",
                    "identity": "",
                    "website": "decimalchain.com",
                    "security_contact": "",
                    "details": "Declaring validator on dev-node-fra1-02"
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
                "signature": "bLd92//oiK0B2cwCRr9XTuymyqN/orxzyKSW4yDVm3khGla3TKM8Ic/p+o2iZVS1VEW9ENrWbUGM0Z2kr5R/fA=="
              }
            ],
            "memo": "e0e7a88de0b39bd2adceb3516d353582ff94ec15@164.90.211.234:26656"
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
                    "moniker": "dev-node-tor1-01",
                    "identity": "",
                    "website": "decimalchain.com",
                    "security_contact": "",
                    "details": "Declaring validator on dev-node-tor1-01"
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
                "signature": "g4L04c6+v2SmtpFL8aP0Yz9y2RVpq5ka1E7Vkvt7+SRGJ3FTjzfmcEW2PU5gKoSLv70JhCTxBLCvQv+6MUh/jg=="
              }
            ],
            "memo": "27fcfef145b3717c5d639ec72fb12f9c43da98f0@167.99.182.218:26656"
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
            "public_key": "",
            "account_number": 0,
            "sequence": 0
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
            "public_key": "",
            "account_number": 0,
            "sequence": 0
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
            "public_key": "",
            "account_number": 0,
            "sequence": 0
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
            "public_key": "",
            "account_number": 0,
            "sequence": 0
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
      "initial_volume": "340000000000000000000000000"
    },
    "check": {}
  }
}
`

const testNetGenesis = `
{
  "genesis_time": "2020-07-28T11:30:00.000000Z",
  "chain_id": "decimal-testnet-07-28-18-30",
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
                "signature": "jtew4RO2QP8GU4xs89kIk/b2CY9X/SuzEEOWWdhtQ9QpZBqcK9/W/OIjoygf4Qv5iGG5HTodMRUegdO0thmAgA=="
              }
            ],
            "memo": "bf7a6b366e3c451a3c12b3a6c01af7230fb92fc7@139.59.133.148:26656"
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
                "signature": "w9nwHWSIloqIXQg54gzOycgvgkBX/Y1P3QwVCtuN2T0Zh1ObKcq3wRUWC5K52u6PdhdUm9i0HKy0pPL42SPVlw=="
              }
            ],
            "memo": "c0b9b6c9a0f95e3d2f4aed890806739fc77faefd@64.225.110.228:26656"
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
                "signature": "egnSsqoPqYUp3O/IRIRhYzaPqyFCrt2hnl4gxd4gEYdxyBXUFFEHvA6yd/SZhTAuEVMmiOE+32Cddy3cQk5E0w=="
              }
            ],
            "memo": "76b81a4b817b39d63a3afe1f3a294f2a8f5c55b0@64.225.56.107:26656"
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
                "signature": "oLYlwAdAX3avgwdlyr+gZfXvev7JTrABbVVj3ttjHtF6i3iiV01rbdmmM2bOQXC/jsw6j4SNFJIKEBhY3WdpGg=="
              }
            ],
            "memo": "29e566c41d51be90fa53340ba4edccefbebe8cb2@139.59.192.48:26656"
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

const mainNetGenesis = `
`
