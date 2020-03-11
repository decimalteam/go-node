package validator

import (
	"fmt"

	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates an sdk.Handler for all the validator type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateValidator:
			return handleMsgCreateValidator(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateValidator(ctx sdk.Context, msg types.MsgCreateValidator, k Keeper) sdk.Result {

	err := k.CreateValidator(ctx, msg.ValidatorAddr)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
