// nolint
// DONTCOVER
package gov

import (
	"bytes"
	"sort"
	"testing"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/nft"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

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
	valTokens  = validator.TokensFromConsensusPower(42)
	initTokens = validator.TokensFromConsensusPower(100000)
	valCoins   = sdk.NewCoins(sdk.NewCoin(validator.DefaultBondDenom, valTokens))
	initCoins  = sdk.NewCoins(sdk.NewCoin(validator.DefaultBondDenom, initTokens))
)

type testInput struct {
	ctx      sdk.Context
	keeper   keep.Keeper
	router   types.Router
	vk       validator.Keeper
	ck       coin.Keeper
	sk       supply.Keeper
	addrs    []sdk.AccAddress
	pubKeys  []crypto.PubKey
	privKeys []crypto.PrivKey
}

func getTestInput(t *testing.T, numGenAccs int, genState types.GenesisState, genAccs []authexported.Account) testInput {

	_config := sdk.GetConfig()
	_config.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	_config.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	_config.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)

	keyValidator := sdk.NewKVStoreKey(validator.StoreKey)
	tkeyValidator := sdk.NewTransientStoreKey(validator.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyCoin := sdk.NewKVStoreKey(coin.StoreKey)
	keyMultisig := sdk.NewKVStoreKey(multisig.StoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyNFT := sdk.NewKVStoreKey(nft.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(tkeyValidator, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyValidator, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyCoin, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMultisig, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, false, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)

	cdc := codec.New()
	validator.RegisterCodec(cdc)
	types.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	coin.RegisterCodec(cdc)
	cdc.RegisterInterface((*authexported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/gov/base_account", nil)
	codec.RegisterCrypto(cdc)

	govAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Burner)
	notBondedPool := supply.NewEmptyModuleAccount(validator.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(validator.BondedPoolName, supply.Burner, supply.Staking)
	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName, supply.Burner, supply.Minter)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[govAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true
	blacklistedAddrs[feeCollectorAcc.String()] = true

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)

	accountKeeper := auth.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		pk.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount, // prototype
	)

	bk := bank.NewBaseKeeper(
		accountKeeper,
		pk.Subspace(bank.DefaultParamspace),
		blacklistedAddrs,
	)

	rtr := types.NewRouter()

	maccPerms := map[string][]string{
		auth.FeeCollectorName:       {supply.Burner, supply.Minter},
		types.ModuleName:            {supply.Burner},
		validator.NotBondedPoolName: {supply.Burner, supply.Staking},
		validator.BondedPoolName:    {supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bk, maccPerms)

	coinKeeper := coin.NewKeeper(cdc, keyCoin, pk.Subspace(coin.DefaultParamspace), accountKeeper, bk, supplyKeeper, config.GetDefaultConfig(config.ChainID))
	coinConfig := config.GetDefaultConfig(config.ChainID)
	coinKeeper.SetCoin(ctx, coin.Coin{
		Title:  coinConfig.TitleBaseCoin,
		Symbol: coinConfig.SymbolBaseCoin,
		Volume: coinConfig.InitialVolumeBaseCoin,
	})

	nftKeeper := nft.NewKeeper(cdc, keyNFT, supplyKeeper, validator.DefaultBondDenom)

	multisigKeeper := multisig.NewKeeper(cdc, keyMultisig, pk.Subspace(multisig.DefaultParamspace), accountKeeper, coinKeeper, bk)
	sk := validator.NewKeeper(cdc, keyValidator, pk.Subspace(validator.DefaultParamSpace), coinKeeper, accountKeeper, supplyKeeper, multisigKeeper, nftKeeper, auth.FeeCollectorName)
	sk.SetParams(ctx, validator.DefaultParams())

	keeper := keep.NewKeeper(
		cdc, keyCoin, pk.Subspace(DefaultParamspace).WithKeyTable(ParamKeyTable()), supplyKeeper, sk, rtr,
	)
	keeper.SetTallyParams(ctx, types.DefaultParams().TallyParams)
	keeper.SetProposalID(ctx, 1)

	var (
		addrs    []sdk.AccAddress
		pubKeys  []crypto.PubKey
		privKeys []crypto.PrivKey
	)

	if genAccs == nil || len(genAccs) == 0 {
		genAccs, addrs, pubKeys, privKeys = mock.CreateGenAccounts(numGenAccs, valCoins)
	}

	initCoins := sdk.NewCoins(sdk.NewCoin(validator.DefaultBondDenom, initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(validator.DefaultBondDenom, initTokens.MulRaw(int64(len(addrs)))))

	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))
	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)

	for _, addr := range addrs {
		_, err := bk.AddCoins(ctx, addr, initCoins)
		require.NoError(t, err)
	}

	return testInput{ctx, keeper, rtr, sk, coinKeeper, supplyKeeper, addrs, pubKeys, privKeys}
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

		totalSupply := sdk.NewCoins(sdk.NewCoin(validator.DefaultBondDenom, initTokens.MulRaw(int64(len(mapp.GenesisAccounts)))))
		supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

		// set module accounts
		for _, macc := range blacklistedAddrs {
			supplyKeeper.SetModuleAccount(ctx, macc)
		}

		validators := validator.InitGenesis(ctx, stakingKeeper, supplyKeeper, stakingGenesis)
		if genState.IsEmpty() {
			InitGenesis(ctx, keeper, types.DefaultGenesisState())
		} else {
			InitGenesis(ctx, keeper, genState)
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
		panic("not fail-able with `bytes.Comparable` bounded [-1, 1].")
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

		valTokens := validator.TokensFromConsensusPower(powerAmt[i])
		valCreateMsg := validator.NewMsgDeclareCandidate(
			addrs[i], pubkeys[i], sdk.ZeroDec(), sdk.NewCoin(validator.DefaultBondDenom, valTokens),
			validator.Description{}, sdk.AccAddress(addrs[i]),
		)

		res, err := stakingHandler(ctx, valCreateMsg)
		require.NoError(t, err)
		require.NotNil(t, res)
	}
}
