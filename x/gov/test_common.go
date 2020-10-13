// nolint
// DONTCOVER
package gov

import (
	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bytes"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"log"
	"sort"
	"testing"

	keep "bitbucket.org/decimalteam/go-node/x/gov/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	"bitbucket.org/decimalteam/go-node/x/validator"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/supply"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var (
	valTokens  = sdk.TokensFromConsensusPower(42)
	initTokens = sdk.TokensFromConsensusPower(100000)
	valCoins   = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, valTokens))
	initCoins  = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens))
)

type testInput struct {
	mApp     *mock.App
	keeper   keep.Keeper
	router   types.Router
	sk       validator.Keeper
	addrs    []sdk.AccAddress
	pubKeys  []crypto.PubKey
	privKeys []crypto.PrivKey
}

func getMockApp(t *testing.T, numGenAccs int, genState types.GenesisState, genAccs []authexported.Account) testInput {

	mApp := mock.NewApp()

	validator.RegisterCodec(mApp.Cdc)
	types.RegisterCodec(mApp.Cdc)
	supply.RegisterCodec(mApp.Cdc)
	coin.RegisterCodec(mApp.Cdc)

	keyValidator := sdk.NewKVStoreKey(validator.StoreKey)
	keyGov := sdk.NewKVStoreKey(types.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyCoin := sdk.NewKVStoreKey(coin.StoreKey)
	keyMultisig := sdk.NewKVStoreKey(multisig.StoreKey)

	govAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Burner)
	notBondedPool := supply.NewEmptyModuleAccount(validator.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(validator.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[govAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	pk := mApp.ParamsKeeper

	rtr := types.NewRouter()

	bk := bank.NewBaseKeeper(mApp.AccountKeeper, mApp.ParamsKeeper.Subspace(bank.DefaultParamspace), blacklistedAddrs)

	maccPerms := map[string][]string{
		types.ModuleName:            {supply.Burner},
		validator.NotBondedPoolName: {supply.Burner, supply.Staking},
		validator.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(mApp.Cdc, keySupply, mApp.AccountKeeper, bk, maccPerms)
	coinKeeper := coin.NewKeeper(mApp.Cdc, keyCoin, pk.Subspace(coin.DefaultParamspace), mApp.AccountKeeper, bk, config.GetDefaultConfig(config.ChainID))
	multisigKeeper := multisig.NewKeeper(mApp.Cdc, keyMultisig, pk.Subspace(multisig.DefaultParamspace), mApp.AccountKeeper, coinKeeper, bk)
	sk := validator.NewKeeper(
		mApp.Cdc, keyValidator, pk.Subspace(validator.DefaultParamSpace), coinKeeper, mApp.AccountKeeper, supplyKeeper, multisigKeeper, auth.FeeCollectorName,
	)

	keeper := keep.NewKeeper(
		mApp.Cdc, keyGov, pk.Subspace(DefaultParamspace).WithKeyTable(ParamKeyTable()), supplyKeeper, sk, rtr,
	)

	mApp.Router().AddRoute(types.RouterKey, NewHandler(keeper))
	mApp.QueryRouter().AddRoute(types.QuerierRoute, keep.NewQuerier(keeper))

	mApp.SetEndBlocker(getEndBlocker(keeper))
	mApp.SetInitChainer(getInitChainer(mApp, keeper, sk, supplyKeeper, genState, []supplyexported.ModuleAccountI{govAcc, notBondedPool, bondPool}))

	require.NoError(t, mApp.CompleteSetup(keyValidator, keyGov, keySupply))

	var (
		addrs    []sdk.AccAddress
		pubKeys  []crypto.PubKey
		privKeys []crypto.PrivKey
	)

	if genAccs == nil || len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(numGenAccs, valCoins)
	}

	mock.SetGenesis(mApp, genAccs)

	return testInput{mApp, keeper, rtr, sk, addrs, pubKeys, privKeys}
}

// gov and staking endblocker
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		EndBlocker(ctx, keeper)
		return abci.ResponseEndBlock{}
	}
}

// gov and staking initchainer
func getInitChainer(mapp *mock.App, keeper Keeper, stakingKeeper validator.Keeper, supplyKeeper supply.Keeper, genState GenesisState, blacklistedAddrs []supplyexported.ModuleAccountI) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		stakingGenesis := validator.DefaultGenesisState()

		totalSupply := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initTokens.MulRaw(int64(len(mapp.GenesisAccounts)))))
		supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		validators := validator.InitGenesis(ctx, stakingKeeper, supplyKeeper, stakingGenesis)
		if genState.IsEmpty() {
			InitGenesis(ctx, keeper, supplyKeeper, types.DefaultGenesisState())
		} else {
			InitGenesis(ctx, keeper, supplyKeeper, genState)
		}
		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}

// SortAddresses - Sorts Addresses
func SortAddresses(addrs []sdk.AccAddress) {
	var byteAddrs [][]byte
	for _, addr := range addrs {
		byteAddrs = append(byteAddrs, addr.Bytes())
	}
	SortByteArrays(byteAddrs)
	for i, byteAddr := range byteAddrs {
		addrs[i] = byteAddr
	}
}

// implement `Interface` in sort package.
type sortByteArrays [][]byte

func (b sortByteArrays) Len() int {
	return len(b)
}

func (b sortByteArrays) Less(i, j int) bool {
	// bytes package already implements Comparable for []byte.
	switch bytes.Compare(b[i], b[j]) {
	case -1:
		return true
	case 0, 1:
		return false
	default:
		log.Panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
		return false
	}
}

func (b sortByteArrays) Swap(i, j int) {
	b[j], b[i] = b[i], b[j]
}

// SortByteArrays - sorts the provided byte array
func SortByteArrays(src [][]byte) [][]byte {
	sorted := sortByteArrays(src)
	sort.Sort(sorted)
	return sorted
}

const contextKeyBadProposal = "contextKeyBadProposal"

var (
	pubkeys = []crypto.PubKey{
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
		ed25519.GenPrivKey().PubKey(),
	}
)

func createValidators(t *testing.T, stakingHandler sdk.Handler, ctx sdk.Context, addrs []sdk.ValAddress, powerAmt []int64) {
	require.True(t, len(addrs) <= len(pubkeys), "Not enough pubkeys specified at top of file.")

	for i := 0; i < len(addrs); i++ {

		valTokens := sdk.TokensFromConsensusPower(powerAmt[i])
		valCreateMsg := validator.NewMsgDeclareCandidate(
			addrs[i], pubkeys[i], sdk.ZeroDec(), sdk.NewCoin(sdk.DefaultBondDenom, valTokens),
			validator.Description{}, sdk.AccAddress(addrs[i]),
		)

		res, err := stakingHandler(ctx, valCreateMsg)
		require.NoError(t, err)
		require.NotNil(t, res)
	}
}
