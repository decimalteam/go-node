package keeper

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/multisig/internal/types"
)

var (
	Addrs = createTestAddrs(500)
)

// create a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(types.MsgCreateWallet{}, "test/coin/create_wallet", nil)
	cdc.RegisterConcrete(types.MsgCreateTransaction{}, "test/coin/create_transaction", nil)
	cdc.RegisterConcrete(types.MsgSignTransaction{}, "test/coin/sign_transaction", nil)

	// Register AppAccount
	cdc.RegisterInterface((*authexported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/coin/base_account", nil)
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

// Hogpodge of all sorts of input required for testing.
// `initPower` is converted to an amount of tokens.
// If `initPower` is 0, no addrs get created.
func CreateTestInput(t *testing.T, isCheckTx bool) (sdk.Context, Keeper, coin.Keeper, auth.AccountKeeper, bank.Keeper) {
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	tKeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyCoin := sdk.NewKVStoreKey(coin.StoreKey)
	keyMultisig := sdk.NewKVStoreKey(types.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyCoin, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMultisig, sdk.StoreTypeIAVL, db)
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

	pk := params.NewKeeper(cdc, keyParams, tKeyParams)

	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), blacklistedAddrs)
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bankKeeper, make(map[string][]string))
	coinKeeper := coin.NewKeeper(cdc, keyCoin, pk.Subspace(coin.DefaultParamspace), accountKeeper, bankKeeper, supplyKeeper, config.GetDefaultConfig(config.ChainID))

	coinConfig := config.GetDefaultConfig(config.ChainID)
	coinKeeper.SetCoin(ctx, coin.Coin{
		Title:  coinConfig.TitleBaseCoin,
		Symbol: coinConfig.SymbolBaseCoin,
		Volume: coinConfig.InitialVolumeBaseCoin,
	})

	multisigKeeper := NewKeeper(cdc, keyMultisig, pk.Subspace(types.DefaultParamspace), accountKeeper, coinKeeper, bankKeeper)

	return ctx, multisigKeeper, coinKeeper, accountKeeper, bankKeeper
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
