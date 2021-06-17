package swap

import (
	"bitbucket.org/decimalteam/go-node/x/swap/internal/types"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"math/big"
	"runtime/debug"
	"strconv"
	"strings"
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
		case types.MsgSwapInitialize:
			return handleMsgSwapInitialize(ctx, keeper, msg)
		case types.MsgRedeemV2:
			return handleMsgRedeemV2(ctx, keeper, msg)
		case types.MsgChainActivate:
			return handleMsgChainActivate(ctx, keeper, msg)
		case types.MsgChainDeactivate:
			return handleMsgChainDeactivate(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgHTLT(ctx sdk.Context, k Keeper, msg types.MsgHTLT) (*sdk.Result, error) {
	if k.HasSwap(ctx, msg.HashedSecret) {
		return nil, types.ErrSwapAlreadyExist(
			hex.EncodeToString(msg.HashedSecret[:]))
	}

	if msg.TransferType == types.TransferTypeOut {
		ok, err := k.CheckBalance(ctx, msg.From, msg.Amount)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, sdkerrors.ErrInsufficientFunds
		}
	} else if msg.TransferType == types.TransferTypeIn {
		ok, err := k.CheckPoolFunds(ctx, msg.Amount)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, sdkerrors.ErrInsufficientFunds
		}
	}

	var lockedTime int64
	if msg.TransferType == types.TransferTypeOut {
		lockedTime = ctx.BlockTime().Add(k.LockedTimeOut(ctx)).UnixNano()
	} else if msg.TransferType == types.TransferTypeIn {
		lockedTime = ctx.BlockTime().Add(k.LockedTimeIn(ctx)).UnixNano()
	}

	swap := types.NewSwap(
		msg.TransferType,
		msg.HashedSecret,
		msg.From,
		msg.Recipient,
		msg.Amount,
		uint64(lockedTime),
	)

	k.SetSwap(ctx, swap)

	if msg.TransferType == types.TransferTypeOut {
		err := k.LockFunds(ctx, swap.From, swap.Amount)
		if err != nil {
			return nil, err
		}
	} else if msg.TransferType == types.TransferTypeIn {
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

	if swap.Redeemed {
		return nil, types.ErrAlreadyRedeemed()
	}

	if ctx.BlockTime().UnixNano() >= int64(swap.Timestamp) {
		return nil, types.ErrExpired()
	}

	if getHash(msg.Secret) != swap.HashedSecret {
		return nil, types.ErrWrongSecret()
	}

	var err error
	if swap.TransferType == types.TransferTypeIn {
		var recipientAddr sdk.AccAddress
		if swap.Recipient == "" {
			recipientAddr = msg.From
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

func handleMsgSwapInitialize(ctx sdk.Context, k Keeper, msg types.MsgSwapInitialize) (*sdk.Result, error) {
	if !k.HasChain(ctx, msg.DestChain) {
		return nil, types.ErrChainNotExist(strconv.Itoa(msg.DestChain))
	}
	if !k.HasChain(ctx, msg.FromChain) {
		return nil, types.ErrChainNotExist(strconv.Itoa(msg.FromChain))
	}

	funds := sdk.NewCoins(sdk.NewCoin(strings.ToLower(msg.TokenSymbol), msg.Amount))

	ok, err := k.CheckBalance(ctx, msg.From, funds)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, sdkerrors.ErrInsufficientFunds
	}

	err = k.LockFunds(ctx, msg.From, funds)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, msg.From.String()),
			sdk.NewAttribute(types.AttributeKeyDestChain, strconv.Itoa(msg.DestChain)),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyTransactionNumber, msg.TransactionNumber),
			sdk.NewAttribute(types.AttributeKeyTokenName, msg.TokenName),
			sdk.NewAttribute(types.AttributeKeyTokenSymbol, msg.TokenSymbol),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRedeemV2(ctx sdk.Context, k Keeper, msg types.MsgRedeemV2) (*sdk.Result, error) {
	hash, err := types.GetHash(msg.TransactionNumber, msg.TokenName, msg.TokenSymbol, msg.Amount, msg.Recipient, msg.DestChain)
	if err != nil {
		return nil, err
	}

	if k.HasSwapV2(ctx, hash) {
		return nil, fmt.Errorf("swap already redeemed")
	}

	R := big.NewInt(0)
	R.SetBytes(msg.R[:])

	S := big.NewInt(0)
	S.SetBytes(msg.S[:])

	address, err := types.Ecrecover(hash, R, S, sdk.NewInt(int64(msg.V)).BigInt())
	if err != nil {
		return nil, err
	}

	if !address.Equals(types.SwapServiceAddress()) {
		return nil, fmt.Errorf("invalid ecrecover address")
	}

	k.SetSwapV2(ctx, hash)

	funds := sdk.NewCoins(sdk.NewCoin(strings.ToLower(msg.TokenSymbol), msg.Amount))

	ok, err := k.CheckPoolFunds(ctx, funds)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, sdkerrors.ErrInsufficientFunds
	}

	err = k.UnlockFunds(ctx, address, funds)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeKeyFrom, msg.From),
			sdk.NewAttribute(types.AttributeKeyDestChain, strconv.Itoa(msg.DestChain)),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyTransactionNumber, msg.TransactionNumber),
			sdk.NewAttribute(types.AttributeKeyTokenName, msg.TokenName),
			sdk.NewAttribute(types.AttributeKeyTokenSymbol, msg.TokenSymbol),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgChainActivate(ctx sdk.Context, k Keeper, msg types.MsgChainActivate) (*sdk.Result, error) {
	chain, found := k.GetChain(ctx, msg.ChainNumber)
	if found {
		chain.Active = true
	} else {
		chain = types.NewChain(msg.ChainName, true)
	}

	k.SetChain(ctx, msg.ChainNumber, chain)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgChainDeactivate(ctx sdk.Context, k Keeper, msg types.MsgChainDeactivate) (*sdk.Result, error) {
	chain, found := k.GetChain(ctx, msg.ChainNumber)
	if !found {
		return nil, types.ErrChainNotExist(strconv.Itoa(msg.ChainNumber))
	}

	chain.Active = false
	k.SetChain(ctx, msg.ChainNumber, chain)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func getHash(secret []byte) [32]byte {
	return sha256.Sum256(secret)
}
