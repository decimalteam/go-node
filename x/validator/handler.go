package validator

import (
	"fmt"

	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler creates an sdk.Handler for all the validator type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgDeclareCandidate:
			return handleMsgDeclareCandidate(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgDeclareCandidate(ctx sdk.Context, k Keeper, msg types.MsgDeclareCandidate) sdk.Result {
	// check to see if the pubkey or sender has been registered before
	if _, err := k.GetValidator(ctx, sdk.ValAddress(msg.ValidatorAddr)); err != nil {
		return types.ErrValidatorOwnerExists(k.Codespace()).Result()
	}

	if _, err := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(msg.PubKey)); err != nil {
		return types.ErrValidatorPubKeyExists(k.Codespace()).Result()
	}

	val := types.NewValidator(sdk.ValAddress(msg.ValidatorAddr), msg.PubKey, msg.Stake, msg.Commission)
	err := k.SetValidator(ctx, val)
	if err != nil {
		return types.ErrInvalidStruct(k.Codespace()).Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeDeclareCandidate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Stake.Amount.String()),
		),
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
