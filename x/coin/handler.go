package coin

import (
	"fmt"

	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates an sdk.Handler for all the coin type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateCoin:
			return handleMsgCreateCoin(ctx, k, msg)
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
		InitialAmount:        msg.InitialAmount,
		InitialReserve:       msg.InitialReserve,
		LimitAmount:          msg.LimitAmount,
	}
	existCoin, _ := k.GetCoin(ctx, coin.Symbol)
	if existCoin.Symbol != "" {
		return sdk.NewError(DefaultCodespace, types.CoinAlreadyExists, fmt.Sprintf("Coin with symbol %s already exists", coin.Symbol)).Result()
	}
	k.SetCoin(ctx, coin)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeCreateCoin),
			sdk.NewAttribute(types.AttributeSymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeTitle, msg.Title),
			sdk.NewAttribute(types.AttributeInitAmount, msg.InitialAmount.String()),
			sdk.NewAttribute(types.AttributeInitReserve, msg.InitialReserve.String()),
			sdk.NewAttribute(types.AttributeCRR, string(msg.ConstantReserveRatio)),
			sdk.NewAttribute(types.AttributeLimitAmount, msg.LimitAmount.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
