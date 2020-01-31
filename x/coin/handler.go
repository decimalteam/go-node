package coin

import (
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateCoin(ctx sdk.Context, k Keeper, msg types.MsgCreateCoin) sdk.Result {
	var coin = types.Coin{
		Title:                msg.Title,
		ConstantReserveRatio: msg.ConstantReserveRatio,
		Symbol:               msg.Symbol,
		Reserve:              msg.InitialReserve,
		LimitVolume:          msg.LimitVolume,
		Volume:               msg.InitialVolume,
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
	_ = k.AddCoin(ctx, msg.CoinToBuy, msg.AmountToBuy, msg.Buyer)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeBuyCoin),
			sdk.NewAttribute(types.AttributeCoinToBuy, msg.CoinToBuy),
			sdk.NewAttribute(types.AttributeCoinToSell, msg.CoinToSell),
			sdk.NewAttribute(types.AttributeAmountToBuy, msg.AmountToBuy.String()),
			sdk.NewAttribute(types.AttributeMaxAmountToSell, msg.MaximumAmountToSell.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
