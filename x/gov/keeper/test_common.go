//nolint

package keeper
/*
//DONTCOVER

import (
	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/coin"
	types2 "bitbucket.org/decimalteam/go-node/x/gov/types"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/nft"
	"bytes"
	"encoding/hex"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"


	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"



	"bitbucket.org/decimalteam/go-node/x/validator"
)
// dummy addresses used for testing
var (
	delPk1   = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB51")
	delPk2   = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB50")
	delPk3   = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB52")
	delAddr1 = sdk.AccAddress(delPk1.Address())
	delAddr2 = sdk.AccAddress(delPk2.Address())
	delAddr3 = sdk.AccAddress(delPk3.Address())

	valOpPk1    = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB53")
	valOpPk2    = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB54")
	valOpPk3    = newPubKey("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AFB55")
	valOpAddr1  = sdk.ValAddress(valOpPk1.Address())
	valOpAddr2  = sdk.ValAddress(valOpPk2.Address())
	valOpAddr3  = sdk.ValAddress(valOpPk3.Address())
	valAccAddr1 = sdk.AccAddress(valOpPk1.Address())
	valAccAddr2 = sdk.AccAddress(valOpPk2.Address())
	valAccAddr3 = sdk.AccAddress(valOpPk3.Address())

	TestAddrs = []sdk.AccAddress{
		delAddr1, delAddr2, delAddr3,
		valAccAddr1, valAccAddr2, valAccAddr3,
	}

	emptyDelAddr sdk.AccAddress
	emptyValAddr sdk.ValAddress
	emptyPubkey  crypto.PubKey

	TestProposal = types2.NewProposal(types2.Content{Title: "Title", Description: "Description"}, 1, 1, 10)
)

// TODO move to common testing framework
func newPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd
}

func makeTestCodec() *codec.LegacyAmino {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	types2.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	validator.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	coin.RegisterCodec(cdc)

	return cdc
}

func createTestInput(t *testing.T, isCheckTx bool, initPower int64) (sdk.Context, auth.AccountKeeper, Keeper, validator.Keeper, supply.Keeper, coin.Keeper) {

	initTokens := validator.TokensFromConsensusPower(initPower)

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyGov := sdk.NewKVStoreKey(types2.StoreKey)
	keyStaking := sdk.NewKVStoreKey(validator.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(validator.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyCoin := sdk.NewKVStoreKey(coin.StoreKey)
	keyMultisig := sdk.NewKVStoreKey(multisig.StoreKey)
	keyNFT := sdk.NewKVStoreKey(nft.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyGov, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyCoin, sdk.StoreTypeIAVL, db)
	require.Nil(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, tmproto.Header{ChainID: "gov-chain"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&tmproto.ConsensusParams{
			Validator: &tmproto.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := makeTestCodec()

	maccPerms := map[string][]string{
		auth.FeeCollectorName:       {supply.Burner},
		types2.ModuleName:           nil,
		validator.NotBondedPoolName: {supply.Burner, supply.Staking},
		validator.BondedPoolName:    {supply.Burner, supply.Staking},
	}

	 create module accounts
	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName, supply.Burner, supply.Minter)
	govAcc := supply.NewEmptyModuleAccount(types2.ModuleName, supply.Burner)
	notBondedPool := supply.NewEmptyModuleAccount(validator.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(validator.BondedPoolName, supply.Burner, supply.Staking)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.GetAddress().String()] = true
	blacklistedAddrs[govAcc.GetAddress().String()] = true
	blacklistedAddrs[notBondedPool.GetAddress().String()] = true
	blacklistedAddrs[bondPool.GetAddress().String()] = true

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), blacklistedAddrs)
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	coinKeeper := coin.NewKeeper(cdc, keyCoin, pk.Subspace(coin.DefaultParamspace), accountKeeper, bankKeeper, config.GetDefaultConfig(config.ChainID))

	coinConfig := config.GetDefaultConfig(config.ChainID)
	coinKeeper.SetCoin(ctx, coin.Coin{
		Title:  coinConfig.TitleBaseCoin,
		Symbol: coinConfig.SymbolBaseCoin,
		Volume: coinConfig.InitialVolumeBaseCoin,
	})

	multisigKeeper := multisig.NewKeeper(cdc, keyMultisig, pk.Subspace(multisig.DefaultParamspace), accountKeeper, coinKeeper, bankKeeper)

	nftKeeper := nft.NewKeeper(cdc, keyNFT, supplyKeeper, validator.DefaultBondDenom)

	sk := validator.NewKeeper(cdc, keyStaking, pk.Subspace(validator.DefaultParamSpace), coinKeeper, accountKeeper, supplyKeeper, multisigKeeper, nftKeeper, auth.FeeCollectorName)
	sk.SetParams(ctx, validator.DefaultParams())

	rtr := types2.NewRouter()

	keeper := NewKeeper(
		cdc, keyGov, pk.Subspace(types2.DefaultParamspace).WithKeyTable(types2.ParamKeyTable()), supplyKeeper, sk, rtr,
	)

	keeper.SetProposalID(ctx, types2.DefaultStartingProposalID)
	keeper.SetTallyParams(ctx, types2.DefaultTallyParams())

	initCoins := sdk.NewCoins(sdk.NewCoin(validator.DefaultBondDenom, initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(validator.DefaultBondDenom, initTokens.MulRaw(int64(len(TestAddrs)))))
	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	for _, addr := range TestAddrs {
		_, err := bankKeeper.AddCoins(ctx, addr, initCoins)
		require.Nil(t, err)
	}

	keeper.accountKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	keeper.accountKeeper.SetModuleAccount(ctx, govAcc)
	keeper.accountKeeper.SetModuleAccount(ctx, bondPool)
	keeper.accountKeeper.SetModuleAccount(ctx, notBondedPool)

	return ctx, accountKeeper, keeper, sk, supplyKeeper, coinKeeper
}

 //ProposalEqual checks if two proposals are equal (note: slow, for tests only)
func ProposalEqual(proposalA types2.Proposal, proposalB types2.Proposal) bool {
	return bytes.Equal(types2.ModuleCdc.MustMarshal(proposalA),
		types2.ModuleCdc.MustMarshal(proposalB))
}

func createValidators(ctx sdk.Context, vk validator.Keeper, coinKeeper coin.Keeper, supplyKeeper supply.Keeper, powers []int64) {
	val1 := validator.NewValidator(valOpAddr1, valOpPk1, sdk.ZeroDec(), sdk.AccAddress(valOpAddr1), validator.Description{})
	val2 := validator.NewValidator(valOpAddr2, valOpPk2, sdk.ZeroDec(), sdk.AccAddress(valOpAddr2), validator.Description{})
	val3 := validator.NewValidator(valOpAddr3, valOpPk3, sdk.ZeroDec(), sdk.AccAddress(valOpAddr3), validator.Description{})

	handler := validator.NewHandler(vk)

	valTokens := validator.TokensFromConsensusPower(powers[0])
	valCreateMsg := validator.NewMsgDeclareCandidate(
		valOpAddr1, valOpPk1, sdk.ZeroDec(), sdk.NewCoin(vk.BondDenom(ctx), valTokens),
		validator.Description{}, sdk.AccAddress(valOpAddr1),
	)

	_, err := handler(ctx, valCreateMsg)
	if err != nil {
		panic(err)
	}

	valTokens = validator.TokensFromConsensusPower(powers[1])
		valCreateMsg = validator.NewMsgDeclareCandidate(
		valOpAddr2, valOpPk2, sdk.ZeroDec(), sdk.NewCoin(vk.BondDenom(ctx), valTokens),
		validator.Description{}, sdk.AccAddress(valOpAddr2),
	)

	_, err = handler(ctx, valCreateMsg)
	if err != nil {
		panic(err)
	}

	valTokens = validator.TokensFromConsensusPower(powers[2])
	valCreateMsg = validator.NewMsgDeclareCandidate(
		valOpAddr3, valOpPk3, sdk.ZeroDec(), sdk.NewCoin(vk.BondDenom(ctx), valTokens),
		validator.Description{}, sdk.AccAddress(valOpAddr3),
	)

	_, err = handler(ctx, valCreateMsg)
	if err != nil {
		panic(err)
	}

	_ = validator.EndBlocker(ctx, vk, coinKeeper, supplyKeeper, false)
}
*/