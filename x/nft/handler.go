package nft

import (
	"fmt"
	"runtime/debug"
	"strconv"

	"bitbucket.org/decimalteam/go-node/utils/updates"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/x/nft/internal/keeper"
	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
)

// GenericHandler routes the messages to the handlers
func GenericHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("stacktrace from panic: %s \n%s\n", r, string(debug.Stack()))
			}
		}()
		switch msg := msg.(type) {
		case types.MsgTransferNFT:
			return HandleMsgTransferNFT(ctx, msg, k)
		case types.MsgEditNFTMetadata:
			return HandleMsgEditNFTMetadata(ctx, msg, k)
		case types.MsgMintNFT:
			return HandleMsgMintNFT(ctx, msg, k)
		case types.MsgBurnNFT:
			return HandleMsgBurnNFT(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("unrecognized nft message type: %T", msg))
		}
	}
}

// HandleMsgTransferNFT handler for MsgTransferNFT
func HandleMsgTransferNFT(ctx sdk.Context, msg types.MsgTransferNFT, k keeper.Keeper,
) (*sdk.Result, error) {
	nft, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err != nil {
		return nil, err
	}

	nft, err = types.TransferNFT(nft, msg.Sender, msg.Recipient, msg.SubTokenIDs)
	if err != nil {
		return nil, err
	}

	collection, found := k.GetCollection(ctx, msg.Denom)
	if !found {
		return nil, ErrUnknownCollection(msg.Denom)
	}

	collection.NFTs, _ = collection.NFTs.Update(msg.ID, nft)

	k.SetCollection(ctx, msg.Denom, collection)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTransfer,
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.ID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// HandleMsgEditNFTMetadata handler for MsgEditNFTMetadata
func HandleMsgEditNFTMetadata(ctx sdk.Context, msg types.MsgEditNFTMetadata, k keeper.Keeper,
) (*sdk.Result, error) {
	nft, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err != nil {
		return nil, err
	}

	if !nft.GetCreator().Equals(msg.Sender) {
		return nil, ErrNotAllowedMint()
	}

	// update NFT
	nft = nft.EditMetadata(msg.TokenURI)
	err = k.UpdateNFT(ctx, msg.Denom, nft)

	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditNFTMetadata,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.ID),
			sdk.NewAttribute(types.AttributeKeyNFTTokenURI, msg.TokenURI),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// HandleMsgMintNFT handles MsgMintNFT
func HandleMsgMintNFT(ctx sdk.Context, msg types.MsgMintNFT, k keeper.Keeper,
) (*sdk.Result, error) {
	nft, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err == nil {
		if !nft.GetCreator().Equals(msg.Sender) || !nft.GetAllowMint() {
			return nil, ErrNotAllowedMint()
		}
	} else {
		if k.ExistTokenURI(ctx, msg.TokenURI) {
			return nil, ErrNotUniqueTokenURI()
		}
		if k.ExistTokenID(ctx, msg.ID) {
			return nil, ErrNotUniqueTokenID()
		}
		if ctx.BlockHeight() >= updates.Update2Block {
			if msg.Reserve.LT(types.NewMinReserve2) {
				return nil, types.ErrInvalidReserve(msg.Reserve.String())
			}
		} else if ctx.BlockHeight() >= updates.Update1Block {
			if msg.Reserve.LT(types.NewMinReserve) {
				return nil, types.ErrInvalidReserve(msg.Reserve.String())
			}
		} else {
			if msg.Reserve.LT(types.MinReserve) {
				return nil, types.ErrInvalidReserve(msg.Reserve.String())
			}
		}
	}

	lastSubTokenID, err := k.MintNFT(ctx, msg.Denom, msg.ID, msg.Reserve, msg.Quantity, msg.Sender, msg.Recipient, msg.TokenURI, msg.AllowMint)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMintNFT,
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.ID),
			sdk.NewAttribute(types.AttributeKeyNFTTokenURI, msg.TokenURI),
			sdk.NewAttribute(types.AttributeKeySubTokenIDStartRange, strconv.FormatInt(lastSubTokenID-msg.Quantity.Int64(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// HandleMsgBurnNFT handles MsgBurnNFT
func HandleMsgBurnNFT(ctx sdk.Context, msg types.MsgBurnNFT, k keeper.Keeper,
) (*sdk.Result, error) {
	nft, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err != nil {
		return nil, err
	}

	if !nft.GetCreator().Equals(msg.Sender) {
		return nil, ErrNotAllowedBurn()
	}

	// remove NFT
	err = k.DeleteNFT(ctx, msg.Denom, msg.ID, msg.SubTokenIDs)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBurnNFT,
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.ID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
