package keeper

import (
	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/types"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/nft/exported"
	nftTypes "bitbucket.org/decimalteam/go-node/x/nft/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	exported2 "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abcitypes "github.com/tendermint/tendermint/abci/types"
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

	ctx := sdk.NewContext(ms, abcitypes.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abcitypes.ConsensusParams{
			Validator: &abcitypes.ValidatorParams{
				PubKeyTypes: []string{types3.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := MakeTestCodec()

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName, supply.Burner, supply.Minter)
	notBondedPool := supply.NewEmptyModuleAccount(types.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(types.BondedPoolName, supply.Burner, supply.Staking)

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
		auth.FeeCollectorName:   nil,
		types.NotBondedPoolName: {supply.Burner, supply.Staking},
		types.BondedPoolName:    {supply.Burner, supply.Staking},
		nftTypes.ReservedPool:   {supply.Burner},
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bk, maccPerms)

	initTokens := types.TokensFromConsensusPower(initPower)
	initCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultBondDenom, initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(types.DefaultBondDenom, initTokens.MulRaw(int64(len(nftTypes.Addrs)))))

	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	coinKeeper := coin.NewKeeper(cdc, keyCoin, pk.Subspace(coin.DefaultParamspace), accountKeeper, bk, config.GetDefaultConfig(config.ChainID))

	coinConfig := config.GetDefaultConfig(config.ChainID)
	coinKeeper.SetCoin(ctx, coin.Coin{
		Title:  coinConfig.TitleBaseCoin,
		Symbol: coinConfig.SymbolBaseCoin,
		Volume: coinConfig.InitialVolumeBaseCoin,
	})
	nftkeeper := NewKeeper(cdc, keyCoin, supplyKeeper, types.DefaultBondDenom)

	// set module accounts
	err = notBondedPool.SetCoins(totalSupply)
	require.NoError(t, err)

	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, bondPool)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// fill all the Addrs with some coins, set the loose pool tokens simultaneously
	for _, addr := range nftTypes.Addrs {
		_, err := bk.AddCoins(ctx, addr, initCoins)
		if err != nil {
			panic(err)
		}
	}

	return ctx, nftkeeper
}

// MakeTestCodec creates a codec for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()

	// Register Msgs
	cdc.RegisterInterface((*exported.NFT)(nil), nil)
	cdc.RegisterInterface((*exported.TokenOwners)(nil), nil)
	cdc.RegisterInterface((*exported.TokenOwner)(nil), nil)
	cdc.RegisterConcrete(&nftTypes.BaseNFT{}, "nft/BaseNFT", nil)
	cdc.RegisterConcrete(&nftTypes.IDCollection{}, "nft/IDCollection", nil)
	cdc.RegisterConcrete(&nftTypes.Collection{}, "nft/Collection", nil)
	cdc.RegisterConcrete(&nftTypes.Owner{}, "nft/Owner", nil)
	cdc.RegisterConcrete(&nftTypes.TokenOwner{}, "nft/TokenOwner", nil)
	cdc.RegisterConcrete(&nftTypes.TokenOwners{}, "nft/TokenOwners", nil)
	cdc.RegisterConcrete(nftTypes.MsgTransferNFT{}, "nft/msg_transfer", nil)
	cdc.RegisterConcrete(nftTypes.MsgEditNFTMetadata{}, "nft/msg_edit_metadata", nil)
	cdc.RegisterConcrete(nftTypes.MsgMintNFT{}, "nft/msg_mint", nil)
	cdc.RegisterConcrete(nftTypes.MsgBurnNFT{}, "nft/msg_burn", nil)

	// Register AppAccount
	cdc.RegisterInterface((*exported2.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/coin/base_account", nil)
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}
