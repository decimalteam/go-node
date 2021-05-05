package keeper // noalias

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"strconv"
	"testing"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/nft"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// dummy addresses used for testing
// nolint: unused deadcode
var (
	Addrs = createTestAddrs(500)
	PKs   = createTestPubKeys(500)

	addrDels = []sdk.AccAddress{
		Addrs[0],
		Addrs[1],
	}
	addrVals = []sdk.ValAddress{
		sdk.ValAddress(Addrs[2]),
		sdk.ValAddress(Addrs[3]),
		sdk.ValAddress(Addrs[4]),
		sdk.ValAddress(Addrs[5]),
		sdk.ValAddress(Addrs[6]),
	}
)

//_______________________________________________________________________________________

// intended to be used with require/assert:  require.True(ValEq(...))
func ValEq(t *testing.T, exp, got types.Validator) (*testing.T, bool, string, types.Validator, types.Validator) {
	return t, exp.TestEquivalent(got), "expected:\t%v\ngot:\t\t%v", exp, got
}

//_______________________________________________________________________________________

// create a codec used only for testing
func MakeTestCodec() *codec.Codec {
	var cdc = codec.New()

	// Register Msgs
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(bank.MsgSend{}, "test/validator/send", nil)
	cdc.RegisterConcrete(types.MsgDeclareCandidate{}, "test/validator/declare_candidate", nil)
	cdc.RegisterConcrete(types.MsgEditCandidate{}, "test/validator/edit_candidate", nil)
	cdc.RegisterConcrete(types.MsgSetOnline{}, "test/validator/set_online", nil)
	cdc.RegisterConcrete(types.MsgSetOffline{}, "test/validator/set_offline", nil)
	cdc.RegisterConcrete(types.MsgUnbond{}, "test/validator/unbond", nil)
	cdc.RegisterConcrete(types.MsgDelegate{}, "test/validator/delegate", nil)

	// Register AppAccount
	cdc.RegisterInterface((*authexported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/validator/base_account", nil)
	supply.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc
}

// Hogpodge of all sorts of input required for testing.
// `initPower` is converted to an amount of tokens.
// If `initPower` is 0, no addrs get created.
func CreateTestInput(t *testing.T, isCheckTx bool, initPower int64) (sdk.Context, auth.AccountKeeper, Keeper, supply.Keeper, coin.Keeper) {
	keyStaking := sdk.NewKVStoreKey(types.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(types.TStoreKey)
	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyCoin := sdk.NewKVStoreKey(coin.StoreKey)
	keyMultisig := sdk.NewKVStoreKey(multisig.StoreKey)
	keyNFT := sdk.NewKVStoreKey(nft.StoreKey)

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
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

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "foochainid"}, isCheckTx, log.NewNopLogger())
	ctx = ctx.WithConsensusParams(
		&abci.ConsensusParams{
			Validator: &abci.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
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
	}
	supplyKeeper := supply.NewKeeper(cdc, keySupply, accountKeeper, bk, maccPerms)

	initTokens := types.TokensFromConsensusPower(initPower)
	initCoins := sdk.NewCoins(sdk.NewCoin(types.DefaultBondDenom, initTokens))
	totalSupply := sdk.NewCoins(sdk.NewCoin(types.DefaultBondDenom, initTokens.MulRaw(int64(len(Addrs)))))

	supplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	coinKeeper := coin.NewKeeper(cdc, keyCoin, pk.Subspace(coin.DefaultParamspace), accountKeeper, bk, config.GetDefaultConfig(config.ChainID))

	coinConfig := config.GetDefaultConfig(config.ChainID)
	coinKeeper.SetCoin(ctx, coin.Coin{
		Title:  coinConfig.TitleBaseCoin,
		Symbol: coinConfig.SymbolBaseCoin,
		Volume: coinConfig.InitialVolumeBaseCoin,
	})

	multisigKeeper := multisig.NewKeeper(cdc, keyMultisig, pk.Subspace(multisig.DefaultParamspace), accountKeeper, coinKeeper, bk)

	nftKeeper := nft.NewKeeper(cdc, keyNFT, supplyKeeper, types.DefaultBondDenom)

	keeper := NewKeeper(cdc, keyStaking, pk.Subspace(DefaultParamspace), coinKeeper, accountKeeper, supplyKeeper, multisigKeeper, nftKeeper, auth.FeeCollectorName)
	keeper.SetParams(ctx, types.DefaultParams())

	// set module accounts
	err = notBondedPool.SetCoins(totalSupply)
	require.NoError(t, err)

	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, bondPool)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range Addrs {
		_, err := bk.AddCoins(ctx, addr, initCoins)
		if err != nil {
			panic(err)
		}
	}

	return ctx, accountKeeper, keeper, supplyKeeper, coinKeeper
}

func NewPubKey(pk string) (res crypto.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	//res, err = crypto.PubKeyFromBytes(pkBytes)
	var pkEd ed25519.PubKeyEd25519
	copy(pkEd[:], pkBytes[:])
	return pkEd
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

// nolint: unparam
func createTestPubKeys(numPubKeys int) []crypto.PubKey {
	var publicKeys []crypto.PubKey
	var buffer bytes.Buffer

	//start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF") //base pubkey string
		buffer.WriteString(numString)                                                       //adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, NewPubKey(buffer.String()))
		buffer.Reset()
	}
	return publicKeys
}

//_____________________________________________________________________________________

// does a certain by-power index record exist
func ValidatorByPowerIndexExists(ctx sdk.Context, keeper Keeper, power []byte) bool {
	store := ctx.KVStore(keeper.storeKey)
	return store.Has(power)
}

// update validator for testing
func TestingUpdateValidator(keeper Keeper, ctx sdk.Context, validator types.Validator, apply bool) types.Validator {
	keeper.SetValidator(ctx, validator)

	// Remove any existing power key for validator.
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{types.ValidatorsByPowerIndexKey})
	defer iterator.Close()
	deleted := false
	for ; iterator.Valid(); iterator.Next() {
		valAddr := types.ParseValidatorPowerRankKey(iterator.Key())
		if bytes.Equal(valAddr, validator.ValAddress) {
			if deleted {
				panic("found duplicate power index key")
			} else {
				deleted = true
			}
			store.Delete(iterator.Key())
		}
	}

	keeper.SetValidatorByPowerIndex(ctx, validator)
	if apply {
		_, err := keeper.ApplyAndReturnValidatorSetUpdates(ctx)
		if err != nil {
			panic(err)
		}
		validator, err := keeper.GetValidator(ctx, validator.ValAddress)
		if err != nil {
			panic("validator expected but not found")
		}
		return validator
	}
	cachectx, _ := ctx.CacheContext()
	_, err := keeper.ApplyAndReturnValidatorSetUpdates(cachectx)
	if err != nil {
		panic(err)
	}
	validator, err = keeper.GetValidator(cachectx, validator.ValAddress)
	if err != nil {
		panic("validator expected but not found")
	}
	return validator
}

// nolint: deadcode unused
func validatorByPowerIndexExists(k Keeper, ctx sdk.Context, power []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(power)
}

// RandomValidator returns a random validator given access to the keeper and ctx
func RandomValidator(r *rand.Rand, keeper Keeper, ctx sdk.Context) types.Validator {
	vals := keeper.GetAllValidators(ctx)
	i := r.Intn(len(vals))
	return vals[i]
}

// RandomBondedValidator returns a random bonded validator given access to the keeper and ctx
func RandomBondedValidator(r *rand.Rand, keeper Keeper, ctx sdk.Context) types.Validator {
	vals := keeper.GetBondedValidatorsByPower(ctx)
	i := r.Intn(len(vals))
	return vals[i]
}
