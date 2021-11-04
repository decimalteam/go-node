package coin

/*
import (
	"bitbucket.org/decimalteam/go-node/utils/helpers"
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	keep "bitbucket.org/decimalteam/go-node/x/coin/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func createCoin(ctx sdk.Context, keeper Keeper) Coin {
	volume := helpers.BipToPip(sdk.NewInt(100000))
	reserve := helpers.BipToPip(sdk.NewInt(100000))

	coin := Coin{
		Title:       "TEST COIN",
		CRR:         10,
		Symbol:      "test",
		Reserve:     reserve,
		LimitVolume: volume.Mul(sdk.NewInt(10)),
		Volume:      volume,
	}

	keeper.SetCoin(ctx, coin)

	return coin
}

func TestBuyCoinTxBaseToCustom(t *testing.T) {
	ctx, keeper, accountKeeper := keep.CreateTestInput(t, false)

	coin := createCoin(ctx, keeper)

	initBalance := helpers.BipToPip(sdk.NewInt(1000000))

	account := accountKeeper.NewAccountWithAddress(ctx, keep.Addrs[0])
	err := account.SetCoins(sdk.NewCoins(sdk.NewCoin(cliUtils.GetBaseCoin(), initBalance)))
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, account)

	toBuy := helpers.BipToPip(sdk.NewInt(10))
	maxValToSell, ok := sdk.NewIntFromString("159374246010000000000")
	require.True(t, ok)

	buyCoinMsg := NewMsgBuyCoin(keep.Addrs[0], sdk.NewCoin(coin.Symbol, toBuy), sdk.NewCoin(cliUtils.GetBaseCoin(), maxValToSell))
	_, err = handleMsgBuyCoin(ctx, keeper, buyCoinMsg)
	require.NoError(t, err)

	targetBalance, ok := sdk.NewIntFromString("999899954987997899747979")
	require.True(t, ok)

	account = accountKeeper.GetAccount(ctx, keep.Addrs[0])
	require.NotNil(t, account)

	balance := account.GetCoins().AmountOf(cliUtils.GetBaseCoin())
	if !balance.Equal(targetBalance) {
		t.Fatalf("Target %s initBalance is not correct. Expected %s, got %s", cliUtils.GetBaseCoin(), targetBalance, balance)
	}

	testBalance := account.GetCoins().AmountOf(coin.Symbol)
	if !testBalance.Equal(toBuy) {
		t.Fatalf("Target %s balance is not correct. Expected %s, got %s", coin.Symbol, toBuy, testBalance)
	}
}

func TestBuyCoinTxInsufficientFunds(t *testing.T) {
	ctx, keeper, accountKeeper := keep.CreateTestInput(t, false)

	coin := createCoin(ctx, keeper)

	initBalance := helpers.BipToPip(sdk.NewInt(1))

	account := accountKeeper.NewAccountWithAddress(ctx, keep.Addrs[0])
	err := account.SetCoins(sdk.NewCoins(sdk.NewCoin(cliUtils.GetBaseCoin(), initBalance)))
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, account)

	toBuy := helpers.BipToPip(sdk.NewInt(10))
	maxValToSell, ok := sdk.NewIntFromString("159374246010000000000")
	require.True(t, ok)

	buyCoinMsg := NewMsgBuyCoin(keep.Addrs[0], sdk.NewCoin(coin.Symbol, toBuy), sdk.NewCoin(cliUtils.GetBaseCoin(), maxValToSell))
	_, err = handleMsgBuyCoin(ctx, keeper, buyCoinMsg)
	require.EqualError(t, err, types.ErrInsufficientFunds("100045012002100252021", "1000000000000000000").Error())
}

func TestBuyCoinTxEqualCoins(t *testing.T) {
	ctx, keeper, accountKeeper := keep.CreateTestInput(t, false)

	coin := createCoin(ctx, keeper)

	initBalance := helpers.BipToPip(sdk.NewInt(100000))

	account := accountKeeper.NewAccountWithAddress(ctx, keep.Addrs[0])
	err := account.SetCoins(sdk.NewCoins(sdk.NewCoin(cliUtils.GetBaseCoin(), initBalance)))
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, account)

	toBuy := helpers.BipToPip(sdk.NewInt(10))
	maxValToSell, ok := sdk.NewIntFromString("159374246010000000000")
	require.True(t, ok)

	buyCoinMsg := NewMsgBuyCoin(keep.Addrs[0], sdk.NewCoin(coin.Symbol, toBuy), sdk.NewCoin(coin.Symbol, maxValToSell))
	err = buyCoinMsg.ValidateBasic()
	require.EqualError(t, err, types.ErrSameCoin().Error())
}

func TestBuyCoinTxNotExistsBuyCoin(t *testing.T) {
	ctx, keeper, accountKeeper := keep.CreateTestInput(t, false)

	coin := createCoin(ctx, keeper)

	initBalance := helpers.BipToPip(sdk.NewInt(1))

	account := accountKeeper.NewAccountWithAddress(ctx, keep.Addrs[0])
	err := account.SetCoins(sdk.NewCoins(sdk.NewCoin(cliUtils.GetBaseCoin(), initBalance)))
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, account)

	toBuy := helpers.BipToPip(sdk.NewInt(10))
	maxValToSell, ok := sdk.NewIntFromString("159374246010000000000")
	require.True(t, ok)

	buyCoinMsg := NewMsgBuyCoin(keep.Addrs[0], sdk.NewCoin("invalid", toBuy), sdk.NewCoin(coin.Symbol, maxValToSell))
	_, err = handleMsgBuyCoin(ctx, keeper, buyCoinMsg)
	require.EqualError(t, err, types.ErrCoinDoesNotExist("invalid").Error())
}

func TestBuyCoinTxNotExistsSellCoin(t *testing.T) {
	ctx, keeper, accountKeeper := keep.CreateTestInput(t, false)

	coin := createCoin(ctx, keeper)

	initBalance := helpers.BipToPip(sdk.NewInt(1))

	account := accountKeeper.NewAccountWithAddress(ctx, keep.Addrs[0])
	err := account.SetCoins(sdk.NewCoins(sdk.NewCoin(cliUtils.GetBaseCoin(), initBalance)))
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, account)

	toBuy := helpers.BipToPip(sdk.NewInt(10))
	maxValToSell, ok := sdk.NewIntFromString("159374246010000000000")
	require.True(t, ok)

	buyCoinMsg := NewMsgBuyCoin(keep.Addrs[0], sdk.NewCoin(coin.Symbol, toBuy), sdk.NewCoin("invalid", maxValToSell))
	_, err = handleMsgBuyCoin(ctx, keeper, buyCoinMsg)
	require.EqualError(t, err, types.ErrCoinDoesNotExist("invalid").Error())
}

func TestBuyCoinTxCustomToBase(t *testing.T) {
	ctx, keeper, accountKeeper := keep.CreateTestInput(t, false)

	coin := createCoin(ctx, keeper)

	initBalance := helpers.BipToPip(sdk.NewInt(10000000))

	account := accountKeeper.NewAccountWithAddress(ctx, keep.Addrs[0])
	err := account.SetCoins(sdk.NewCoins(sdk.NewCoin(coin.Symbol, initBalance)))
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, account)

	toBuy := helpers.BipToPip(sdk.NewInt(10))
	maxValToSell, ok := sdk.NewIntFromString("159374246010000000000")
	require.True(t, ok)

	buyCoinMsg := NewMsgBuyCoin(keep.Addrs[0], sdk.NewCoin(cliUtils.GetBaseCoin(), toBuy), sdk.NewCoin(coin.Symbol, maxValToSell))
	_, err = handleMsgBuyCoin(ctx, keeper, buyCoinMsg)
	require.NoError(t, err)

	targetBalance, ok := sdk.NewIntFromString("9999998999954997149793304")
	require.True(t, ok)

	account = accountKeeper.GetAccount(ctx, keep.Addrs[0])
	require.NotNil(t, account)

	balance := account.GetCoins().AmountOf(coin.Symbol)
	require.Equal(t, balance, targetBalance, "Target %s balance is not correct. Expected %s, got %s", coin.Symbol, targetBalance, balance)

	baseBalance := account.GetCoins().AmountOf(cliUtils.GetBaseCoin())
	require.Equal(t, baseBalance, toBuy, "Target %s balance is not correct. Expected %s, got %s", cliUtils.GetBaseCoin(), toBuy, baseBalance)

	coin, err = keeper.GetCoin(ctx, coin.Symbol)
	require.NoError(t, err)

	targetReserve, ok := sdk.NewIntFromString("99990000000000000000000")
	require.True(t, ok)
	require.Equal(t, targetReserve, coin.Reserve, "Target %s reserve is not correct. Expected %s, got %s", coin.Symbol, targetReserve, coin.Reserve)

	targetVolume, ok := sdk.NewIntFromString("99998999954997149793304")
	require.True(t, ok)
	require.Equal(t, targetVolume, coin.Volume, "Target %s volume is not correct. Expected %s, got %s", coin.Symbol, targetVolume, coin.Volume)
}

func TestBuyCoinReserveUnderflow(t *testing.T) {
	ctx, keeper, accountKeeper := keep.CreateTestInput(t, false)

	coin := createCoin(ctx, keeper)

	initBalance := helpers.BipToPip(sdk.NewInt(10000000))

	account := accountKeeper.NewAccountWithAddress(ctx, keep.Addrs[0])
	err := account.SetCoins(sdk.NewCoins(sdk.NewCoin(cliUtils.GetBaseCoin(), initBalance)))
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, account)

	toBuy := helpers.BipToPip(sdk.NewInt(99900))
	maxValToSell, ok := sdk.NewIntFromString("49881276637272773421684")
	require.True(t, ok)

	buyCoinMsg := NewMsgBuyCoin(keep.Addrs[0], sdk.NewCoin(cliUtils.GetBaseCoin(), toBuy), sdk.NewCoin(coin.Symbol, maxValToSell))
	_, err = handleMsgBuyCoin(ctx, keeper, buyCoinMsg)
	require.EqualError(t, err, types.ErrTxBreaksMinReserveRule(MinCoinReserve(ctx).String(), toBuy.String()).Error())
}

func TestSellCoinTxBaseToCustom(t *testing.T) {
	ctx, keeper, accountKeeper := keep.CreateTestInput(t, false)

	coin := createCoin(ctx, keeper)

	initBalance := helpers.BipToPip(sdk.NewInt(1000000))

	account := accountKeeper.NewAccountWithAddress(ctx, keep.Addrs[0])
	err := account.SetCoins(sdk.NewCoins(sdk.NewCoin(cliUtils.GetBaseCoin(), initBalance)))
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, account)

	valueToSell := helpers.BipToPip(sdk.NewInt(10))
	minValToBuy, ok := sdk.NewIntFromString("957658277688702625")
	require.True(t, ok)

	sellCoinMsg := NewMsgSellCoin(keep.Addrs[0], sdk.NewCoin(cliUtils.GetBaseCoin(), valueToSell), sdk.NewCoin(coin.Symbol, minValToBuy))
	_, err = handleMsgSellCoin(ctx, keeper, sellCoinMsg, false)
	require.NoError(t, err)

	targetBalance, ok := sdk.NewIntFromString("999990000000000000000000")
	require.True(t, ok)

	account = accountKeeper.GetAccount(ctx, keep.Addrs[0])
	require.NotNil(t, account)

	balance := account.GetCoins().AmountOf(cliUtils.GetBaseCoin())
	require.Equal(t, balance, targetBalance, "Target %s balance is not correct. Expected %s, got %s", coin.Symbol, targetBalance, balance)

	targetTestBalance, ok := sdk.NewIntFromString("999955002849793446")
	require.True(t, ok)

	testBalance := account.GetCoins().AmountOf(coin.Symbol)
	require.Equal(t, testBalance, targetTestBalance, "Target %s balance is not correct. Expected %s, got %s", coin.Symbol, targetTestBalance, testBalance)
}

func TestSellAllCoinTx(t *testing.T) {
	ctx, keeper, accountKeeper := keep.CreateTestInput(t, false)

	coin := createCoin(ctx, keeper)

	initBalance := helpers.BipToPip(sdk.NewInt(1000000))

	account := accountKeeper.NewAccountWithAddress(ctx, keep.Addrs[0])
	err := account.SetCoins(sdk.NewCoins(sdk.NewCoin(cliUtils.GetBaseCoin(), initBalance)))
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, account)

	minValToBuy, ok := sdk.NewIntFromString("151191152412701306252")
	require.True(t, ok)

	sellCoinMsg := types.NewMsgSellCoin(keep.Addrs[0], sdk.NewCoin(cliUtils.GetBaseCoin(), sdk.ZeroInt()), sdk.NewCoin(coin.Symbol, minValToBuy))
	_, err = handleMsgSellCoin(ctx, keeper, sellCoinMsg, true)
	require.NoError(t, err)

	account = accountKeeper.GetAccount(ctx, keep.Addrs[0])
	require.NotNil(t, account)

	balance := account.GetCoins().AmountOf(cliUtils.GetBaseCoin())
	require.Equal(t, balance, sdk.ZeroInt(), "Target %s balance is not correct. Expected %s, got %s", coin.Symbol, sdk.ZeroInt(), balance)

	targetTestBalance, ok := sdk.NewIntFromString("27098161521014065552356")
	require.True(t, ok)

	testBalance := account.GetCoins().AmountOf(coin.Symbol)
	require.Equal(t, testBalance, targetTestBalance, "Target %s balance is not correct. Expected %s, got %s", coin.Symbol, targetTestBalance, testBalance)
}

func TestCreateCoinTx(t *testing.T) {
	ctx, keeper, accountKeeper := keep.CreateTestInput(t, false)

	initBalance := helpers.BipToPip(sdk.NewInt(1000000))

	account := accountKeeper.NewAccountWithAddress(ctx, keep.Addrs[0])
	err := account.SetCoins(sdk.NewCoins(sdk.NewCoin(cliUtils.GetBaseCoin(), initBalance)))
	require.NoError(t, err)
	accountKeeper.SetAccount(ctx, account)

	reserve := helpers.BipToPip(sdk.NewInt(10000))
	volume := helpers.BipToPip(sdk.NewInt(100))
	crr := uint(50)
	title := "My Test Coin"
	symbol := "ABCDEF"

	sellCoinMsg := types.NewMsgCreateCoin(keep.Addrs[0], title, symbol, crr, volume, reserve, volume.MulRaw(10), "")
	_, err = handleMsgCreateCoin(ctx, keeper, sellCoinMsg)
	require.NoError(t, err)

	account = accountKeeper.GetAccount(ctx, keep.Addrs[0])
	require.NotNil(t, account)

	targetBalance, ok := sdk.NewIntFromString("989000000000000000000000")
	require.True(t, ok)

	balance := account.GetCoins().AmountOf(cliUtils.GetBaseCoin())
	require.Equal(t, balance, targetBalance, "Target %s balance is not correct. Expected %s, got %s", cliUtils.GetBaseCoin(), targetBalance, balance)

	targetTestBalance := volume

	testBalance := account.GetCoins().AmountOf(strings.ToLower(symbol))
	require.Equal(t, testBalance, targetTestBalance, "Target %s balance is not correct. Expected %s, got %s", symbol, targetTestBalance, testBalance)
}
*/
