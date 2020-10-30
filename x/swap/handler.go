package swap

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	"crypto/sha256"
	"errors"
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
			return handleMsgClaim(ctx, keeper, msg)
		case types.MsgRefund:
			return handleMsgRefund(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgHTLT(ctx sdk.Context, k Keeper, msg types.MsgHTLT) (*sdk.Result, error) {
	if k.HasSwap(ctx, msg.Hash) {
		return nil, types.ErrSwapAlreadyExist(msg.Hash)
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
		msg.Hash,
		msg.From,
		msg.Recipient,
		msg.Amount,
		uint64(ctx.BlockTime().UnixNano()),
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

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgClaim(ctx sdk.Context, k Keeper, msg types.MsgRedeem) (*sdk.Result, error) {
	hash := sha256.Sum256(msg.Secret[:])

	swap, ok := k.GetSwap(ctx, hash)
	if !ok {
		return nil, types.ErrSwapNotFound()
	}

	if !swap.From.Equals(msg.From) {
		return nil, types.ErrFromFieldNotEqual(msg.From, swap.From)
	}

	if swap.Claimed {
		return nil, errors.New("already claimed")
	}

	if ctx.BlockTime().Sub(time.Unix(0, int64(swap.Timestamp))) >= k.LockedTime(ctx) {
		return nil, errors.New("expired")
	}

	if sha256.Sum256(msg.Secret[:]) != swap.Hash {
		return nil, errors.New("wrong secret")
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

	swap.Claimed = true
	k.SetSwap(ctx, swap)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRefund(ctx sdk.Context, k Keeper, msg types.MsgRefund) (*sdk.Result, error) {
	swap, ok := k.GetSwap(ctx, msg.Hash)
	if !ok {
		return nil, types.ErrSwapNotFound()
	}

	if ctx.BlockTime().Sub(time.Unix(0, int64(swap.Timestamp))) < k.LockedTime(ctx) {
		return nil, errors.New("swap not expired")
	}

	if !swap.From.Equals(msg.From) {
		return nil, errors.New("'from' field not equal")
	}

	err := k.UnlockFunds(ctx, swap.From, swap.Amount)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func getHash(secret [32]byte) [32]byte {
	return sha256.Sum256(secret[:])
}
