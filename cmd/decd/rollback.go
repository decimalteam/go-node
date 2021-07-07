package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/node"
	protostore "github.com/tendermint/tendermint/proto/tendermint/store"
	"github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/types"
	db "github.com/tendermint/tm-db"
)

func fixAppHashError(ctx *server.Context, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback [100]",
		Short: "",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			cfg := config.DefaultConfig()
			cfg.SetRoot(viper.GetString(cli.HomeFlag))

			var countBlocks int64
			var err error
			if len(args) == 1 {
				countBlocks, err = strconv.ParseInt(args[0], 10, 64)
				if err != nil {
					return err
				}
			}

			blockStoreDB, err := node.DefaultDBProvider(&node.DBContext{ID: "blockstore", Config: cfg})
			if err != nil {
				return err
			}

			stateDB, err := node.DefaultDBProvider(&node.DBContext{ID: "state", Config: cfg})
			if err != nil {
				return err
			}

			newst := state.NewStore(stateDB)
			rootState, err := newst.Load()

			if err != nil {
				return err
			}
			blockState := store.LoadBlockStoreState(stateDB)
			height := blockState.GetHeight()

			if countBlocks > (blockState.GetHeight() - rootState.LastBlockHeight/100*100) {
				countBlocks = rootState.LastBlockHeight - rootState.LastBlockHeight/100*100
			}

			blockStore := store.NewBlockStore(blockStoreDB)
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

			block := blockStore.LoadBlock(height - 1)

			rootState.LastBlockHeight = height - 1
			rootState.LastBlockID = block.LastBlockID
			rootState.AppHash = block.AppHash
			rootState.LastResultsHash = block.LastResultsHash
			rootState.LastBlockTime = time.Unix(0, block.Time.UnixNano()-time.Second.Nanoseconds()*5)
			rootState.LastHeightValidatorsChanged = rootState.LastBlockHeight - 3

			newst.Save(rootState)

			return nil
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")

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

	for i := 0; i < int(meta.BlockID.PartSetHeader.Total); i++ {
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


	store.SaveBlockStoreState(&protostore.BlockStoreState{Height: block.Height - 1}, db)
	//store.BlockStoreStateJSON{Height: block.Height - 1}.Save(db)

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
