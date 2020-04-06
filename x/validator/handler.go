package validator

import (
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

// NewHandler creates an sdk.Handler for all the validator type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgDeclareCandidate:
			return handleMsgDeclareCandidate(ctx, keeper, msg)
		case types.MsgDelegate:
			return handleMsgDelegate(ctx, keeper, msg)
		case types.MsgUnbond:
			return handleMsgUnbond(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgDeclareCandidate(ctx sdk.Context, k Keeper, msg types.MsgDeclareCandidate) sdk.Result {
	// check to see if the pubkey or sender has been registered before
	if _, err := k.GetValidator(ctx, msg.ValidatorAddr); err == nil {
		return types.ErrValidatorOwnerExists(k.Codespace()).Result()
	}

	if _, err := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(msg.PubKey)); err == nil {
		return types.ErrValidatorPubKeyExists(k.Codespace()).Result()
	}

	val := types.NewValidator(msg.ValidatorAddr, msg.PubKey, msg.Commission)
	err := k.SetValidator(ctx, val)
	if err != nil {
		return types.ErrInvalidStruct(k.Codespace()).Result()
	}
	err = k.SetValidatorByConsAddr(ctx, val)
	if err != nil {
		return types.ErrInvalidStruct(k.Codespace()).Result()
	}
	err = k.SetNewValidatorByPowerIndex(ctx, val)
	if err != nil {
		return types.ErrInvalidStruct(k.Codespace()).Result()
	}

	k.AfterValidatorCreated(ctx, val.ValAddress)

	_, err = k.Delegate(ctx, sdk.AccAddress(msg.ValidatorAddr), msg.Stake, types.Unbonded, val, true)
	if err != nil {
		return sdk.NewError(k.Codespace(), types.CodeInvalidDelegation, err.Error()).Result()
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

func handleMsgDelegate(ctx sdk.Context, k Keeper, msg types.MsgDelegate) sdk.Result {
	val, err := k.GetValidator(ctx, msg.ValidatorAddress)
	if err != nil {
		return types.ErrValidatorOwnerExists(k.Codespace()).Result()
	}

	_, err = k.Delegate(ctx, msg.DelegatorAddress, msg.Amount, types.Unbonded, val, true)
	if err != nil {
		return sdk.NewError(k.Codespace(), 1, err.Error()).Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	})

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgUnbond(ctx sdk.Context, k Keeper, msg types.MsgUnbond) sdk.Result {
	completionTime, err := k.Undelegate(ctx, msg.DelegatorAddress, msg.ValidatorAddress, msg.Amount)
	if err != nil {
		return sdk.NewError(k.Codespace(), 1, err.Error()).Result()
	}

	completionTimeBz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(completionTime)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbond,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	})

	return sdk.Result{Data: completionTimeBz, Events: ctx.EventManager().Events()}
}
