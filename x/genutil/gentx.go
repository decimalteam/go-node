package genutil

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
	"github.com/cosmos/cosmos-sdk/x/bank/exported"
	"path/filepath"

	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	tos "github.com/tendermint/tendermint/libs/os"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"

	"bitbucket.org/decimalteam/go-node/x/validator"
)

// SetGenTxsInAppGenesisState - sets the genesis transactions in the app genesis state
func SetGenTxsInAppGenesisState(cdc codec.JSONCodec, txJSONEncoder sdk.TxEncoder, appGenesisState map[string]json.RawMessage,
	genTxs []sdk.Tx) (map[string]json.RawMessage, error) {

	genesisState := GetGenesisStateFromAppState(cdc, appGenesisState)
	// convert all the GenTxs to JSON
	var genTxsBz []json.RawMessage
	for _, genTx := range genTxs {
		txBz, err := txJSONEncoder(genTx)
		if err != nil {
			return appGenesisState, err
		}
		genTxsBz = append(genTxsBz, txBz)
	}

	genesisState.GenTxs = genTxsBz
	return SetGenesisStateInAppState(cdc, appGenesisState, genesisState), nil
}

// ValidateAccountInGenesis checks that the provided key has sufficient
// coins in the genesis accounts
func ValidateAccountInGenesis(appGenesisState map[string]json.RawMessage,
	genBalIterator types.GenesisBalancesIterator,
	key sdk.AccAddress, coins sdk.Coins, cdc codec.JSONCodec) error {

	accountIsInGenesis := false

	genUtilDataBz := appGenesisState[ModuleName]
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(genUtilDataBz, &genesisState)

	var err error
	genBalIterator.IterateGenesisBalances(cdc, appGenesisState,
		func(bal exported.GenesisBalance) (stop bool) {
			accAddress := bal.GetAddress()
			accCoins := bal.GetCoins()

			// ensure that account is in genesis
			if accAddress.Equals(key) {
				// ensure account contains enough funds of default bond denom
				if coins.AmountOf(validator.DefaultBondDenom).GT(accCoins.AmountOf(validator.DefaultBondDenom)) {
					err = fmt.Errorf(
						"account %s has a balance in genesis, but it only has %v%s available to stake, not %v%s",
						key, accCoins.AmountOf(validator.DefaultBondDenom), validator.DefaultBondDenom, coins.AmountOf(validator.DefaultBondDenom), validator.DefaultBondDenom,
					)

					return true
				}

				accountIsInGenesis = true
				return true
			}

			return false
		},
	)
	if err != nil {
		return err
	}

	if !accountIsInGenesis {
		return fmt.Errorf("account %s in not in the app_state.accounts array of genesis.json", key)
	}

	return nil
}

type deliverTxfn func(abci.RequestDeliverTx) abci.ResponseDeliverTx

// DeliverGenTxs - deliver a genesis transaction
func DeliverGenTxs(ctx sdk.Context, cdc codec.JSONCodec, genTxs []json.RawMessage,
	validatorKeeper validator.Keeper, deliverTx deliverTxfn, txEncodingConfig client.TxEncodingConfig) ([]abci.ValidatorUpdate, error) {

	for _, genTx := range genTxs {
		//var tx legacytx.StdTx
		//cdc.MustUnmarshalJSON(genTx, &tx)
		//bz := cdc.MustMarshalJSON(tx)
		tx, err := txEncodingConfig.TxJSONDecoder()(genTx)
		if err != nil {
			panic(err)
		}
		bz, err := txEncodingConfig.TxEncoder()(tx)
		if err != nil {
			panic(err)
		}
		//cdc.MustUnmarshalJSON(genTx, &tx)
		//bz := cdc.MustMarshalJSON(tx)
		res := deliverTx(abci.RequestDeliverTx{Tx: bz})
		if !res.IsOK() {
			panic(res.Log)
		}
	}
	return validatorKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
}

// InitializeNodeValidatorFiles creates private validator and p2p configuration files.
func InitializeNodeValidatorFiles(config *cfg.Config,
) (nodeID string, valPubKey crypto.PubKey, err error) {

	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return nodeID, valPubKey, err
	}

	nodeID = string(nodeKey.ID())
	server.UpgradeOldPrivValFile(config)

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := tos.EnsureDir(filepath.Dir(pvKeyFile), 0777); err != nil {
		return nodeID, valPubKey, nil
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := tos.EnsureDir(filepath.Dir(pvStateFile), 0777); err != nil {
		return nodeID, valPubKey, nil
	}

	valPubKey, err = privval.LoadOrGenFilePV(pvKeyFile, pvStateFile).GetPubKey()

	if err != nil {
		return nodeID, valPubKey, err
	}

	return nodeID, valPubKey, nil
}
