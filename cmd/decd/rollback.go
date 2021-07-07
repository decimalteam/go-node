package main

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"
	"strconv"
	"strings"
	"time"
)

func fixAppHashError(ctx *server.Context, defaultNodeHome string) *cobra.Command {
	const flagSetStateHeight = "state-height"

	cmd := &cobra.Command{
		Use:   "rollback [100]",
		Short: "",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg := config.DefaultConfig()
			cfg.SetRoot(viper.GetString(cli.HomeFlag))

			stateDB, err := node.DefaultDBProvider(&node.DBContext{ID: "state", Config: cfg})
			if err != nil {
				return err
			}

			st := state.LoadState(stateDB)

			blockStoreDB, err := node.DefaultDBProvider(&node.DBContext{ID: "blockstore", Config: cfg})
			if err != nil {
				return err
			}

			blockStore := store.NewBlockStore(blockStoreDB)

			stateHeightFlag := viper.GetString(flagSetStateHeight)
			if stateHeightFlag != "" {
				stateHeight, err := strconv.ParseInt(stateHeightFlag, 10, 64)
				if err != nil {
					return err
				}

				st.LastBlockHeight = stateHeight

				state.SaveState(stateDB, st)
				return nil
			}

			var countBlocks int64
			if len(args) == 1 {
				countBlocks, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return err
				}
			}
			height := st.LastBlockHeight

			if countBlocks > (st.LastBlockHeight - st.LastBlockHeight/100*100) {
				countBlocks = st.LastBlockHeight - st.LastBlockHeight/100*100
			}

			for i := int64(0); i < countBlocks; i++ {
				block := blockStore.LoadBlock(height)
				if block == nil {
					break
				}

				err = DeleteBlock(blockStoreDB, blockStore, block)
				if err != nil {
					return err
				}

				height--
			}

			block := blockStore.LoadBlock(height)

			st.LastBlockHeight = height - 1
			st.LastBlockID = block.LastBlockID
			st.AppHash = block.AppHash
			st.LastResultsHash = block.LastResultsHash
			st.LastBlockTime = time.Unix(0, block.Time.UnixNano()-time.Second.Nanoseconds()*5)
			st.LastHeightValidatorsChanged = st.LastBlockHeight
			fmt.Println(strings.ToUpper(hex.EncodeToString(st.Validators.Hash())))
			for i, validator := range st.Validators.Validators {
				if validator.Address.String() == "BA1B262312BBDF500C5410F26CA80AD63CFC3F81" {
					fmt.Println("done")
					validatorSet := st.Validators.Copy()
					(*validatorSet.Validators[i]).VotingPower = 4568124
					st.Validators = validatorSet
					validatorSet = st.Validators.Copy()
					(*validatorSet.Validators[i]).VotingPower = 4568124
					st.NextValidators = validatorSet
				}
			}

			for _, validator := range st.Validators.Validators {
				if validator.Address.String() == "BA1B262312BBDF500C5410F26CA80AD63CFC3F81" {
					fmt.Println(validator.VotingPower)
				}
			}
			for _, validator := range st.NextValidators.Validators {
				if validator.Address.String() == "BA1B262312BBDF500C5410F26CA80AD63CFC3F81" {
					fmt.Println(validator.VotingPower)
				}
			}
			fmt.Println(strings.ToUpper(hex.EncodeToString(st.Validators.Hash())))
			fmt.Println(strings.ToUpper(hex.EncodeToString(st.NextValidators.Hash())))
			fmt.Println(block.ValidatorsHash.String())
			fmt.Println(block.NextValidatorsHash.String())

			state.SaveState(stateDB, st)
			return nil
		},
	}

	FsSetStateHeight := flag.NewFlagSet("", flag.ContinueOnError)
	FsSetStateHeight.String(flagSetStateHeight, "", "Set state height")

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().AddFlagSet(FsSetStateHeight)

	return cmd
}

func DeleteBlock(db db.DB, blockStore *store.BlockStore, block *types.Block) error {
	meta := blockStore.LoadBlockMeta(block.Height)

	err := db.Delete(calcBlockMetaKey(block.Height))
	if err != nil {
		return err
	}

	err = db.Delete(calcBlockHashKey(block.Hash()))
	if err != nil {
		return err
	}

	for i := 0; i < meta.BlockID.PartsHeader.Total; i++ {
		err = db.Delete(calcBlockPartKey(block.Height, i))
		if err != nil {
			return err
		}
	}

	err = db.Delete(calcBlockCommitKey(block.Height - 1))
	if err != nil {
		return err
	}

	err = db.Delete(calcSeenCommitKey(block.Height))
	if err != nil {
		return err
	}

	store.BlockStoreStateJSON{Height: block.Height - 1}.Save(db)

	err = db.SetSync(nil, nil)
	if err != nil {
		return err
	}

	return nil
}

//-----------------------------------------------------------------------------

func calcBlockMetaKey(height int64) []byte {
	return []byte(fmt.Sprintf("H:%v", height))
}

func calcBlockPartKey(height int64, partIndex int) []byte {
	return []byte(fmt.Sprintf("P:%v:%v", height, partIndex))
}

func calcBlockCommitKey(height int64) []byte {
	return []byte(fmt.Sprintf("C:%v", height))
}

func calcSeenCommitKey(height int64) []byte {
	return []byte(fmt.Sprintf("SC:%v", height))
}

func calcBlockHashKey(hash []byte) []byte {
	return []byte(fmt.Sprintf("BH:%x", hash))
}
