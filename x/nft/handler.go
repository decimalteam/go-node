package nft

import (
	keeper2 "bitbucket.org/decimalteam/go-node/x/nft/keeper"
	types2 "bitbucket.org/decimalteam/go-node/x/nft/types"
	"fmt"
	"runtime/debug"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GenericHandler routes the messages to the handlers
func GenericHandler(k keeper2.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("stacktrace from panic: %s \n%s\n", r, string(debug.Stack()))
			}
		}()
		switch msg := msg.(type) {
		case types2.MsgTransferNFT:
			return HandleMsgTransferNFT(ctx, msg, k)
		case types2.MsgEditNFTMetadata:
			return HandleMsgEditNFTMetadata(ctx, msg, k)
		case types2.MsgMintNFT:
			return HandleMsgMintNFT(ctx, msg, k)
		case types2.MsgBurnNFT:
			return HandleMsgBurnNFT(ctx, msg, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("unrecognized nft message type: %T", msg))
		}
	}
}

// HandleMsgTransferNFT handler for MsgTransferNFT
func HandleMsgTransferNFT(ctx sdk.Context, msg types2.MsgTransferNFT, k keeper2.Keeper,
) (*sdk.Result, error) {
	nft, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err != nil {
		return nil, err
	}

	nft, err = types2.TransferNFT(nft, msg.Sender, msg.Recipient, msg.Quantity)
	if err != nil {
		return nil, err
	}
	fmt.Println(nft)

	collection, found := k.GetCollection(ctx, msg.Denom)
	if !found {
		return nil, ErrUnknownCollection
	}
	collection.NFTs.Update(msg.ID, nft)
	k.SetCollection(ctx, msg.Denom, collection)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types2.EventTypeTransfer,
			sdk.NewAttribute(types2.AttributeKeyRecipient, msg.Recipient.String()),
			sdk.NewAttribute(types2.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types2.AttributeKeyNFTID, msg.ID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// HandleMsgEditNFTMetadata handler for MsgEditNFTMetadata
func HandleMsgEditNFTMetadata(ctx sdk.Context, msg types2.MsgEditNFTMetadata, k keeper2.Keeper,
) (*sdk.Result, error) {
	nft, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err != nil {
		return nil, err
	}

	// update NFT
	nft.EditMetadata(msg.TokenURI)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types2.EventTypeEditNFTMetadata,
			sdk.NewAttribute(types2.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types2.AttributeKeyNFTID, msg.ID),
			sdk.NewAttribute(types2.AttributeKeyNFTTokenURI, msg.TokenURI),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// HandleMsgMintNFT handles MsgMintNFT
func HandleMsgMintNFT(ctx sdk.Context, msg types2.MsgMintNFT, k keeper2.Keeper,
) (*sdk.Result, error) {
	nft, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err == nil {
		if !nft.GetCreator().Equals(msg.Sender) {
			return nil, ErrNotAllowedMint
		}
		if !nft.GetAllowMint() {
			return nil, ErrNotAllowedMint
		}
	}
	nft = types2.NewBaseNFT(msg.ID, msg.Sender, msg.Recipient, msg.TokenURI, msg.Quantity, msg.Reserve, msg.AllowMint)
	err = k.MintNFT(ctx, msg.Denom, nft)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types2.EventTypeMintNFT,
			sdk.NewAttribute(types2.AttributeKeyRecipient, msg.Recipient.String()),
			sdk.NewAttribute(types2.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types2.AttributeKeyNFTID, msg.ID),
			sdk.NewAttribute(types2.AttributeKeyNFTTokenURI, msg.TokenURI),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// HandleMsgBurnNFT handles MsgBurnNFT
func HandleMsgBurnNFT(ctx sdk.Context, msg types2.MsgBurnNFT, k keeper2.Keeper,
) (*sdk.Result, error) {
	nft, err := k.GetNFT(ctx, msg.Denom, msg.ID)
	if err != nil {
		return nil, err
	}

	if !nft.GetCreator().Equals(msg.Sender) {
		return nil, ErrNotAllowedBurn
	}

	// remove  NFT
	err = k.DeleteNFT(ctx, msg.Denom, msg.ID, msg.Quantity)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types2.EventTypeBurnNFT,
			sdk.NewAttribute(types2.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types2.AttributeKeyNFTID, msg.ID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
