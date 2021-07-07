package swap

import (
	types2 "bitbucket.org/decimalteam/go-node/x/swap/types"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"runtime/debug"
	"time"
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
		case *types2.MsgHTLT:
			return handleMsgHTLT(ctx, keeper, *msg)
		case *types2.MsgRedeem:
			return handleMsgRedeem(ctx, keeper, *msg)
		case *types2.MsgRefund:
			return handleMsgRefund(ctx, keeper, *msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types2.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgHTLT(ctx sdk.Context, k Keeper, msg types2.MsgHTLT) (*sdk.Result, error) {
	if k.HasSwap(ctx, msg.HashedSecret) {
		return nil, types2.ErrSwapAlreadyExist(msg.HashedSecret)
	}

	if msg.TransferType == types2.TransferTypeOut {
		fromaddr, err := sdk.AccAddressFromBech32(msg.From)
		if err != nil {
			return nil, err
		}

		ok, err := k.CheckBalance(ctx, fromaddr, msg.Amount)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, sdkerrors.ErrInsufficientFunds
		}
	} else if msg.TransferType == types2.TransferTypeIn {
		ok, err := k.CheckPoolFunds(ctx, msg.Amount)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, sdkerrors.ErrInsufficientFunds
		}
	}

	var lockedTime int64
	if msg.TransferType == types2.TransferTypeOut {
		lockedTime = ctx.BlockTime().Add(k.LockedTimeOut(ctx)).UnixNano()
	} else if msg.TransferType == types2.TransferTypeIn {
		lockedTime = ctx.BlockTime().Add(k.LockedTimeIn(ctx)).UnixNano()
	}

	swap := types2.NewSwap(
		msg.TransferType,
		msg.HashedSecret,
		msg.From,
		msg.Recipient,
		msg.Amount,
		uint64(lockedTime),
	)

	k.SetSwap(ctx, swap)

	fromaddr, err := sdk.AccAddressFromBech32(swap.From)
	if err != nil {
		return nil, err
	}

	if msg.TransferType == types2.TransferTypeOut {
		err = k.LockFunds(ctx, fromaddr, swap.Amount)
		if err != nil {
			return nil, err
		}
	} else if msg.TransferType == types2.TransferTypeIn {
		if swap.Recipient != "" {
			_, err := sdk.AccAddressFromBech32(swap.Recipient)
			if err != nil {
				return nil, sdkerrors.ErrInvalidAddress
			}
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, fromaddr.String()),
			sdk.NewAttribute(types2.AttributeKeyTimeLocked, time.Unix(0, int64(swap.Timestamp)).String()),
			sdk.NewAttribute(types2.AttributeKeyHashedSecret, hex.EncodeToString(swap.HashedSecret[:])),
			sdk.NewAttribute(types2.AttributeKeyRecipient, swap.Recipient),
			sdk.NewAttribute(types2.AttributeKeyAmount, swap.Amount.String()),
			sdk.NewAttribute(types2.AttributeKeyTransferType, swap.TransferType.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRedeem(ctx sdk.Context, k Keeper, msg types2.MsgRedeem) (*sdk.Result, error) {
	hash := sha256.Sum256(msg.Secret)

	swap, ok := k.GetSwap(ctx, hash)
	if !ok {
		return nil, types2.ErrSwapNotFound()
	}

	if swap.Redeemed {
		return nil, types2.ErrAlreadyRedeemed()
	}

	if ctx.BlockTime().UnixNano() >= int64(swap.Timestamp) {
		return nil, types2.ErrExpired()
	}

	if getHash(msg.Secret) != *swap.HashedSecret {
		return nil, types2.ErrWrongSecret()
	}

	fromaddr, err := sdk.AccAddressFromBech32(swap.From)
	if err != nil {
		return nil, err
	}

	if swap.TransferType == types2.TransferTypeIn {
		var recipientAddr sdk.AccAddress
		if swap.Recipient == "" {
			recipientAddr = fromaddr
		} else {
			recipientAddr, err = sdk.AccAddressFromBech32(swap.Recipient)
			if err != nil {
				return nil, sdkerrors.ErrInvalidAddress
			}
		}
		err = k.UnlockFunds(ctx, recipientAddr, swap.Amount)
		if err != nil {
			return nil, err
		}
	}

	swap.Redeemed = true
	k.SetSwap(ctx, swap)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, fromaddr.String()),
			sdk.NewAttribute(types2.AttributeKeySecret, hex.EncodeToString(msg.Secret)),
			sdk.NewAttribute(types2.AttributeKeyHashedSecret, hex.EncodeToString(swap.HashedSecret[:])),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRefund(ctx sdk.Context, k Keeper, msg types2.MsgRefund) (*sdk.Result, error) {
	swap, ok := k.GetSwap(ctx, msg.HashedSecret)
	if !ok {
		return nil, types2.ErrSwapNotFound()
	}

	if ctx.BlockTime().UnixNano() < int64(swap.Timestamp) {
		return nil, types2.ErrNotExpired()
	}

	if swap.Refunded {
		return nil, types2.ErrAlreadyRefunded()
	}

	fromaddr, err := sdk.AccAddressFromBech32(swap.From)
	if err != nil {
		return nil, err
	}

	err = k.UnlockFunds(ctx, fromaddr, swap.Amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, fromaddr.String()),
			sdk.NewAttribute(types2.AttributeKeyHashedSecret, hex.EncodeToString(msg.HashedSecret[:])),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func getHash(secret []byte) [32]byte {
	return sha256.Sum256(secret)
}
