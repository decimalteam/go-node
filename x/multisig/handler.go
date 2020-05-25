package multisig

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/utils/helpers"
	"bitbucket.org/decimalteam/go-node/x/multisig/internal/types"
)

// NewHandler creates an sdk.Handler for all the multisig type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		// ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateWallet:
			return handleMsgCreateWallet(ctx, k, msg)
		case MsgCreateTransaction:
			return handleMsgCreateTransaction(ctx, k, msg)
		case MsgSignTransaction:
			return handleMsgSignTransaction(ctx, k, msg, true)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreateWallet(ctx sdk.Context, keeper Keeper, msg MsgCreateWallet) (*sdk.Result, error) {

	// Create new multisig wallet
	wallet, err := NewWallet(msg.Owners, msg.Weights, msg.Threshold)
	if err != nil {
		msgError := fmt.Sprintf("Unable to create new multi-signature wallet: %s", err.Error())
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Ensure multisig wallet with the address does not exist
	existingWallet := keeper.GetWallet(ctx, wallet.Address.String())
	if !existingWallet.Address.Empty() {
		msgError := fmt.Sprintf("Multi-signature wallet with address %s already exists", existingWallet.Address)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Ensure account with multisig address does not exist
	existingAccount := keeper.AccountKeeper.GetAccount(ctx, wallet.Address)
	if existingAccount != nil && !existingAccount.GetAddress().Empty() {
		msgError := fmt.Sprintf("Account with address %s already exists", existingAccount.GetAddress())
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Save created multisig wallet to the KVStore
	keeper.SetWallet(ctx, *wallet)

	// Emit transaction events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeCreateWallet),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator.String()),
			sdk.NewAttribute(types.AttributeKeyOwners, helpers.JoinAccAddresses(msg.Owners)),
			sdk.NewAttribute(types.AttributeKeyWeights, helpers.JoinUints(msg.Weights)),
			sdk.NewAttribute(types.AttributeKeyThreshold, strconv.FormatUint(uint64(msg.Threshold), 10)),
			sdk.NewAttribute(types.AttributeKeyWallet, wallet.Address.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgCreateTransaction(ctx sdk.Context, keeper Keeper, msg MsgCreateTransaction) (*sdk.Result, error) {

	// Retrieve multisig wallet from the KVStore
	wallet := keeper.GetWallet(ctx, msg.Wallet.String())
	if wallet.Address.Empty() {
		msgError := fmt.Sprintf("No registered multi-signature wallet with address %s", msg.Wallet)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// TODO: Check wallet has enough coins

	// Create new multisig transaction
	transaction, err := NewTransaction(
		msg.Wallet,
		msg.Receiver,
		msg.Coins,
		make([]sdk.AccAddress, len(wallet.Owners)),
		ctx.BlockHeight(),
	)

	// Save created multisig transaction to the KVStore
	keeper.SetTransaction(ctx, *transaction)

	// Sign created multisig transaction by the creator
	signEvents, err := handleMsgSignTransaction(ctx, keeper, MsgSignTransaction{
		Signer: msg.Creator,
		TxID:   transaction.ID,
	}, false)
	if err != nil {
		msgError := fmt.Sprintf("Unable to sign created multi-signature transaction with ID %s by it's creator: %s", transaction.ID, err.Error())
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Emit transaction events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeCreateTransaction),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator.String()),
			sdk.NewAttribute(types.AttributeKeyWallet, msg.Wallet.String()),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.Receiver.String()),
			sdk.NewAttribute(types.AttributeKeyCoins, msg.Coins.String()),
			sdk.NewAttribute(types.AttributeKeyTransaction, transaction.ID),
		),
	})
	ctx.EventManager().EmitEvents(signEvents.Events)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSignTransaction(ctx sdk.Context, keeper Keeper, msg MsgSignTransaction, emitEvents bool) (*sdk.Result, error) {

	// Retrieve multisig transaction from the KVStore
	transaction := keeper.GetTransaction(ctx, msg.TxID)
	if transaction.Wallet.Empty() {
		msgError := fmt.Sprintf("No registered multi-signature transaction with ID %s", msg.TxID)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Retrieve multisig wallet from the KVStore
	wallet := keeper.GetWallet(ctx, transaction.Wallet.String())
	if wallet.Address.Empty() {
		msgError := fmt.Sprintf("No registered multi-signature wallet with address %s", transaction.Wallet)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Calculate current weight of signatures
	transactionWeight := uint(0)
	for i, c := 0, len(wallet.Owners); i < c; i++ {
		if !transaction.Signers[i].Empty() {
			transactionWeight += wallet.Weights[i]
		}
	}

	// Ensure current weight of signatures is not enough
	if transactionWeight >= wallet.Threshold {
		msgError := fmt.Sprintf("Multi-signature transaction already has enough signatures (%d >= %d)", transactionWeight, wallet.Threshold)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Append the signature to the multisig transaction
	weight := uint(0)
	for i, c := 0, len(wallet.Owners); i < c; i++ {
		if wallet.Owners[i].Equals(msg.Signer) {
			if !transaction.Signers[i].Empty() {
				msgError := fmt.Sprintf("Unable to sign multi-signature transaction since signer with address %s is already signed it", msg.Signer)
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
			}
			weight = wallet.Weights[i]
			transactionWeight += weight
			transaction.Signers[i] = msg.Signer
			break
		}
		if i == c-1 {
			msgError := fmt.Sprintf("Unable to sign multi-signature transaction since signer with address %s is not an owner of the wallet", msg.Signer)
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
		}
	}

	// Save updated multisig transaction to the KVStore
	keeper.SetTransaction(ctx, transaction)

	// Check if new weight of signatures is enough to perform multisig transaction
	approved := (transactionWeight >= wallet.Threshold)
	if approved {
		// Perform transaction
		err := keeper.BankKeeper.SendCoins(ctx, wallet.Address, transaction.Receiver, transaction.Coins)
		if err != nil {
			msgError := fmt.Sprintf("Unable to perform multi-signature transaction %s: %s", transaction.ID, err.Error())
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
		}
	}

	// Emit transaction events
	events := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeSignTransaction),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer.String()),
			sdk.NewAttribute(types.AttributeKeySignerWeight, strconv.FormatUint(uint64(weight), 10)),
			sdk.NewAttribute(types.AttributeKeyTransaction, msg.TxID),
			sdk.NewAttribute(types.AttributeKeyTransactionWeight, strconv.FormatUint(uint64(transactionWeight), 10)),
			sdk.NewAttribute(types.AttributeKeyApproved, strconv.FormatBool(approved)),
		),
	}
	if !emitEvents {
		return &sdk.Result{Events: events}, nil
	}
	ctx.EventManager().EmitEvents(events)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
