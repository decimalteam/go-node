package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/config"
	cryptoamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/libs/cli"
	tmos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"
	"strconv"
	"time"
)

func fixAppHashError(ctx *server.Context, defaultNodeHome string) *cobra.Command {
	const flagSetStateHeight = "state-height"

	cmd := &cobra.Command{
		Use:   "rollback [100]",
		Short: "",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cryptoamino.RegisterAmino(cdc)

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

			valInfo := loadValidatorsInfo(stateDB, height)
			if valInfo.ValidatorSet == nil {
				valInfo = loadValidatorsInfo(stateDB, valInfo.LastHeightChanged)
				if valInfo.ValidatorSet == nil {
					panic(valInfo)
				}
			}

			block := blockStore.LoadBlock(height)

			st.LastBlockHeight = height - 1
			st.LastBlockID = block.LastBlockID
			st.AppHash = block.AppHash
			st.LastResultsHash = block.LastResultsHash
			st.LastBlockTime = time.Unix(0, block.Time.UnixNano()-time.Second.Nanoseconds()*5)
			st.LastHeightValidatorsChanged = valInfo.LastHeightChanged
			st.LastValidators = valInfo.ValidatorSet
			st.Validators = valInfo.ValidatorSet
			st.NextValidators = valInfo.ValidatorSet

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

var cdc = amino.NewCodec()

func loadValidatorsInfo(db db.DB, height int64) *state.ValidatorsInfo {
	buf, err := db.Get(calcValidatorsKey(height))
	if err != nil {
		panic(err)
	}
	if len(buf) == 0 {
		return nil
	}

	v := new(state.ValidatorsInfo)
	err = cdc.UnmarshalBinaryBare(buf, v)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		tmos.Exit(fmt.Sprintf(`LoadValidators: Data has been corrupted or its spec has changed:
                %v\n`, err))
	}
	// TODO: ensure that buf is completely read.

	return v
}

func calcValidatorsKey(height int64) []byte {
	return []byte(fmt.Sprintf("validatorsKey:%v", height))
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
