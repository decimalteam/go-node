package coin

import (
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// NewHandler creates an sdk.Handler for all the coin type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateCoin:
			return handleMsgCreateCoin(ctx, k, msg)
		case types.MsgBuyCoin:
			return handleMsgBuyCoin(ctx, k, msg)
		case types.MsgSellCoin:
			return handleMsgSellCoin(ctx, k, msg)
		case types.MsgSendCoin:
			return handleMsgSendCoin(ctx, k, msg)
		case types.MsgMultiSendCoin:
			return handleMsgMultiSendCoin(ctx, k, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateCoin(ctx sdk.Context, k Keeper, msg types.MsgCreateCoin) sdk.Result {
	var coin = types.Coin{
		Title:       msg.Title,
		CRR:         msg.ConstantReserveRatio,
		Symbol:      msg.Symbol,
		Reserve:     msg.InitialReserve,
		LimitVolume: msg.LimitVolume,
		Volume:      msg.InitialVolume,
	}

	k.SetCoin(ctx, coin)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeCreateCoin),
			sdk.NewAttribute(types.AttributeSymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeTitle, msg.Title),
			sdk.NewAttribute(types.AttributeInitVolume, msg.InitialVolume.String()),
			sdk.NewAttribute(types.AttributeInitReserve, msg.InitialReserve.String()),
			sdk.NewAttribute(types.AttributeCRR, string(msg.ConstantReserveRatio)),
			sdk.NewAttribute(types.AttributeLimitVolume, msg.LimitVolume.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgBuyCoin(ctx sdk.Context, k Keeper, msg types.MsgBuyCoin) sdk.Result {
	// TODO: Commission

	coinToBuy, _ := k.GetCoin(ctx, msg.CoinToBuy)
	if coinToBuy.Symbol != msg.CoinToBuy {
		return sdk.Result{
			Code:      types.CoinToBuyNotExists,
			Codespace: types.DefaultCodespace,
		}
	}
	// Check if coin to sell exists
	coinToSell, _ := k.GetCoin(ctx, msg.CoinToSell)
	if coinToSell.Symbol != msg.CoinToSell {
		return sdk.Result{
			Code:      types.CoinToSellNotExists,
			Codespace: types.DefaultCodespace,
		}
	}

	amountBuy, amountSell, calcErr := cliUtils.BuyCoinCalculateAmounts(coinToBuy, coinToSell, msg.AmountToBuy, msg.AmountToSell)

	if calcErr != nil {
		return sdk.Result{Codespace: calcErr.Codespace(), Code: calcErr.Code()}
	}

	acc := k.AccountKeeper.GetAccount(ctx, msg.Buyer)
	balance := acc.GetCoins()
	if balance.AmountOf(strings.ToLower(msg.CoinToSell)).LT(amountSell) {
		return sdk.Result{
			Code:      types.InsufficientCoinToSell,
			Codespace: types.DefaultCodespace,
		}
	}

	err := k.UpdateBalance(ctx, msg.CoinToBuy, amountBuy, msg.Buyer)
	if err != nil {
		return sdk.Result{
			Code:      types.UpdateBalanceError,
			Codespace: types.DefaultCodespace,
		}
	}
	err = k.UpdateBalance(ctx, msg.CoinToSell, amountSell.Neg(), msg.Buyer)
	if err != nil {
		return sdk.Result{
			Code:      types.UpdateBalanceError,
			Codespace: types.DefaultCodespace,
		}
	}
	amountBuyInBaseCoin := formulas.CalculateSaleReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountBuy)
	amountSellInBaseCoin := formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountSell)

	k.UpdateCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountSellInBaseCoin), coinToSell.Volume.Sub(amountSell))
	k.UpdateCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountBuyInBaseCoin), coinToBuy.Volume.Add(amountBuy))
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeBuyCoin),
			sdk.NewAttribute(types.AttributeCoinToBuy, msg.CoinToBuy),
			sdk.NewAttribute(types.AttributeCoinToSell, msg.CoinToSell),
			sdk.NewAttribute(types.AttributeAmountToBuy, amountBuy.String()),
			sdk.NewAttribute(types.AttributeAmountToSell, amountSell.String()),
			sdk.NewAttribute(types.AttributeAmountToBuyInBaseCoin, amountBuyInBaseCoin.String()),
			sdk.NewAttribute(types.AttributeAmountToSellInBaseCoin, amountSellInBaseCoin.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSellCoin(ctx sdk.Context, k Keeper, msg types.MsgSellCoin) sdk.Result {
	coinToBuy, _ := k.GetCoin(ctx, msg.CoinToBuy)
	if coinToBuy.Symbol != msg.CoinToBuy {
		return sdk.Result{
			Code:      types.CoinToBuyNotExists,
			Codespace: types.DefaultCodespace,
		}
	}
	// Check if coin to sell exists
	coinToSell, _ := k.GetCoin(ctx, msg.CoinToSell)
	if coinToSell.Symbol != msg.CoinToSell {
		return sdk.Result{
			Code:      types.CoinToSellNotExists,
			Codespace: types.DefaultCodespace,
		}
	}

	amountBuy, amountSell, calcErr := cliUtils.SellCoinCalculateAmounts(coinToBuy, coinToSell, msg.AmountToBuy, msg.AmountToSell)

	if calcErr != nil {
		return sdk.Result{Codespace: calcErr.Codespace(), Code: calcErr.Code()}
	}

	acc := k.AccountKeeper.GetAccount(ctx, msg.Seller)
	balance := acc.GetCoins()
	if balance.AmountOf(strings.ToLower(msg.CoinToSell)).LT(amountSell) {
		return sdk.Result{
			Code:      types.InsufficientCoinToSell,
			Codespace: types.DefaultCodespace,
		}
	}

	err := k.UpdateBalance(ctx, msg.CoinToBuy, amountBuy, msg.Seller)
	if err != nil {
		return sdk.Result{
			Code:      types.UpdateBalanceError,
			Codespace: types.DefaultCodespace,
		}
	}
	err = k.UpdateBalance(ctx, msg.CoinToSell, amountSell.Neg(), msg.Seller)
	if err != nil {
		return sdk.Result{
			Code:      types.UpdateBalanceError,
			Codespace: types.DefaultCodespace,
		}
	}
	amountBuyInBaseCoin := formulas.CalculateSaleReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountBuy)
	amountSellInBaseCoin := formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountSell)
	k.UpdateCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountBuyInBaseCoin), coinToSell.Volume.Sub(amountSell))
	k.UpdateCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountSellInBaseCoin), coinToBuy.Volume.Add(amountBuy))

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeSellCoin),
			sdk.NewAttribute(types.AttributeCoinToBuy, msg.CoinToBuy),
			sdk.NewAttribute(types.AttributeCoinToSell, msg.CoinToSell),
			sdk.NewAttribute(types.AttributeAmountToBuy, amountBuy.String()),
			sdk.NewAttribute(types.AttributeAmountToSell, amountSell.String()),
			sdk.NewAttribute(types.AttributeAmountToBuyInBaseCoin, amountBuyInBaseCoin.String()),
			sdk.NewAttribute(types.AttributeAmountToSellInBaseCoin, amountSellInBaseCoin.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgSendCoin(ctx sdk.Context, k Keeper, msg types.MsgSendCoin) sdk.Result {
	// TODO: commission
	_ = k.BankKeeper.SendCoins(ctx, msg.Sender, msg.Receiver, sdk.Coins{sdk.NewCoin(strings.ToLower(msg.Coin), msg.Amount)})
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeSendCoin),
			sdk.NewAttribute(types.AttributeCoin, msg.Coin),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeReceiver, string(msg.Receiver)),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}

// TODO: add it (don't know how to add it to cli)
func handleMsgMultiSendCoin(ctx sdk.Context, k Keeper, msg types.MsgMultiSendCoin) sdk.Result {
	for i := range msg.Coins {
		// TODO: Commission
		_ = k.BankKeeper.SendCoins(ctx, msg.Sender, msg.Coins[i].Receiver, sdk.Coins{sdk.NewCoin(strings.ToLower(msg.Coins[i].Coin), msg.Coins[i].Amount)})
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeMultiSendCoin),
				sdk.NewAttribute(types.AttributeCoin, msg.Coins[i].Coin),
				sdk.NewAttribute(types.AttributeAmount, msg.Coins[i].Amount.String()),
				sdk.NewAttribute(types.AttributeReceiver, string(msg.Coins[i].Receiver)),
			),
		)
	}
	return sdk.Result{Events: ctx.EventManager().Events()}

}
