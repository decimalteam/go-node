package keeper

import (
	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	nftTypes "bitbucket.org/decimalteam/go-node/x/nft/internal/types"
	"bitbucket.org/decimalteam/go-node/x/validator"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	types2 "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	types3 "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"testing"
)

// Hogpodge of all sorts of input required for testing.
// `initPower` is converted to an amount of tokens.
// If `initPower` is 0, no addrs get created.
func CreateTestInput(t *testing.T, isCheckTx bool, initPower int64) (sdk.Context, Keeper) {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyCoin := sdk.NewKVStoreKey(coin.StoreKey)
	keyMultisig := sdk.NewKVStoreKey(multisig.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyCoin, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMultisig, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	_config := sdk.GetConfig()
	_config.SetCoinType(60)
	_config.SetFullFundraiserPath("44'/60'/0'/0/0")
	_config.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	_config.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	_config.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)

	ctx := sdk.NewContext(ms, types2.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&types2.ConsensusParams{
			Validator: &types2.ValidatorParams{
				PubKeyTypes: []string{types3.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := keeper.MakeTestCodec()

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName, supply.Burner, supply.Minter)
	notBondedPool := supply.NewEmptyModuleAccount(validator.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(validator.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.String()] = true
	blacklistedAddrs[notBondedPool.String()] = true
	blacklistedAddrs[bondPool.String()] = true

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

	maccPerms := map[string][]string{
		auth.FeeCollectorName:       nil,
		validator.NotBondedPoolName: {supply.Burner, supply.Staking},
		validator.BondedPoolName:    {supply.Burner, supply.Staking},
		nftTypes.ReservedPool:       {supply.Burner},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bk, maccPerms)

	initTokens := validator.TokensFromConsensusPower(initPower)
	initCoins := sdk.NewCoins(sdk.NewCoin(validator.DefaultBondDenom, initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(validator.DefaultBondDenom, initTokens.MulRaw(int64(len(keeper.Addrs)))))

	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	coinKeeper := coin.NewKeeper(cdc, keyCoin, pk.Subspace(coin.DefaultParamspace), accountKeeper, bk, config.GetDefaultConfig(config.ChainID))

	coinConfig := config.GetDefaultConfig(config.ChainID)
	coinKeeper.SetCoin(ctx, coin.Coin{
		Title:  coinConfig.TitleBaseCoin,
		Symbol: coinConfig.SymbolBaseCoin,
		Volume: coinConfig.InitialVolumeBaseCoin,
	})
	nftkeeper := NewKeeper(cdc, keyCoin, supplyKeeper, validator.DefaultBondDenom)

	// set module accounts
	err = notBondedPool.SetCoins(totalSupply)
	require.NoError(t, err)

	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, bondPool)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range keeper.Addrs {
		_, err := bk.AddCoins(ctx, addr, initCoins)
		if err != nil {
			panic(err)
		}
	}

	return ctx, nftkeeper
}
