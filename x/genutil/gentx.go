package genutil

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
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
func SetGenTxsInAppGenesisState(cdc codec.JSONCodec, appGenesisState map[string]json.RawMessage,
	genTxs []legacytx.StdTx) (map[string]json.RawMessage, error) {

	genesisState := GetGenesisStateFromAppState(cdc, appGenesisState)
	// convert all the GenTxs to JSON
	var genTxsBz []json.RawMessage
	for _, genTx := range genTxs {
		txBz, err := cdc.MarshalJSON(genTx)
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
	genAccIterator types.GenesisAccountsIterator,
	key sdk.AccAddress, coin sdk.Coin, cdc codec.JSONCodec) error {

	accountIsInGenesis := false

	genUtilDataBz := appGenesisState[ModuleName]
	var genesisState GenesisState
	cdc.MustUnmarshalJSON(genUtilDataBz, &genesisState)

	var err error
	genAccIterator.IterateGenesisAccounts(cdc, appGenesisState,
		func(acc authtypes.AccountI) (stop bool) {
			accAddress := acc.GetAddress()

			fmt.Println(coin)

			// Ensure that account is in genesis
			if accAddress.Equals(key) {

				// Ensure account contains enough funds of default bond denom
				if coin.Amount.GT(accCoins.AmountOf(coin.Denom)) {
					err = fmt.Errorf(
						"account %v is in genesis, but it only has %v%v available to stake, not %v%v",
						key.String(), accCoins.AmountOf(coin.Denom), coin.Denom, coin.Amount, coin.Denom,
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
	validatorKeeper validator.Keeper, deliverTx deliverTxfn) ([]abci.ValidatorUpdate, error) {

	for _, genTx := range genTxs {
		var tx legacytx.StdTx
		cdc.MustUnmarshalJSON(genTx, &tx)
		bz := cdc.MustMarshalJSON(tx)
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
