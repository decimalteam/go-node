package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/codec"
	"runtime/debug"
	"time"

	tmstrings "github.com/tendermint/tendermint/libs/strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/x/validator/types"
)

// NewHandler creates an sdk.Handler for all the validator type messages
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("stacktrace from panic: %s \n%s\n", r, string(debug.Stack()))
			}
		}()
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgDeclareCandidate:
			return handleMsgDeclareCandidate(ctx, keeper, *msg)
		case *types.MsgDelegate:
			return handleMsgDelegate(ctx, keeper, *msg)
		case *types.MsgDelegateNFT:
			return handleMsgDelegateNFT(ctx, keeper, *msg)
		case *types.MsgUnbond:
			return handleMsgUnbond(ctx, keeper, *msg)
		case *types.MsgUnbondNFT:
			return handleMsgUnbondNFT(ctx, keeper, *msg)
		case *types.MsgEditCandidate:
			return handleMsgEditCandidate(ctx, keeper, *msg)
		case *types.MsgSetOnline:
			return handleMsgSetOnline(ctx, keeper, *msg)
		case *types.MsgSetOffline:
			return handleMsgSetOffline(ctx, keeper, *msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgDeclareCandidate(ctx sdk.Context, k Keeper, msg types.MsgDeclareCandidate) (*sdk.Result, error) {
	validatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	// check to see if the pubkey or sender has been registered before
	if _, err := k.GetValidator(ctx, validatorAddr); err == nil {
		return nil, types.ErrValidatorOwnerExists()
	}

	if _, err := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(msg.PubKey)); err == nil {
		return nil, types.ErrValidatorPubKeyExists()
	}

	if ctx.ConsensusParams() != nil {
		tmPubKey, err := codec.ToTmProtoPublicKey(msg.PubKey)

		if err != nil {
			return nil, err
		}

		if !tmstrings.StringInSlice(tmPubKey.String(), ctx.ConsensusParams().Validator.PubKeyTypes) {
			return nil, sdkerrors.Wrapf(
				types.ErrValidatorPubKeyTypeNotSupported(),
				"got: %s, valid: %s", tmPubKey.String(), ctx.ConsensusParams().Validator.PubKeyTypes,
			)
		}
	}

	val := types.NewValidator(validatorAddr.String(), msg.PubKey, msg.Commission, msg.RewardAddr, msg.Description)
	err = k.SetValidator(ctx, val)
	if err != nil {
		return nil, types.ErrInvalidStruct()
	}
	k.SetValidatorByConsAddr(ctx, val)
	k.SetNewValidatorByPowerIndex(ctx, validatorAddr, val)

	k.AfterValidatorCreated(ctx, validatorAddr)

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
		sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
		sdk.NewAttribute(types.AttributeKeyCoin, msg.Stake.String()),
		sdk.NewAttribute(types.AttributeKeyPubKey, msg.PubKey.Address().String()),
		sdk.NewAttribute(types.AttributeKeyCommission, msg.Commission.String()),
		sdk.NewAttribute(types.AttributeKeyDescriptionMoniker, msg.Description.Moniker),
		sdk.NewAttribute(types.AttributeKeyDescriptionIdentity, msg.Description.Identity),
		sdk.NewAttribute(types.AttributeKeyDescriptionWebsite, msg.Description.Website),
		sdk.NewAttribute(types.AttributeKeyDescriptionSecurityContact, msg.Description.SecurityContact),
		sdk.NewAttribute(types.AttributeKeyDescriptionDetails, msg.Description.Details),
	))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDelegate(ctx sdk.Context, k Keeper, msg types.MsgDelegate) (*sdk.Result, error) {
	validatoraddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}
	delegatoraddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	val, err := k.GetValidator(ctx, validatoraddr)
	if err != nil {
		return nil, types.ErrNoValidatorFound()
	}

	ok, err := k.IsDelegatorStakeSufficient(ctx, validatoraddr, delegatoraddr, msg.Coin)
	if err != nil {
		return nil, types.ErrCoinDoesNotExist(msg.Coin.Denom)
	}
	if !ok {
		return nil, types.ErrDelegatorStakeIsTooLow()
	}

	err = k.Delegate(ctx, delegatoraddr, msg.Coin, types.Unbonded, val, true)
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
		sdk.NewAttribute(sdk.AttributeKeySender, delegatoraddr.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, validatoraddr.String()),
		sdk.NewAttribute(types.AttributeKeyCoin, msg.Coin.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgDelegateNFT(ctx sdk.Context, k Keeper, msg types.MsgDelegateNFT) (*sdk.Result, error) {
	delegatorAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil ,err
	}
	validatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil ,err
	}

	val, err := k.GetValidator(ctx, validatorAddr)
	if err != nil {
		return nil, types.ErrNoValidatorFound()
	}

	err = k.DelegateNFT(ctx, delegatorAddr, msg.TokenID, msg.Denom, msg.Quantity, val)
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
		sdk.NewAttribute(sdk.AttributeKeySender, delegatorAddr.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUnbondNFT(ctx sdk.Context, k Keeper, msg types.MsgUnbondNFT) (*sdk.Result, error) {
	delegatorAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil ,err
	}
	validatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil ,err
	}

	completionTime, err := k.UndelegateNFT(ctx, delegatorAddr, validatorAddr, msg.TokenID, msg.Denom, msg.Quantity)
	if err != nil {
		e := sdkerrors.Error{}
		if errors.As(err, &e) {
			return nil, e
		} else {
			return nil, types.ErrInternal(err.Error())
		}
	}

	completionTimeBz, err := json.Marshal(completionTime)
	if err != nil {
		panic(err)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, delegatorAddr.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
		sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
	))

	return &sdk.Result{Data: completionTimeBz, Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUnbond(ctx sdk.Context, k Keeper, msg types.MsgUnbond) (*sdk.Result, error) {
	delegatorAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil ,err
	}
	validatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil ,err
	}

	completionTime, err := k.Undelegate(ctx, delegatorAddr, validatorAddr, msg.Coin)
	if err != nil {
		e := sdkerrors.Error{}
		if errors.As(err, &e) {
			return nil, e
		} else {
			return nil, types.ErrInternal(err.Error())
		}
	}

	//completionTimeBz := types.ModuleCdc.MustMarshalLengthPrefixed(completionTime)
	completionTimeBz, err := json.Marshal(completionTime)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, delegatorAddr.String()),
		sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
		sdk.NewAttribute(types.AttributeKeyCoin, msg.Coin.String()),
		sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
	))

	return &sdk.Result{Data: completionTimeBz, Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgEditCandidate(ctx sdk.Context, k Keeper, msg types.MsgEditCandidate) (*sdk.Result, error) {
	var validator types.Validator

	validatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil ,err
	}

	validator, err = k.GetValidator(ctx, validatorAddr)
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
		sdk.NewAttribute(types.AttributeKeyRewardAddress, msg.RewardAddress),
		sdk.NewAttribute(types.AttributeKeyDescriptionMoniker, msg.Description.Moniker),
		sdk.NewAttribute(types.AttributeKeyDescriptionDetails, msg.Description.Details),
		sdk.NewAttribute(types.AttributeKeyDescriptionIdentity, msg.Description.Identity),
		sdk.NewAttribute(types.AttributeKeyDescriptionWebsite, msg.Description.Website),
		sdk.NewAttribute(types.AttributeKeyDescriptionSecurityContact, msg.Description.SecurityContact),
	))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgSetOnline(ctx sdk.Context, k Keeper, msg types.MsgSetOnline) (*sdk.Result, error) {
	validatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil ,err
	}
	validator, err := k.GetValidator(ctx, validatorAddr)
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
		sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgSetOffline(ctx sdk.Context, k Keeper, msg types.MsgSetOffline) (*sdk.Result, error) {
	validatorAddr, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil ,err
	}

	validator, err := k.GetValidator(ctx, validatorAddr)
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
		sdk.NewAttribute(types.AttributeKeyValidator, validatorAddr.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
