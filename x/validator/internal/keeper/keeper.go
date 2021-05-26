package keeper

import (
	appTypes "bitbucket.org/decimalteam/go-node/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"

	"bitbucket.org/decimalteam/go-node/x/coin"
	"bitbucket.org/decimalteam/go-node/x/multisig"
	"bitbucket.org/decimalteam/go-node/x/nft"
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

const aminoCacheSize = 500

// Keeper of the validator store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	paramSpace       types.ParamSubspace
	CoinKeeper       coin.Keeper
	AccountKeeper    auth.AccountKeeper
	supplyKeeper     supply.Keeper
	multisigKeeper   multisig.Keeper
	nftKeeper        nft.Keeper
	hooks            types.ValidatorHooks
	FeeCollectorName string
}

// NewKeeper creates a validator keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace types.ParamSubspace, coinKeeper coin.Keeper, accountKeeper auth.AccountKeeper, supplyKeeper supply.Keeper, multisigKeeper multisig.Keeper, nftKeeper nft.Keeper, feeCollectorName string) Keeper {

	// ensure bonded and not bonded module accounts are set
	if addr := supplyKeeper.GetModuleAddress(appTypes.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", appTypes.BondedPoolName))
	}

	if addr := supplyKeeper.GetModuleAddress(appTypes.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", appTypes.NotBondedPoolName))
	}

	keeper := Keeper{
		storeKey:         key,
		cdc:              cdc,
		paramSpace:       paramSpace.WithKeyTable(ParamKeyTable()),
		CoinKeeper:       coinKeeper,
		AccountKeeper:    accountKeeper,
		supplyKeeper:     supplyKeeper,
		multisigKeeper:   multisigKeeper,
		nftKeeper:        nftKeeper,
		FeeCollectorName: feeCollectorName,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Codespace() string {
	return types.DefaultCodespace
}

// Get returns the pubkey from the adddress-pubkey relation
func (k Keeper) Get(ctx sdk.Context, key []byte, value interface{}) error {
	store := ctx.KVStore(k.storeKey)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(key), &value)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) set(ctx sdk.Context, key []byte, value interface{}) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(value)
	if err != nil {
		return err
	}
	store.Set(key, bz)
	return nil
}

func (k Keeper) delete(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(key)
}

// Load the last total validator power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) sdk.Int {
	power := sdk.Int{}
	err := k.Get(ctx, []byte{types.LastTotalPowerKey}, &power)
	if err != nil {
		return sdk.ZeroInt()
	}
	return power
}

// Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) error {
	return k.set(ctx, []byte{types.LastTotalPowerKey}, power)
}

func (k Keeper) GetCoin(ctx sdk.Context, symbol string) (coin.Coin, error) {
	if symbol == "tdel" {
		symbol = "tDEL"
	} else {
		symbol = strings.ToUpper(symbol)
	}
	return k.CoinKeeper.GetCoin(ctx, symbol)
}
