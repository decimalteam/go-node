package validator

import (
	"bitbucket.org/decimalteam/go-node/utils/helpers"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmstrings "github.com/tendermint/tendermint/libs/strings"
	tmtypes "github.com/tendermint/tendermint/types"
	"strings"
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
)

const (
	declareCandidateFee = 10000
	editCandidateFee    = 10000
	delegateFee         = 200
	unbondFee           = 200
	setOnlineFee        = 100
	setOfflineFee       = 100
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
	commission, feeCoin, err := k.CoinKeeper.GetCommission(ctx, helpers.UnitToPip(declareCandidateFee))
	if err != nil {
		return nil, types.ErrCalculateCommission(err)
	}

	acc := k.AccountKeeper.GetAccount(ctx, sdk.AccAddress(msg.ValidatorAddr))
	balance := acc.GetCoins()

	if balance.AmountOf(k.BondDenom(ctx)).LT(commission) {
		return nil, types.ErrInsufficientCoinToPayCommission(commission.String())
	}

	if msg.Stake.Denom == k.BondDenom(ctx) {
		if balance.AmountOf(k.BondDenom(ctx)).LT(commission.Add(msg.Stake.Amount)) {
			return nil, types.ErrInsufficientFunds(commission.Add(msg.Stake.Amount).String())
		}
	}

	// check to see if the pubkey or sender has been registered before
	if _, err := k.GetValidator(ctx, msg.ValidatorAddr); err == nil {
		return nil, types.ErrValidatorOwnerExists(k.Codespace())
	}

	if _, err := k.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(msg.PubKey)); err == nil {
		return nil, types.ErrValidatorPubKeyExists(k.Codespace())
	}

	if ctx.ConsensusParams() != nil {
		tmPubKey := tmtypes.TM2PB.PubKey(msg.PubKey)
		if !tmstrings.StringInSlice(tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes) {
			return nil, sdkerrors.Wrapf(
				types.ErrValidatorPubKeyTypeNotSupported(k.Codespace()),
				"got: %s, valid: %s", tmPubKey.Type, ctx.ConsensusParams().Validator.PubKeyTypes,
			)
		}
	}

	val := types.NewValidator(msg.ValidatorAddr, msg.PubKey, msg.Commission, msg.RewardAddr, msg.Description)
	err = k.SetValidator(ctx, val)
	if err != nil {
		return nil, types.ErrInvalidStruct(k.Codespace())
	}
	k.SetValidatorByConsAddr(ctx, val)
	k.SetNewValidatorByPowerIndex(ctx, val)

	k.AfterValidatorCreated(ctx, val.ValAddress)

	err = k.Delegate(ctx, sdk.AccAddress(msg.ValidatorAddr), msg.Stake, types.Unbonded, val, true)
	if err != nil {
		return nil, sdkerrors.New(k.Codespace(), types.CodeInvalidDelegation, err.Error())
	}

	err = k.CoinKeeper.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), sdk.AccAddress(msg.ValidatorAddr))
	if err != nil {
		return nil, types.ErrUpdateBalance(err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeclareCandidate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Stake.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Stake.Denom),
			sdk.NewAttribute(types.AttributeKeyPubKey, msg.PubKey.Address().String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(msg.ValidatorAddr).String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgDelegate(ctx sdk.Context, k Keeper, msg types.MsgDelegate) (*sdk.Result, error) {
	commission, feeCoin, err := k.CoinKeeper.GetCommission(ctx, helpers.UnitToPip(delegateFee))
	if err != nil {
		return nil, types.ErrCalculateCommission(err)
	}

	acc := k.AccountKeeper.GetAccount(ctx, msg.DelegatorAddress)
	balance := acc.GetCoins()

	if balance.AmountOf(k.BondDenom(ctx)).LT(commission) {
		return nil, types.ErrInsufficientCoinToPayCommission(commission.String())
	}

	if msg.Amount.Denom == k.BondDenom(ctx) {
		if balance.AmountOf(k.BondDenom(ctx)).LT(commission.Add(msg.Amount.Amount)) {
			return nil, types.ErrInsufficientFunds(commission.Add(msg.Amount.Amount).String())
		}
	}

	val, err := k.GetValidator(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, types.ErrNoValidatorFound(k.Codespace())
	}

	if !k.IsDelegatorStakeSufficient(ctx, val, msg.DelegatorAddress, msg.Amount) {
		return nil, types.ErrDelegatorStakeIsTooLow(k.Codespace())
	}

	err = k.Delegate(ctx, msg.DelegatorAddress, msg.Amount, types.Unbonded, val, true)
	if err != nil {
		return nil, sdkerrors.New(k.Codespace(), 1, err.Error())
	}

	err = k.CoinKeeper.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), msg.DelegatorAddress)
	if err != nil {
		return nil, types.ErrUpdateBalance(err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegate,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Amount.Denom),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgUnbond(ctx sdk.Context, k Keeper, msg types.MsgUnbond) (*sdk.Result, error) {
	commission, feeCoin, err := k.CoinKeeper.GetCommission(ctx, helpers.UnitToPip(unbondFee))
	if err != nil {
		return nil, types.ErrCalculateCommission(err)
	}

	acc := k.AccountKeeper.GetAccount(ctx, msg.DelegatorAddress)
	balance := acc.GetCoins()

	if balance.AmountOf(k.BondDenom(ctx)).LT(commission) {
		return nil, types.ErrInsufficientCoinToPayCommission(commission.String())
	}

	completionTime, err := k.Undelegate(ctx, msg.DelegatorAddress, msg.ValidatorAddress, msg.Amount)
	if err != nil {
		return nil, sdkerrors.New(k.Codespace(), 1, err.Error())
	}

	completionTimeBz := types.ModuleCdc.MustMarshalBinaryLengthPrefixed(completionTime)

	err = k.CoinKeeper.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), msg.DelegatorAddress)
	if err != nil {
		return nil, types.ErrUpdateBalance(err)
	}

	return &sdk.Result{Data: completionTimeBz, Events: sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbond,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Amount.Denom),
			sdk.NewAttribute(types.AttributeKeyDelegator, msg.DelegatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress.String()),
		),
	}}, nil
}

func handleMsgEditCandidate(ctx sdk.Context, k Keeper, msg types.MsgEditCandidate) (*sdk.Result, error) {
	commission, feeCoin, err := k.CoinKeeper.GetCommission(ctx, helpers.UnitToPip(editCandidateFee))
	if err != nil {
		return nil, types.ErrCalculateCommission(err)
	}

	acc := k.AccountKeeper.GetAccount(ctx, sdk.AccAddress(msg.ValidatorAddress))
	balance := acc.GetCoins()

	if balance.AmountOf(k.BondDenom(ctx)).LT(commission) {
		return nil, types.ErrInsufficientCoinToPayCommission(commission.String())
	}

	validator, err := k.GetValidatorByConsAddr(ctx, sdk.ConsAddress(msg.PubKey.Address()))
	if err != nil {
		return nil, types.ErrNoValidatorFound(k.Codespace())
	}

	validator.ValAddress = msg.ValidatorAddress
	validator.RewardAddress = msg.RewardAddress

	k.SetValidatorByConsAddr(ctx, validator)
	err = k.SetValidator(ctx, validator)
	if err != nil {
		return nil, sdkerrors.New(k.Codespace(), 1, err.Error())
	}

	err = k.CoinKeeper.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), sdk.AccAddress(msg.ValidatorAddress))
	if err != nil {
		return nil, types.ErrUpdateBalance(err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditCandidate,
			sdk.NewAttribute(types.AttributeKeyRewardAddress, msg.RewardAddress.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(validator.ValAddress).String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSetOnline(ctx sdk.Context, k Keeper, msg types.MsgSetOnline) (*sdk.Result, error) {
	commission, feeCoin, err := k.CoinKeeper.GetCommission(ctx, helpers.UnitToPip(setOnlineFee))
	if err != nil {
		return nil, types.ErrCalculateCommission(err)
	}

	acc := k.AccountKeeper.GetAccount(ctx, sdk.AccAddress(msg.ValidatorAddress))
	balance := acc.GetCoins()

	if balance.AmountOf(k.BondDenom(ctx)).LT(commission) {
		return nil, types.ErrInsufficientCoinToPayCommission(commission.String())
	}

	validator, err := k.GetValidator(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, types.ErrNoValidatorFound(k.Codespace())
	}

	if validator.Online {
		if !validator.Jailed {
			return nil, sdkerrors.New(k.Codespace(), 1, "Validator already online")
		}
	}

	validator.Online = true
	validator.Jailed = false
	err = k.SetValidator(ctx, validator)
	if err != nil {
		return nil, sdkerrors.New(k.Codespace(), 1, err.Error())
	}
	k.SetValidatorByPowerIndex(ctx, validator)

	err = k.CoinKeeper.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), sdk.AccAddress(msg.ValidatorAddress))
	if err != nil {
		return nil, types.ErrUpdateBalance(err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetOnline,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(msg.ValidatorAddress).String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSetOffline(ctx sdk.Context, k Keeper, msg types.MsgSetOffline) (*sdk.Result, error) {
	commission, feeCoin, err := k.CoinKeeper.GetCommission(ctx, helpers.UnitToPip(setOfflineFee))
	if err != nil {
		return nil, types.ErrCalculateCommission(err)
	}

	acc := k.AccountKeeper.GetAccount(ctx, sdk.AccAddress(msg.ValidatorAddress))
	balance := acc.GetCoins()

	if balance.AmountOf(k.BondDenom(ctx)).LT(commission) {
		return nil, types.ErrInsufficientCoinToPayCommission(commission.String())
	}

	validator, err := k.GetValidator(ctx, msg.ValidatorAddress)
	if err != nil {
		return nil, types.ErrNoValidatorFound(k.Codespace())
	}

	if !validator.Online {
		return nil, sdkerrors.New(k.Codespace(), 1, "Validator already offline")
	}

	validator.Online = false

	err = k.SetValidator(ctx, validator)
	if err != nil {
		return nil, sdkerrors.New(k.Codespace(), 1, err.Error())
	}
	k.DeleteValidatorByPowerIndex(ctx, validator)

	err = k.CoinKeeper.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), sdk.AccAddress(msg.ValidatorAddress))
	if err != nil {
		return nil, types.ErrUpdateBalance(err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetOffline,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sdk.AccAddress(msg.ValidatorAddress).String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
