package keeper

import (
	"bitbucket.org/decimalteam/go-node/config"
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"bytes"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codec2 "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"strconv"
	"testing"
)

var (
	Addrs = createTestAddrs(500)
)

// create a codec used only for testing
func MakeTestCodec() *codec.LegacyAmino {
	var cdc = codec.NewLegacyAmino()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(types2.MsgCreateCoin{}, "test/coin/create_coin", nil)
	cdc.RegisterConcrete(types2.MsgBuyCoin{}, "test/coin/buy_coin", nil)
	cdc.RegisterConcrete(types2.MsgSellCoin{}, "test/coin/sell_coin", nil)
	cdc.RegisterConcrete(types2.MsgSendCoin{}, "test/coin/send_coin", nil)
	cdc.RegisterConcrete(types2.MsgSellAllCoin{}, "test/coin/sell_all_coin", nil)
	cdc.RegisterConcrete(types2.MsgMultiSendCoin{}, "test/coin/multi_send_coin", nil)

	// Register AppAccount
	cdc.RegisterInterface((*client.Account)(nil), nil)
	cdc.RegisterConcrete(&authtypes.BaseAccount{}, "test/coin/base_account", nil)
	codec2.RegisterCrypto(cdc)

	return cdc
}

// Hogpodge of all sorts of input required for testing.
// `initPower` is converted to an amount of tokens.
// If `initPower` is 0, no addrs get created.
func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, Keeper, authkeeper.AccountKeeper) {
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(authtypes.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(.TStoreKey)
	keyCoin := sdk.NewKVStoreKey(types2.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyCoin, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
	)
	cdc := MakeTestCodec()

	blacklistedAddrs := make(map[string]bool)

	pk := paramskeeper.NewKeeper(codec.BinaryCodec, cdc, keyParams, tKeyParams)

	accountKeeper := authkeeper.NewAccountKeeper(
		cdc,    // amino codec
		keyAcc, // target store
		pk.Subspace(authtypes.DefaultParams()),
		authtypes.ProtoBaseAccount, // prototype
	)

	bk := bankkeeper.NewBaseKeeper(
		accountKeeper,
		pk.Subspace(bank.DefaultParamspace),
		blacklistedAddrs,
	)

	coinKeeper := NewKeeper(cdc, keyCoin, pk.Subspace(types2.DefaultParamspace), accountKeeper, bk, config.GetDefaultConfig(config.ChainID))

	coinConfig := config.GetDefaultConfig(config.ChainID)
	coinKeeper.SetCoin(ctx, types2.Coin{
		Title:  coinConfig.TitleBaseCoin,
		Symbol: coinConfig.SymbolBaseCoin,
		Volume: coinConfig.InitialVolumeBaseCoin,
	})

	return ctx, coinKeeper, accountKeeper
}

// nolint: unparam
func createTestAddrs(numAddrs int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (numAddrs + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addresses = append(addresses, TestAddr(buffer.String(), bech))
		buffer.Reset()
	}
	return addresses
}

// for incode address generation
func TestAddr(addr string, bech string) sdk.AccAddress {

	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		panic(err)
	}
	bechexpected := res.String()
	if bech != bechexpected {
		panic("Bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(bechres, res) {
		panic("Bech decode and hex decode don't match")
	}

	return res
}
