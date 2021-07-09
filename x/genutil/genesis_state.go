package genutil

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/pkg/errors"

	abci "github.com/tendermint/tendermint/abci/types"
	tos "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/types"

	"bitbucket.org/decimalteam/go-node/x/genutil/types"
	"bitbucket.org/decimalteam/go-node/x/validator"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ModuleName = "genutil"

// GenesisState defines the raw genesis transaction in JSON
//type GenesisState struct {
//	GenTxs []json.RawMessage `json:"gentxs" yaml:"gentxs"`
//}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(genTxs []json.RawMessage) *types.GenesisState {
	if len(genTxs) == 0 {
		genTxs = make([]json.RawMessage, 0)
	}

	return &types.GenesisState{
		GenTxs: genTxs,
	}
}

// NewGenesisStateFromStdTx creates a new GenesisState object
// from auth transactions
func NewGenesisStateFromStdTx(codec *codec.LegacyAmino, genTxs []sdk.Tx) *types.GenesisState {
	genTxsBz := make([]json.RawMessage, len(genTxs))
	for i, genTx := range genTxs {
		genTxsBz[i] = codec.MustMarshalJSON(genTx)
	}
	return NewGenesisState(genTxsBz)
}

// GetGenesisStateFromAppState gets the genutil genesis state from the expected app state
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) *types.GenesisState {
	var genesisState types.GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return &genesisState
}

// SetGenesisStateInAppState sets the genutil genesis state within the expected app state
func SetGenesisStateInAppState(cdc codec.JSONCodec,
	appState map[string]json.RawMessage, genesisState *types.GenesisState) map[string]json.RawMessage {

	genesisStateBz := cdc.MustMarshalJSON(genesisState)
	appState[ModuleName] = genesisStateBz
	return appState
}

// GenesisStateFromGenDoc creates the core parameters for genesis initialization
// for the application.
//
// NOTE: The pubkey input is this machines pubkey.
func GenesisStateFromGenDoc(cdc codec.JSONCodec, genDoc tmtypes.GenesisDoc,
) (genesisState map[string]json.RawMessage, err error) {

	if err = json.Unmarshal(genDoc.AppState, &genesisState); err != nil {
		return genesisState, err
	}
	return genesisState, nil
}

// GenesisStateFromGenFile creates the core parameters for genesis initialization
// for the application.
//
// NOTE: The pubkey input is this machines pubkey.
func GenesisStateFromGenFile(cdc codec.JSONCodec, genFile string,
) (genesisState map[string]json.RawMessage, genDoc *tmtypes.GenesisDoc, err error) {

	if !tos.FileExists(genFile) {
		return genesisState, genDoc,
			fmt.Errorf("%s does not exist, run `init` first", genFile)
	}
	genDoc, err = tmtypes.GenesisDocFromFile(genFile)
	if err != nil {
		return genesisState, genDoc, err
	}

	genesisState, err = GenesisStateFromGenDoc(cdc, *genDoc)
	return genesisState, genDoc, err
}

// ValidateGenesis validates GenTx transactions
func ValidateGenesis(genesisState *types.GenesisState, txJSONDecoder sdk.TxDecoder) error {
	for i, genTx := range genesisState.GenTxs {
		tx, err := txJSONDecoder(genTx)
		if err != nil {
			return err
		}

		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			return errors.New(
				"must provide genesis StdTx with exactly 1 CreateValidator message")
		}

		if _, ok := msgs[0].(*validator.MsgDeclareCandidate); !ok {
			return fmt.Errorf(
				"genesis transaction %v does not contain a MsgCreateValidator", i)
		}
	}
	return nil
}

// InitGenesis - initialize accounts and deliver genesis transactions
func InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, validatorKeeper validator.Keeper,
	deliverTx deliverTxfn, genesisState types.GenesisState, txConfig client.TxConfig) []abci.ValidatorUpdate {

	var validators []abci.ValidatorUpdate
	var err error
	if len(genesisState.GenTxs) > 0 {
		validators, err = DeliverGenTxs(ctx, cdc, genesisState.GenTxs, validatorKeeper, deliverTx, txConfig)
		if err != nil {
			panic(err)
		}
	}
	return validators
}
