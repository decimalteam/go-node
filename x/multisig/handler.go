package multisig

import (
	types2 "bitbucket.org/decimalteam/go-node/x/multisig/types"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/utils/helpers"
)

// NewHandler creates an sdk.Handler for all the multisig type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("stacktrace from panic: %s \n%s\n", r, string(debug.Stack()))
			}
		}()
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *MsgCreateWallet:
			return handleMsgCreateWallet(ctx, k, *msg)
		case *MsgCreateTransaction:
			return handleMsgCreateTransaction(ctx, k, *msg)
		case *MsgSignTransaction:
			return handleMsgSignTransaction(ctx, k, *msg, true)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types2.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreateWallet(ctx sdk.Context, keeper Keeper, msg MsgCreateWallet) (*sdk.Result, error) {
	// Create new multisig wallet
	wallet, err := NewWallet(msg.Owners, msg.Weights, msg.Threshold, ctx.TxBytes())
	if err != nil {
		msgError := fmt.Sprintf("Unable to create new multi-signature wallet: %s", err.Error())
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Ensure multisig wallet with the address does not exist
	existingWallet := keeper.GetWallet(ctx, wallet.Address)

	walladdr, err := sdk.AccAddressFromBech32(existingWallet.Address)
	if err != nil {
		return nil, err
	}

	if !walladdr.Empty() {
		msgError := fmt.Sprintf("Multi-signature wallet with address %s already exists", existingWallet.Address)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	accAddr, err := sdk.AccAddressFromBech32(wallet.Address)
	if err != nil {
		return nil, err
	}

	// Ensure account with multisig address does not exist
	existingAccount := keeper.AccountKeeper.GetAccount(ctx, accAddr)
	if existingAccount != nil && !existingAccount.GetAddress().Empty() {
		msgError := fmt.Sprintf("Account with address %s already exists", existingAccount.GetAddress())
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Save created multisig wallet to the KVStore
	keeper.SetWallet(ctx, *wallet)

	// Emit transaction events
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types2.AttributeKeyOwners, helpers.JoinAccAddresses(msg.Owners)),
		sdk.NewAttribute(types2.AttributeKeyWeights, helpers.JoinUints(msg.Weights)),
		sdk.NewAttribute(types2.AttributeKeyThreshold, strconv.FormatUint(uint64(msg.Threshold), 10)),
		sdk.NewAttribute(types2.AttributeKeyWallet, walladdr.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgCreateTransaction(ctx sdk.Context, keeper Keeper, msg MsgCreateTransaction) (*sdk.Result, error) {
	// Retrieve multisig wallet from the KVStore
	wallet := keeper.GetWallet(ctx, msg.Wallet)
	if len(wallet.Address) == 0 {
		msgError := fmt.Sprintf("No registered multi-signature wallet with address %s", msg.Wallet)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	accAddr, err := sdk.AccAddressFromBech32(wallet.Address)
	if err != nil {
		return nil, err
	}

	// Retrieve coins hold on the multisig wallet
	var walletCoins sdk.Coins
	if walletAccount := keeper.AccountKeeper.GetAccount(ctx, accAddr); walletAccount != nil {
		walletCoins = keeper.BankKeeper.GetAllBalances(ctx, walletAccount.GetAddress())
	} else {
		walletCoins = sdk.NewCoins()
	}

	// Ensure there are enough coins on the multisig wallet
	for _, coin := range msg.Coins {
		if walletCoins.AmountOf(strings.ToLower(coin.Denom)).LT(coin.Amount) {
			return nil, types2.ErrInsufficientFunds(coin.String())
		}
	}

	// Create new multisig transaction
	transaction, err := NewTransaction(
		msg.Wallet,
		msg.Receiver,
		msg.Coins,
		wallet.Owners,
		ctx.BlockHeight(),
		ctx.TxBytes(),
	)
	if err != nil {
		msgError := "Unable to create multi-signature transaction"
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Save created multisig transaction to the KVStore
	keeper.SetTransaction(ctx, *transaction)

	// Sign created multisig transaction by the creator
	signEvents, err := handleMsgSignTransaction(ctx, keeper, MsgSignTransaction{
		Sender: msg.Sender,
		TxID:   transaction.ID,
	}, false)
	if err != nil {
		msgError := fmt.Sprintf("Unable to sign created multi-signature transaction with ID %s by it's creator: %s", transaction.ID, err.Error())
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Emit transaction events
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types2.AttributeKeyWallet, msg.Wallet),
		sdk.NewAttribute(types2.AttributeKeyReceiver, msg.Receiver),
		sdk.NewAttribute(types2.AttributeKeyCoins, msg.Coins.String()),
		sdk.NewAttribute(types2.AttributeKeyTransaction, transaction.ID),
	))
	ctx.EventManager().EmitEvents(signEvents.GetEvents())

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgSignTransaction(ctx sdk.Context, keeper Keeper, msg MsgSignTransaction, emitEvents bool) (*sdk.Result, error) {
	senderacc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Retrieve multisig transaction from the KVStore
	transaction := keeper.GetTransaction(ctx, msg.TxID)

	txwalletaddr, err := sdk.AccAddressFromBech32(transaction.Wallet)
	if err != nil {
		return nil, err
	}

	if txwalletaddr.Empty() {
		msgError := fmt.Sprintf("No registered multi-signature transaction with ID %s", msg.TxID)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Retrieve multisig wallet from the KVStore
	wallet := keeper.GetWallet(ctx, transaction.Wallet)

	walladrr, err := sdk.AccAddressFromBech32(wallet.Address)
	if err != nil {
		return nil, err
	}

	if walladrr.Empty() {
		msgError := fmt.Sprintf("No registered multi-signature wallet with address %s", transaction.Wallet)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Calculate current weight of signatures
	confirmations := uint64(0)
	for i, c := 0, len(wallet.Owners); i < c; i++ {
		if len(transaction.Signers[i]) != 0 {
			confirmations += wallet.Weights[i]
		}
	}

	// Ensure current weight of signatures is not enough
	if confirmations >= wallet.Threshold {
		msgError := fmt.Sprintf("Multi-signature transaction already has enough signatures (%d >= %d)", confirmations, wallet.Threshold)
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
	}

	// Append the signature to the multisig transaction
	weight := uint64(0)
	for i, c := 0, len(wallet.Owners); i < c; i++ {
		owneraddr, err := sdk.AccAddressFromBech32(wallet.Owners[i])
		if err != nil {
			return nil, err
		}

		if owneraddr.Equals(senderacc) {
			txsigneraddr, err := sdk.AccAddressFromBech32(transaction.Signers[i])

			if err != nil {
				return nil, err
			}

			if !txsigneraddr.Empty() {
				msgError := fmt.Sprintf("Unable to sign multi-signature transaction since signer with address %s is already signed it", msg.Sender)
				return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
			}
			weight = wallet.Weights[i]
			confirmations += weight
			transaction.Signers[i] = msg.Sender
			break
		}
		if i == c-1 {
			msgError := fmt.Sprintf("Unable to sign multi-signature transaction since signer with address %s is not an owner of the wallet", msg.Sender)
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
		}
	}

	// Save updated multisig transaction to the KVStore
	keeper.SetTransaction(ctx, transaction)

	txreceiveraddr, err := sdk.AccAddressFromBech32(transaction.Receiver)
	if err != nil {
		return nil, err
	}

	// Check if new weight of signatures is enough to perform multisig transaction
	confirmed := confirmations >= wallet.Threshold
	if confirmed {
		// Perform transaction
		err = keeper.BankKeeper.SendCoins(ctx, walladrr, txreceiveraddr, transaction.Coins)
		if err != nil {
			msgError := fmt.Sprintf("Unable to perform multi-signature transaction %s: %s", transaction.ID, err.Error())
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, msgError)
		}
	}

	// Emit transaction events
	events := sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, senderacc.String()),
			sdk.NewAttribute(types2.AttributeKeyWallet, walladrr.String()),
			sdk.NewAttribute(types2.AttributeKeyTransaction, msg.TxID),
			sdk.NewAttribute(types2.AttributeKeySignerWeight, strconv.FormatUint(weight, 10)),
			sdk.NewAttribute(types2.AttributeKeyConfirmations, strconv.FormatUint(confirmations, 10)),
			sdk.NewAttribute(types2.AttributeKeyConfirmed, strconv.FormatBool(confirmed)),
		),
	}
	if !emitEvents {
		return &sdk.Result{Events: events.ToABCIEvents()}, nil
	}
	ctx.EventManager().EmitEvents(events)
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
