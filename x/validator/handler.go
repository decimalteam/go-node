package validator

import (
	"errors"
	"fmt"
	"log"
	"time"

	tmstrings "github.com/tendermint/tendermint/libs/strings"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

// NewHandler creates an sdk.Handler for all the validator type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgDeclareCandidate:
			return handleMsgDeclareCandidate(ctx, keeper, msg)
		case types.MsgDelegate:
			return handleMsgDelegate(ctx, keeper, msg)
		case types.MsgUnbond:
			return handleMsgUnbond(ctx, keeper, msg)
		case types.MsgEditCandidate:
			return handleMsgEditCandidate(ctx, keeper, msg)
		case types.MsgSetOnline:
			return handleMsgSetOnline(ctx, keeper, msg)
		case types.MsgSetOffline:
			return handleMsgSetOffline(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgDeclareCandidate(ctx sdk.Context, k Keeper, msg types.MsgDeclareCandidate) (*sdk.Result, error) {
	// check to see if the pubkey or sender has been registered before
	if _, err := k.GetValidator(ctx, msg.ValidatorAddr); err == nil {
		return nil, types.ErrValidatorOwnerExists()
	}

	if _, err := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(msg.PubKey)); err == nil {
		return nil, types.ErrValidatorPubKeyExists()
	}

	if ctx.ConsensusParams() != nil {
		tmPubKey := tmtypes.TM2PB.PubKey(msg.PubKey)
		if !tmstrings.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
			return nil, sdkerrors.Wrapf(
				types.ErrValidatorPubKeyTypeNotSupported(),
				"got: %s, valid: %s", tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes,
			)
		}
	}

	val := types.NewValidator(msg.ValidatorAddr, msg.PubKey, msg.Commission, msg.RewardAddr, msg.Description)
	err := k.SetValidator(ctx, val)
	if err != nil {
		return nil, types.ErrInvalidStruct()
	}
	k.SetValidatorByConsAddr(ctx, val)
	k.SetNewValidatorByPowerIndex(ctx, val)

	k.AfterValidatorCreated(ctx, val.ValAddress)

	err = k.Delegate(ctx, sdk.AccAddress(msg.ValidatorAddr), msg.Stake, types.Unbonded, val, true)
	if err != nil {
		e := sdkerrors.Error{}
		if errors.As(err, &e) {
			return nil, e
		} else {
			return nil, types.ErrInternal(err.Error())
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(msg.ValidatorAddr).String()),
		sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddr.String()),
		sdk.NewAttribute(types.AttributeKeyCoin, msg.Stake.String()),
		sdk.NewAttribute(types.AttributeKeyPubKey, msg.PubKey.Address().String()),
		sdk.NewAttribute(types.AttributeKeyCommission, msg.Commission.String()),
		sdk.NewAttribute(types.AttributeKeyDescriptionMoniker, msg.Description.Moniker),
		sdk.NewAttribute(types.AttributeKeyDescriptionIdentity, msg.Description.Identity),
		sdk.NewAttribute(types.AttributeKeyDescriptionWebsite, msg.Description.Website),
		sdk.NewAttribute(types.AttributeKeyDescriptionSecurityContact, msg.Description.SecurityContact),
		sdk.NewAttribute(types.AttributeKeyDescriptionDetails, msg.Description.Details),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDelegate(ctx sdk.Context, k Keeper, msg types.MsgDelegate) (*sdk.Result, error) {
	val, err := k.GetValidator(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, types.ErrNoValidatorFound()
	}

	ok, err := k.IsDelegatorStakeSufficient(ctx, val, msg.DelegatorAddress, msg.Coin)
	if err != nil {
		return nil, types.ErrCoinDoesNotExist(msg.Coin.Denom)
	}
	if !ok {
		return nil, types.ErrDelegatorStakeIsTooLow()
	}

	t := time.Now()
	err = k.Delegate(ctx, msg.DelegatorAddress, msg.Coin, types.Unbonded, val, true)
	if err != nil {
		e := sdkerrors.Error{}
		if errors.As(err, &e) {
			return nil, e
		} else {
			return nil, types.ErrInternal(err.Error())
		}
	}
	log.Println(time.Since(t))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
		sdk.NewAttribute(types.AttributeKeyCoin, msg.Coin.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgUnbond(ctx sdk.Context, k Keeper, msg types.MsgUnbond) (*sdk.Result, error) {
	completionTime, err := k.Undelegate(ctx, msg.DelegatorAddress, msg.ValidatorAddress, msg.Coin)
	if err != nil {
		e := sdkerrors.Error{}
		if errors.As(err, &e) {
			return nil, e
		} else {
			return nil, types.ErrInternal(err.Error())
		}
	}

	completionTimeBz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(completionTime)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
		sdk.NewAttribute(types.AttributeKeyCoin, msg.Coin.String()),
		sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
	))

	return &sdk.Result{Data: completionTimeBz, Events: ctx.EventManager().Events()}, nil
}

func handleMsgEditCandidate(ctx sdk.Context, k Keeper, msg types.MsgEditCandidate) (*sdk.Result, error) {
	var validator types.Validator

	validator, err := k.GetValidator(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, types.ErrNoValidatorFound()
	}

	validator.RewardAddress = msg.RewardAddress
	validator.Description = msg.Description

	err = k.SetValidator(ctx, validator)
	if err != nil {
		return nil, types.ErrInternal(err.Error())
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(validator.ValAddress).String()),
		sdk.NewAttribute(types.AttributeKeyRewardAddress, msg.RewardAddress.String()),
		sdk.NewAttribute(types.AttributeKeyDescriptionMoniker, msg.Description.Moniker),
		sdk.NewAttribute(types.AttributeKeyDescriptionDetails, msg.Description.Details),
		sdk.NewAttribute(types.AttributeKeyDescriptionIdentity, msg.Description.Identity),
		sdk.NewAttribute(types.AttributeKeyDescriptionWebsite, msg.Description.Website),
		sdk.NewAttribute(types.AttributeKeyDescriptionSecurityContact, msg.Description.SecurityContact),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSetOnline(ctx sdk.Context, k Keeper, msg types.MsgSetOnline) (*sdk.Result, error) {
	validator, err := k.GetValidator(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, types.ErrNoValidatorFound()
	}

	if validator.Online {
		if !validator.Jailed {
			return nil, types.ErrValidatorAlreadyOnline()
		}
	}

	validator.Online = true
	validator.Jailed = false
	err = k.SetValidator(ctx, validator)
	if err != nil {
		return nil, types.ErrInternal(err.Error())
	}
	k.SetValidatorByPowerIndex(ctx, validator)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(msg.ValidatorAddress).String()),
		sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSetOffline(ctx sdk.Context, k Keeper, msg types.MsgSetOffline) (*sdk.Result, error) {
	validator, err := k.GetValidator(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, types.ErrNoValidatorFound()
	}

	if !validator.Online {
		return nil, types.ErrValidatorAlreadyOffline()
	}

	validator.Online = false

	err = k.SetValidator(ctx, validator)
	if err != nil {
		return nil, types.ErrInternal(err.Error())
	}
	k.DeleteValidatorByPowerIndex(ctx, validator)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(msg.ValidatorAddress).String()),
		sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
