package swap

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
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
		case types.MsgHTLT:
			return handleMsgHTLT(ctx, keeper, msg)
		case types.MsgRedeem:
			return handleMsgRedeem(ctx, keeper, msg)
		case types.MsgRefund:
			return handleMsgRefund(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgHTLT(ctx sdk.Context, k Keeper, msg types.MsgHTLT) (*sdk.Result, error) {
	if k.HasSwap(ctx, msg.HashedSecret) {
		return nil, types.ErrSwapAlreadyExist(msg.HashedSecret)
	}

	ok, err := k.CheckBalance(ctx, msg.From, msg.Amount)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, sdkerrors.ErrInsufficientFunds
	}

	swap := types.NewSwap(
		msg.TransferType,
		msg.HashedSecret,
		msg.From,
		msg.Recipient,
		msg.Amount,
		uint64(ctx.BlockTime().Add(k.LockedTime(ctx)).UnixNano()),
	)

	k.SetSwap(ctx, swap)

	if msg.TransferType == types.TransferTypeOut {
		err = k.LockFunds(ctx, swap.From, swap.Amount)
		if err != nil {
			return nil, err
		}
	} else if msg.TransferType == types.TransferTypeIn {
		recipient, err := sdk.AccAddressFromBech32(swap.Recipient)
		if err != nil {
			return nil, sdkerrors.ErrInvalidAddress
		}
		err = k.LockFunds(ctx, recipient, swap.Amount)
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
			sdk.NewAttribute(types.AttributeKeyTimeLocked, time.Unix(0, int64(swap.Timestamp)).String()),
			sdk.NewAttribute(types.AttributeKeyHashedSecret, hex.EncodeToString(swap.HashedSecret[:])),
			sdk.NewAttribute(types.AttributeKeyRecipient, swap.Recipient),
			sdk.NewAttribute(types.AttributeKeyAmount, swap.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyTransferType, swap.TransferType.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRedeem(ctx sdk.Context, k Keeper, msg types.MsgRedeem) (*sdk.Result, error) {
	hash := sha256.Sum256(msg.Secret)

	swap, ok := k.GetSwap(ctx, hash)
	if !ok {
		return nil, types.ErrSwapNotFound()
	}

	if !swap.From.Equals(msg.From) {
		return nil, types.ErrFromFieldNotEqual(msg.From, swap.From)
	}

	if swap.Redeemed {
		return nil, types.ErrAlreadyRedeemed()
	}

	if ctx.BlockTime().UnixNano() >= int64(swap.Timestamp) {
		return nil, types.ErrExpired()
	}

	if getHash(msg.Secret) != swap.HashedSecret {
		return nil, types.ErrWrongSecret()
	}

	if swap.TransferType == types.TransferTypeIn {
		recipient, err := sdk.AccAddressFromBech32(swap.Recipient)
		if err != nil {
			return nil, sdkerrors.ErrInvalidAddress
		}
		err = k.UnlockFunds(ctx, recipient, swap.Amount)
		if err != nil {
			return nil, err
		}
	}

	swap.Redeemed = true
	k.SetSwap(ctx, swap)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
			sdk.NewAttribute(types.AttributeKeySecret, hex.EncodeToString(msg.Secret)),
			sdk.NewAttribute(types.AttributeKeyHashedSecret, hex.EncodeToString(swap.HashedSecret[:])),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRefund(ctx sdk.Context, k Keeper, msg types.MsgRefund) (*sdk.Result, error) {
	swap, ok := k.GetSwap(ctx, msg.HashedSecret)
	if !ok {
		return nil, types.ErrSwapNotFound()
	}

	if ctx.BlockTime().UnixNano() < int64(swap.Timestamp) {
		return nil, types.ErrNotExpired()
	}

	if !swap.From.Equals(msg.From) {
		return nil, types.ErrFromFieldNotEqual(msg.From, swap.From)
	}

	if swap.Refunded {
		return nil, types.ErrAlreadyRefunded()
	}

	err := k.UnlockFunds(ctx, swap.From, swap.Amount)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
			sdk.NewAttribute(types.AttributeKeyHashedSecret, hex.EncodeToString(msg.HashedSecret[:])),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func getHash(secret []byte) [32]byte {
	return sha256.Sum256(secret)
}
