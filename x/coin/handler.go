package coin

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

var e18 = big.NewFloat(1000000000000000000)

func floatFromInt(amount sdk.Int) float64 {
	bigFloat := big.NewFloat(0)
	bigFloat.SetInt(amount.BigInt())
	bigFloat = bigFloat.Quo(bigFloat, e18)
	float, _ := bigFloat.Float64()
	return float
}

// NewHandler creates an sdk.Handler for all the coin type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		//ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateCoin:
			return handleMsgCreateCoin(ctx, k, msg)
		case types.MsgSendCoin:
			return handleMsgSendCoin(ctx, k, msg)
		case types.MsgMultiSendCoin:
			return handleMsgMultiSendCoin(ctx, k, msg)
		case types.MsgBuyCoin:
			return handleMsgBuyCoin(ctx, k, msg)
		case types.MsgSellCoin:
			return handleMsgSellCoin(ctx, k, msg, false)
		case types.MsgSellAllCoin:
			msgSell := MsgSellCoin{
				Seller:       msg.Seller,
				CoinToBuy:    msg.CoinToBuy,
				CoinToSell:   msg.CoinToSell,
				AmountToSell: sdk.ZeroInt(),
				AmountToBuy:  msg.AmountToBuy,
			}
			return handleMsgSellCoin(ctx, k, msgSell, true)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

////////////////////////////////////////////////////////////////
// Creating coins handlers
////////////////////////////////////////////////////////////////

func handleMsgCreateCoin(ctx sdk.Context, k Keeper, msg types.MsgCreateCoin) (*sdk.Result, error) {
	var coin = types.Coin{
		Title:       msg.Title,
		CRR:         msg.ConstantReserveRatio,
		Symbol:      msg.Symbol,
		Reserve:     msg.InitialReserve,
		LimitVolume: msg.LimitVolume,
		Volume:      msg.InitialVolume,
	}
	log.Println("Create coin gas: ", ctx.GasMeter().GasConsumed())
	// TODO: take reserve from creator and give it initial volume
	acc := k.AccountKeeper.GetAccount(ctx, msg.Creator)
	balance := acc.GetCoins()
	if balance.AmountOf(strings.ToLower(cliUtils.GetBaseCoin())).LT(msg.InitialReserve) {
		return nil, sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, "")
	}

	err := k.UpdateBalance(ctx, strings.ToLower(cliUtils.GetBaseCoin()), msg.InitialReserve.Neg(), msg.Creator)
	if err != nil {
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, "")
	}

	k.SetCoin(ctx, coin)
	err = k.UpdateBalance(ctx, strings.ToLower(coin.Symbol), msg.InitialVolume, msg.Creator)
	if err != nil {
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, "")
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeCreateCoin),
			sdk.NewAttribute(types.AttributeSymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeTitle, msg.Title),
			sdk.NewAttribute(types.AttributeInitVolume, msg.InitialVolume.String()),
			sdk.NewAttribute(types.AttributeInitReserve, msg.InitialReserve.String()),
			sdk.NewAttribute(types.AttributeCRR, strconv.FormatUint(uint64(msg.ConstantReserveRatio), 10)),
			sdk.NewAttribute(types.AttributeLimitVolume, msg.LimitVolume.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
	})

	log.Println("Create coin gas: ", ctx.GasMeter().GasConsumed())

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

////////////////////////////////////////////////////////////////
// Transfer coins handlers
////////////////////////////////////////////////////////////////

func handleMsgSendCoin(ctx sdk.Context, k Keeper, msg types.MsgSendCoin) (*sdk.Result, error) {
	// TODO: commission
	log.Println("Send coin gas: ", ctx.GasMeter().GasConsumed())
	err := k.BankKeeper.SendCoins(ctx, msg.Sender, msg.Receiver, sdk.Coins{sdk.NewCoin(strings.ToLower(msg.Coin), msg.Amount)})
	if err != nil {
		return nil, sdkerrors.New(types.DefaultCodespace, 6, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeSendCoin),
			sdk.NewAttribute(types.AttributeCoin, msg.Coin),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeReceiver, msg.Receiver.String()),
		),
	})

	log.Println("Send coin gas: ", ctx.GasMeter().GasConsumed())

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgMultiSendCoin(ctx sdk.Context, k Keeper, msg types.MsgMultiSendCoin) (*sdk.Result, error) {
	for i := range msg.Coins {
		// TODO: Commission
		_ = k.BankKeeper.SendCoins(ctx, msg.Sender, msg.Coins[i].Receiver, sdk.Coins{sdk.NewCoin(strings.ToLower(msg.Coins[i].Coin), msg.Coins[i].Amount)})
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeMultiSendCoin),
				sdk.NewAttribute(types.AttributeCoin, msg.Coins[i].Coin),
				sdk.NewAttribute(types.AttributeAmount, msg.Coins[i].Amount.String()),
				sdk.NewAttribute(types.AttributeReceiver, msg.Coins[i].Receiver.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			),
		})
	}
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

////////////////////////////////////////////////////////////////
// Trading coins handlers
////////////////////////////////////////////////////////////////

func handleMsgBuyCoin(ctx sdk.Context, k Keeper, msg types.MsgBuyCoin) (*sdk.Result, error) {
	log.Println("Buy coin gas: ", ctx.GasMeter().GasConsumed())
	// Retrieve buyer account and it's balance of selling coins
	account := k.AccountKeeper.GetAccount(ctx, msg.Buyer)
	balance := account.GetCoins().AmountOf(strings.ToLower(msg.CoinToSell))

	// Retrieve the coin requested to buy
	coinToBuy, err := k.GetCoin(ctx, msg.CoinToBuy)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to buy: %v", msg.CoinToBuy, err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}
	if coinToBuy.Symbol != msg.CoinToBuy {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to buy: retrieved coin %s instead", msg.CoinToBuy, coinToBuy.Symbol)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}

	// Retrieve the coin requested to sell
	coinToSell, err := k.GetCoin(ctx, msg.CoinToSell)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to sell: %v", msg.CoinToSell, err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}
	if coinToSell.Symbol != msg.CoinToSell {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to sell: retrieved coin %s instead", msg.CoinToSell, coinToSell.Symbol)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToSellNotExists, errMsg)
	}

	if !coinToBuy.IsBase() {
		fmt.Printf("####### Coin to buy: %s\n", coinToBuy.Symbol)
		fmt.Printf("####### - reserve: %f (%s)\n", floatFromInt(coinToBuy.Reserve), coinToBuy.Reserve)
		fmt.Printf("####### - supply: %f (%s)\n", floatFromInt(coinToBuy.Volume), coinToBuy.Volume)
		fmt.Printf("####### - limit: %f (%s)\n", floatFromInt(coinToBuy.LimitVolume), coinToBuy.LimitVolume)
	}

	if !coinToSell.IsBase() {
		fmt.Printf("####### Coin to sell: %s\n", coinToSell.Symbol)
		fmt.Printf("####### - reserve: %f (%s)\n", floatFromInt(coinToSell.Reserve), coinToSell.Reserve)
		fmt.Printf("####### - supply: %f (%s)\n", floatFromInt(coinToSell.Volume), coinToSell.Volume)
		fmt.Printf("####### - limit: %f (%s)\n", floatFromInt(coinToSell.LimitVolume), coinToSell.LimitVolume)
	}

	// Ensure supply limit of the coin to buy does not overflow
	if !coinToBuy.IsBase() {
		if coinToBuy.Volume.Add(msg.AmountToBuy).GT(coinToBuy.LimitVolume) {
			errMsg := fmt.Sprintf(
				"Wanted to buy %f %s, but this operation will overflow coin supply limit",
				floatFromInt(msg.AmountToBuy), coinToBuy.Symbol,
			)
			return nil, sdkerrors.New(types.DefaultCodespace, types.TxBreaksVolumeLimit, errMsg)
		}
	}

	// Calculate amount of sell coins which buyer will receive
	amountToBuy, amountToSell, amountInBaseCoin := msg.AmountToBuy, sdk.ZeroInt(), sdk.ZeroInt()
	switch {
	case coinToSell.IsBase():
		// Buyer buys custom coin for base coin
		amountToSell = formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountToBuy)
		amountInBaseCoin = amountToSell
	case coinToBuy.IsBase():
		// Buyer buys base coin for custom coin
		amountToSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountToBuy)
		amountInBaseCoin = amountToBuy
	default:
		// Buyer buys custom coin for custom coin
		amountInBaseCoin = formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountToBuy)
		amountToSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountInBaseCoin)
	}
	fmt.Printf("####### Buy %f (%s) %s for %s\n", floatFromInt(amountToBuy), amountToBuy, coinToBuy.Symbol, coinToSell.Symbol)
	fmt.Printf("####### - calculated amount to buy: %f (%s)\n", floatFromInt(amountToBuy), amountToBuy)
	fmt.Printf("####### - calculated amount to sell: %f (%s)\n", floatFromInt(amountToSell), amountToSell)
	fmt.Printf("####### - calculated amount in base coin: %f (%s)\n", floatFromInt(amountInBaseCoin), amountInBaseCoin)

	// Ensure maximum amount of coins to sell (price guard)
	if amountToSell.GT(msg.AmountToSell) {
		errMsg := fmt.Sprintf(
			"Wanted to sell maximum %f %s, but required to spend %f %s at the moment",
			floatFromInt(msg.AmountToSell), coinToSell.Symbol, floatFromInt(amountToSell), coinToSell.Symbol,
		)
		return nil, sdkerrors.New(types.DefaultCodespace, types.MaximumValueToSellReached, errMsg)
	}

	// Ensure reserve of the coin to sell does not underflow
	if !coinToSell.IsBase() {
		// TODO: Compare with some minimal reserve value?
		if coinToSell.Reserve.Sub(amountInBaseCoin).Sign() < 0 {
			errMsg := fmt.Sprintf(
				"Wanted to sell %f %s, but this operation will underflow coin reserve",
				floatFromInt(amountInBaseCoin), coinToSell.Symbol,
			)
			return nil, sdkerrors.New(types.DefaultCodespace, types.TxBreaksMinReserveLimit, errMsg)
		}
	}

	// Ensure that buyer account holds enough coins to sell
	if balance.LT(amountToSell) {
		errMsg := fmt.Sprintf(
			"Wanted to sell %f %s, but available only %f %s at the moment",
			floatFromInt(balance), coinToSell.Symbol, floatFromInt(amountToSell), coinToSell.Symbol,
		)
		return nil, sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, errMsg)
	}

	// Update buyer account balances
	err = k.UpdateBalance(ctx, msg.CoinToSell, amountToSell.Neg(), msg.Buyer)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to update balance of account %s: %v", account.GetAddress(), err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}
	err = k.UpdateBalance(ctx, msg.CoinToBuy, amountToBuy, msg.Buyer)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to update balance of account %s: %v", account.GetAddress(), err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}

	// Update coins
	if !coinToSell.IsBase() {
		k.UpdateCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountInBaseCoin), coinToSell.Volume.Sub(amountToSell))
	}
	if !coinToBuy.IsBase() {
		k.UpdateCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountInBaseCoin), coinToBuy.Volume.Add(amountToBuy))
	}

	// Emit transaction events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeBuyCoin),
			sdk.NewAttribute(types.AttributeCoinToBuy, msg.CoinToBuy),
			sdk.NewAttribute(types.AttributeCoinToSell, msg.CoinToSell),
			sdk.NewAttribute(types.AttributeAmountToBuy, amountToBuy.String()),
			sdk.NewAttribute(types.AttributeAmountToSell, amountToSell.String()),
			sdk.NewAttribute(types.AttributeAmountInBaseCoin, amountInBaseCoin.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer.String()),
		),
	})

	fmt.Printf("####### Buy transaction successed!\n")

	log.Println("Buy coin gas: ", ctx.GasMeter().GasConsumed())

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSellCoin(ctx sdk.Context, k Keeper, msg types.MsgSellCoin, sellAll bool) (*sdk.Result, error) {
	log.Println("Sell coin gas: ", ctx.GasMeter().GasConsumed())
	// Retrieve seller account and it's balance of selling coins
	account := k.AccountKeeper.GetAccount(ctx, msg.Seller)
	balance := account.GetCoins().AmountOf(strings.ToLower(msg.CoinToSell))

	// Fill amount to sell in case of MsgSellAll
	if sellAll {
		msg.AmountToSell = balance
	}

	// Retrieve the coin requested to sell
	coinToSell, err := k.GetCoin(ctx, msg.CoinToSell)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to sell: %v", msg.CoinToSell, err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToSellNotExists, errMsg)
	}
	if coinToSell.Symbol != msg.CoinToSell {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to sell: retrieved coin %s instead", msg.CoinToSell, coinToSell.Symbol)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToSellNotExists, errMsg)
	}

	// Retrieve the coin requested to buy
	coinToBuy, err := k.GetCoin(ctx, msg.CoinToBuy)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to buy: %v", msg.CoinToBuy, err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}
	if coinToBuy.Symbol != msg.CoinToBuy {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to buy: retrieved coin %s instead", msg.CoinToBuy, coinToBuy.Symbol)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}

	if !coinToBuy.IsBase() {
		fmt.Printf("####### Coin to buy: %s\n", coinToBuy.Symbol)
		fmt.Printf("####### - reserve: %f (%s)\n", floatFromInt(coinToBuy.Reserve), coinToBuy.Reserve)
		fmt.Printf("####### - supply: %f (%s)\n", floatFromInt(coinToBuy.Volume), coinToBuy.Volume)
		fmt.Printf("####### - limit: %f (%s)\n", floatFromInt(coinToBuy.LimitVolume), coinToBuy.LimitVolume)
	}

	if !coinToSell.IsBase() {
		fmt.Printf("####### Coin to sell: %s\n", coinToSell.Symbol)
		fmt.Printf("####### - reserve: %f (%s)\n", floatFromInt(coinToSell.Reserve), coinToSell.Reserve)
		fmt.Printf("####### - supply: %f (%s)\n", floatFromInt(coinToSell.Volume), coinToSell.Volume)
		fmt.Printf("####### - limit: %f (%s)\n", floatFromInt(coinToSell.LimitVolume), coinToSell.LimitVolume)
	}

	// Ensure that seller account holds enough coins to sell
	if balance.LT(msg.AmountToSell) {
		errMsg := fmt.Sprintf(
			"Wanted to sell %f %s, but available only %f %s at the moment",
			floatFromInt(balance), coinToSell.Symbol, floatFromInt(msg.AmountToSell), coinToSell.Symbol,
		)
		return nil, sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, errMsg)
	}

	// Calculate amount of buy coins which seller will receive
	amountToSell, amountToBuy, amountInBaseCoin := msg.AmountToSell, sdk.ZeroInt(), sdk.ZeroInt()
	switch {
	case coinToBuy.IsBase():
		// Seller sells custom coin for base coin
		amountToBuy = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountToSell)
		amountInBaseCoin = amountToBuy
	case coinToSell.IsBase():
		// Seller sells base coin for custom coin
		amountToBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountToSell)
		amountInBaseCoin = amountToSell
	default:
		// Seller sells custom coin for custom coin
		amountInBaseCoin = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountToSell)
		amountToBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountInBaseCoin)
	}
	fmt.Printf("####### Sell %f (%s) %s for %s\n", floatFromInt(amountToSell), amountToSell, coinToSell.Symbol, coinToBuy.Symbol)
	fmt.Printf("####### - calculated amount to sell: %f (%s)\n", floatFromInt(amountToSell), amountToSell)
	fmt.Printf("####### - calculated amount to buy: %f (%s)\n", floatFromInt(amountToBuy), amountToBuy)
	fmt.Printf("####### - calculated amount in base coin: %f (%s)\n", floatFromInt(amountInBaseCoin), amountInBaseCoin)

	// Ensure minimum amount of coins to buy (price guard)
	if amountToBuy.LT(msg.AmountToBuy) {
		errMsg := fmt.Sprintf(
			"Wanted to buy minimum %f %s, but expected to receive %f %s at the moment",
			floatFromInt(msg.AmountToBuy), coinToBuy.Symbol, floatFromInt(amountToBuy), coinToBuy.Symbol,
		)
		return nil, sdkerrors.New(types.DefaultCodespace, types.MinimumValueToBuyReached, errMsg)
	}

	// Ensure reserve of the coin to sell does not underflow
	if !coinToSell.IsBase() {
		// TODO: Compare with some minimal reserve value?
		if coinToSell.Reserve.Sub(amountInBaseCoin).Sign() < 0 {
			errMsg := fmt.Sprintf(
				"Wanted to sell %f %s, but this operation will underflow coin reserve",
				floatFromInt(amountInBaseCoin), coinToSell.Symbol,
			)
			return nil, sdkerrors.New(types.DefaultCodespace, types.TxBreaksMinReserveLimit, errMsg)
		}
	}

	// Ensure supply limit of the coin to buy does not overflow
	if !coinToBuy.IsBase() {
		if coinToBuy.Volume.Add(amountToBuy).GT(coinToBuy.LimitVolume) {
			errMsg := fmt.Sprintf(
				"Wanted to buy %f %s, but this operation will overflow coin supply limit",
				floatFromInt(amountToBuy), coinToBuy.Symbol,
			)
			return nil, sdkerrors.New(types.DefaultCodespace, types.TxBreaksVolumeLimit, errMsg)
		}
	}

	// Update seller account balances
	err = k.UpdateBalance(ctx, msg.CoinToSell, amountToSell.Neg(), msg.Seller)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to update balance of account %s: %v", account.GetAddress(), err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}
	err = k.UpdateBalance(ctx, msg.CoinToBuy, amountToBuy, msg.Seller)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to update balance of account %s: %v", account.GetAddress(), err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}

	// Update coins
	if !coinToSell.IsBase() {
		k.UpdateCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountInBaseCoin), coinToSell.Volume.Sub(amountToSell))
	}
	if !coinToBuy.IsBase() {
		k.UpdateCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountInBaseCoin), coinToBuy.Volume.Add(amountToBuy))
	}

	// Emit transaction events
	eventType := types.EventTypeSellCoin
	if sellAll {
		eventType = types.EventTypeSellAllCoin
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, eventType),
			sdk.NewAttribute(types.AttributeCoinToSell, msg.CoinToSell),
			sdk.NewAttribute(types.AttributeCoinToBuy, msg.CoinToBuy),
			sdk.NewAttribute(types.AttributeAmountToSell, amountToSell.String()),
			sdk.NewAttribute(types.AttributeAmountToBuy, amountToBuy.String()),
			sdk.NewAttribute(types.AttributeAmountInBaseCoin, amountInBaseCoin.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Seller.String()),
		),
	})

	fmt.Printf("####### Sell transaction successed!\n")
	log.Println("Sell coin gas: ", ctx.GasMeter().GasConsumed())

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
