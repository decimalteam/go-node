package coin

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"bitbucket.org/decimalteam/go-node/x/validator/types"
	"bytes"
	"encoding/base64"
	"fmt"
	"math/big"
	"runtime/debug"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"

	"github.com/btcsuite/btcutil/base58"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/utils/helpers"
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
)

var e18 = big.NewFloat(1000000000000000000)

func floatFromInt(amount sdk.Int) float64 {
	bigFloat := big.NewFloat(0)
	bigFloat.SetInt(amount.BigInt())
	bigFloat = bigFloat.Quo(bigFloat, e18)
	float, _ := bigFloat.Float64()
	return float
}

// NewHandler creates an sdk.Handler for all the coin type messages
func NewHandler(k Keeper) sdk.Handler {
	//msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("stacktrace from panic: %s \n%s\n", r, string(debug.Stack()))
			}
		}()
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *MsgCreateCoin:
			return handleMsgCreateCoin(ctx, k, *msg)
		case *MsgUpdateCoin:
			return handleMsgUpdateCoin(ctx, k, *msg)
		case *MsgSendCoin:
			return handleMsgSendCoin(ctx, k, *msg)
		case *MsgMultiSendCoin:
			return handleMsgMultiSendCoin(ctx, k, *msg)
		case *MsgBuyCoin:
			return handleMsgBuyCoin(ctx, k, *msg)
		case *MsgSellCoin:
			return handleMsgSellCoin(ctx, k, *msg, false)
		case *MsgSellAllCoin:
			msgSell := MsgSellCoin{
				Sender:       msg.Sender,
				CoinToSell:   msg.CoinToSell,
				MinCoinToBuy: msg.MinCoinToBuy,
			}
			return handleMsgSellCoin(ctx, k, msgSell, true)
		case *MsgRedeemCheck:
			return handleMsgRedeemCheck(ctx, k, *msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

////////////////////////////////////////////////////////////////
// Creating coins handlers
////////////////////////////////////////////////////////////////

func getCreateCoinCommission(symbol string) sdk.Int {
	switch len(symbol) {
	case 3:
		return sdk.NewInt(1_000_000)
	case 4:
		return sdk.NewInt(100_000)
	case 5:
		return sdk.NewInt(10_000)
	case 6:
		return sdk.NewInt(1000)
	}
	return sdk.NewInt(100)
}

func handleMsgCreateCoin(ctx sdk.Context, k Keeper, msg types2.MsgCreateCoin) (*sdk.Result, error) {
	if msg.InitialReserve.LT(MinCoinReserve(ctx)) {
		return nil, types2.ErrInvalidCoinInitialReserve(ctx)
	}

	var coin = types2.Coin{
		Title:       msg.Title,
		CRR:         msg.ConstantReserveRatio,
		Symbol:      strings.ToLower(msg.Symbol),
		Reserve:     msg.InitialReserve,
		LimitVolume: msg.LimitVolume,
		Volume:      msg.InitialVolume,
		Creator:     msg.Sender,
		Identity:    msg.Identity,
	}

	_, err := k.GetCoin(ctx, strings.ToLower(msg.Symbol))
	if err == nil {
		return nil, types2.ErrCoinAlreadyExist(msg.Symbol)
	}

	commission, feeCoin, err := k.GetCommission(ctx, helpers.BipToPip(getCreateCoinCommission(coin.Symbol)))
	if err != nil {
		return nil, types2.ErrCalculateCommission(err)
	}

	accAddr, err := sdk.AccAddressFromBech32(msg.Sender)

	if err != nil {
		return nil, err
	}

	acc := k.AccountKeeper.GetAccount(ctx, accAddr)
	balance := k.BankKeeper.GetAllBalances(ctx, acc.GetAddress())
	if balance.AmountOf(cliUtils.GetBaseCoin()).LT(msg.InitialReserve) {
		return nil, types2.ErrInsufficientCoinReserve()
	}

	if balance.AmountOf(feeCoin).LT(commission) {
		return nil, types2.ErrInsufficientFundsToPayCommission(commission.String())
	}

	if feeCoin == cliUtils.GetBaseCoin() {
		if balance.AmountOf(cliUtils.GetBaseCoin()).LT(commission.Add(msg.InitialReserve)) {
			return nil, types2.ErrInsufficientFunds(commission.Add(msg.InitialReserve).String(), balance.AmountOf(cliUtils.GetBaseCoin()).String())
		}
	}

	err = k.UpdateBalance(ctx, cliUtils.GetBaseCoin(), msg.InitialReserve.Neg(), acc.GetAddress())
	if err != nil {
		return nil, types2.ErrUpdateBalance(msg.Sender, err.Error())
	}

	k.SetCoin(ctx, coin)
	err = k.UpdateBalance(ctx, coin.Symbol, msg.InitialVolume, acc.GetAddress())
	if err != nil {
		return nil, types2.ErrUpdateBalance(msg.Sender, err.Error())
	}

	err = k.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), acc.GetAddress())
	if err != nil {
		return nil, types2.ErrUpdateBalance(msg.Sender, err.Error())
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types2.AttributeSymbol, coin.Symbol),
		sdk.NewAttribute(types2.AttributeTitle, coin.Title),
		sdk.NewAttribute(types2.AttributeCRR, strconv.FormatUint(uint64(msg.ConstantReserveRatio), 10)),
		sdk.NewAttribute(types2.AttributeInitVolume, msg.InitialVolume.String()),
		sdk.NewAttribute(types2.AttributeInitReserve, msg.InitialReserve.String()),
		sdk.NewAttribute(types2.AttributeLimitVolume, msg.LimitVolume.String()),
		sdk.NewAttribute(types2.AttributeCommissionCreateCoin, sdk.NewCoin(strings.ToLower(feeCoin), commission).String()),
	))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

////////////////////////////////////////////////////////////////
// Updating coin handler
////////////////////////////////////////////////////////////////

func handleMsgUpdateCoin(ctx sdk.Context, k Keeper, msg types2.MsgUpdateCoin) (*sdk.Result, error) {
	coin, err := k.GetCoin(ctx, strings.ToLower(msg.Symbol))
	if err != nil {
		return nil, types2.ErrCoinAlreadyExist(msg.Symbol)
	}

	if coin.Creator != msg.Sender {
		return nil, types2.ErrUpdateOnlyForCreator()
	}

	if coin.Volume.GT(msg.LimitVolume) {
		return nil, types2.ErrLimitVolumeBroken(coin.Volume.String(), msg.LimitVolume.String())
	}

	coin.LimitVolume = msg.LimitVolume
	coin.Identity = msg.Identity

	k.SetCoin(ctx, coin)

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

////////////////////////////////////////////////////////////////
// Transfer coins handlers
////////////////////////////////////////////////////////////////

func handleMsgSendCoin(ctx sdk.Context, k Keeper, msg types2.MsgSendCoin) (*sdk.Result, error) {
	senderaddr, _ := sdk.AccAddressFromBech32(msg.Sender)
	receiveraddr, _ := sdk.AccAddressFromBech32(msg.Receiver)

	_, err := k.GetCoin(ctx, msg.Coin.Denom)
	if err != nil {
		return nil, types.ErrCoinDoesNotExist(msg.Coin.Denom)
	}

	err = k.BankKeeper.SendCoins(ctx, senderaddr, receiveraddr, sdk.Coins{msg.Coin})
	if err != nil {
		return nil, types.ErrInternal(err.Error())
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types2.AttributeCoin, msg.Coin.String()),
		sdk.NewAttribute(types2.AttributeReceiver, msg.Receiver),
	))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgMultiSendCoin(ctx sdk.Context, k Keeper, msg types2.MsgMultiSendCoin) (*sdk.Result, error) {
	senderaddr, _ := sdk.AccAddressFromBech32(msg.Sender)

	for i := range msg.Sends {
		receiveraddr, _ := sdk.AccAddressFromBech32(msg.Sends[i].Receiver)
		_, err := k.GetCoin(ctx, msg.Sends[i].Coin.Denom)
		if err != nil {
			return nil, types.ErrCoinDoesNotExist(msg.Sends[i].Coin.Denom)
		}
		err = k.BankKeeper.SendCoins(ctx, senderaddr, receiveraddr, sdk.Coins{msg.Sends[i].Coin})
		if err != nil {
			return nil, types.ErrInternal(err.Error())
		}

		ctx.EventManager().EmitEvent(sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(types2.AttributeCoin, msg.Sends[i].Coin.String()),
			sdk.NewAttribute(types2.AttributeReceiver, msg.Sends[i].Receiver),
		))
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

////////////////////////////////////////////////////////////////
// Trading coins handlers
////////////////////////////////////////////////////////////////

func handleMsgBuyCoin(ctx sdk.Context, k Keeper, msg types2.MsgBuyCoin) (*sdk.Result, error) {
	// Retrieve buyer account and it's balance of selling coins
	accAddr, err := sdk.AccAddressFromBech32(msg.Sender)

	if err != nil {
		return nil, err
	}

	account := k.AccountKeeper.GetAccount(ctx, accAddr)
	balance := k.BankKeeper.GetAllBalances(ctx, account.GetAddress()).AmountOf(strings.ToLower(msg.MaxCoinToSell.Denom))

	// Retrieve the coin requested to buy
	coinToBuy, err := k.GetCoin(ctx, msg.CoinToBuy.Denom)
	if err != nil {
		return nil, types2.ErrCoinDoesNotExist(msg.CoinToBuy.Denom)
	}
	if coinToBuy.Symbol != msg.CoinToBuy.Denom {
		return nil, types2.ErrRetrievedAnotherCoin(msg.CoinToBuy.Denom, coinToBuy.Symbol)
	}

	// Retrieve the coin requested to sell
	coinToSell, err := k.GetCoin(ctx, msg.MaxCoinToSell.Denom)
	if err != nil {
		return nil, types2.ErrCoinDoesNotExist(msg.MaxCoinToSell.Denom)
	}
	if coinToSell.Symbol != msg.MaxCoinToSell.Denom {
		return nil, types2.ErrRetrievedAnotherCoin(msg.MaxCoinToSell.Denom, coinToSell.Symbol)
	}

	// Ensure supply limit of the coin to buy does not overflow
	if !coinToBuy.IsBase() {
		if coinToBuy.Volume.Add(msg.CoinToBuy.Amount).GT(coinToBuy.LimitVolume) {
			return nil, types2.ErrTxBreaksVolumeLimit(coinToBuy.Volume.Add(msg.CoinToBuy.Amount).String(), coinToBuy.LimitVolume.String())
		}
	}

	// Calculate amount of sell coins which buyer will receive
	amountToBuy, amountToSell, amountInBaseCoin := msg.CoinToBuy.Amount, sdk.ZeroInt(), sdk.ZeroInt()
	switch {
	case coinToSell.IsBase():
		// Buyer buys custom coin for base coin
		amountToSell = formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountToBuy)
		amountInBaseCoin = amountToSell
	case coinToBuy.IsBase():
		// Buyer buys base coin for custom coin
		if msg.CoinToBuy.Amount.GT(coinToSell.Reserve) {
			return nil, types2.ErrInsufficientCoinReserve()
		}

		amountToSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountToBuy)
		amountInBaseCoin = amountToBuy
	default:
		// Buyer buys custom coin for custom coin

		amountInBaseCoin = formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountToBuy)
		if amountInBaseCoin.GT(coinToSell.Reserve) {
			return nil, types2.ErrInsufficientCoinReserve()
		}

		amountToSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountInBaseCoin)
	}

	// Ensure maximum amount of coins to sell (price guard)
	if amountToSell.GT(msg.MaxCoinToSell.Amount) {
		return nil, types2.ErrMaximumValueToSellReached(msg.MaxCoinToSell.Amount.String(), amountToSell.String())
	}

	// Ensure reserve of the coin to sell does not underflow
	if !coinToSell.IsBase() {
		if coinToSell.Reserve.Sub(amountInBaseCoin).LT(types2.MinCoinReserve(ctx)) {
			return nil, types2.ErrTxBreaksMinReserveRule(ctx, amountInBaseCoin.String())
		}
	}

	// Ensure that buyer account holds enough coins to sell
	if balance.LT(amountToSell) {
		return nil, types2.ErrInsufficientFunds(amountToSell.String(), balance.String())
	}

	senderaddr, err := sdk.AccAddressFromBech32(msg.Sender)

	if err != nil {
		return nil, err
	}

	// Update buyer account balances
	err = k.UpdateBalance(ctx, msg.MaxCoinToSell.Denom, amountToSell.Neg(), senderaddr)
	if err != nil {
		return nil, types2.ErrUpdateBalance(account.GetAddress().String(), err.Error())
	}
	err = k.UpdateBalance(ctx, msg.CoinToBuy.Denom, amountToBuy, senderaddr)
	if err != nil {
		return nil, types2.ErrUpdateBalance(account.GetAddress().String(), err.Error())
	}

	// Update coins
	if !coinToSell.IsBase() {
		k.UpdateCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountInBaseCoin), coinToSell.Volume.Sub(amountToSell))
	}
	if !coinToBuy.IsBase() {
		k.UpdateCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountInBaseCoin), coinToBuy.Volume.Add(amountToBuy))
	}

	// Emit transaction events
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types2.AttributeCoinToBuy, sdk.NewCoin(msg.CoinToBuy.Denom, amountToBuy).String()),
		sdk.NewAttribute(types2.AttributeCoinToSell, sdk.NewCoin(msg.MaxCoinToSell.Denom, amountToSell).String()),
		sdk.NewAttribute(types2.AttributeAmountInBaseCoin, amountInBaseCoin.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgSellCoin(ctx sdk.Context, k Keeper, msg types2.MsgSellCoin, sellAll bool) (*sdk.Result, error) {
	// Retrieve seller account and it's balance of selling coins
	senderaddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	account := k.AccountKeeper.GetAccount(ctx, senderaddr)
	balance := k.BankKeeper.GetAllBalances(ctx, account.GetAddress()).AmountOf(strings.ToLower(msg.CoinToSell.Denom))

	// Fill amount to sell in case of MsgSellAll
	if sellAll {
		msg.CoinToSell.Amount = balance
	}

	// Retrieve the coin requested to sell
	coinToSell, err := k.GetCoin(ctx, msg.CoinToSell.Denom)
	if err != nil {
		return nil, types2.ErrCoinDoesNotExist(msg.CoinToSell.Denom)
	}
	if coinToSell.Symbol != msg.CoinToSell.Denom {
		return nil, types2.ErrRetrievedAnotherCoin(msg.CoinToSell.Denom, coinToSell.Symbol)
	}

	// Retrieve the coin requested to buy
	coinToBuy, err := k.GetCoin(ctx, msg.MinCoinToBuy.Denom)
	if err != nil {
		return nil, types2.ErrCoinDoesNotExist(msg.MinCoinToBuy.Denom)
	}
	if coinToBuy.Symbol != msg.MinCoinToBuy.Denom {
		return nil, types2.ErrRetrievedAnotherCoin(msg.MinCoinToBuy.Denom, coinToBuy.Symbol)
	}

	// Ensure that seller account holds enough coins to sell
	if balance.LT(msg.CoinToSell.Amount) {
		return nil, types2.ErrInsufficientFunds(msg.CoinToSell.String(), balance.String())
	}

	// Calculate amount of buy coins which seller will receive
	amountToSell, amountToBuy, amountInBaseCoin := msg.CoinToSell.Amount, sdk.ZeroInt(), sdk.ZeroInt()
	switch {
	case coinToBuy.IsBase():
		// Seller sells custom coin for base coin
		amountToBuy = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountToSell)
		amountInBaseCoin = amountToBuy
	case coinToSell.IsBase():
		// Seller sells base coin for custom coin
		amountToBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountToSell)
		amountInBaseCoin = amountToSell
	default:
		// Seller sells custom coin for custom coin
		amountInBaseCoin = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountToSell)
		amountToBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountInBaseCoin)
	}

	// Ensure minimum amount of coins to buy (price guard)
	if amountToBuy.LT(msg.MinCoinToBuy.Amount) {
		return nil, types2.ErrMinimumValueToBuyReached(amountToBuy.String(), msg.MinCoinToBuy.Amount.String())
	}

	// Ensure reserve of the coin to sell does not underflow
	if !coinToSell.IsBase() {
		if coinToSell.Reserve.Sub(amountInBaseCoin).LT(types2.MinCoinReserve(ctx)) {
			return nil, types2.ErrTxBreaksMinReserveRule(ctx, amountInBaseCoin.String())
		}
	}

	// Ensure supply limit of the coin to buy does not overflow
	if !coinToBuy.IsBase() {
		if coinToBuy.Volume.Add(amountToBuy).GT(coinToBuy.LimitVolume) {
			return nil, types2.ErrTxBreaksVolumeLimit(coinToBuy.Volume.Add(amountToBuy).String(), coinToBuy.LimitVolume.String())
		}
	}

	// Update seller account balances
	err = k.UpdateBalance(ctx, msg.CoinToSell.Denom, amountToSell.Neg(), senderaddr)
	if err != nil {
		return nil, types2.ErrUpdateBalance(msg.Sender, err.Error())
	}
	err = k.UpdateBalance(ctx, msg.MinCoinToBuy.Denom, amountToBuy, senderaddr)
	if err != nil {
		return nil, types2.ErrUpdateBalance(msg.Sender, err.Error())
	}

	// Update coins
	if !coinToSell.IsBase() {
		k.UpdateCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountInBaseCoin), coinToSell.Volume.Sub(amountToSell))
	}
	if !coinToBuy.IsBase() {
		k.UpdateCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountInBaseCoin), coinToBuy.Volume.Add(amountToBuy))
	}

	// Emit transaction events
	// eventType := types.EventTypeSellCoin
	// if sellAll {
	// 	eventType = types.EventTypeSellAllCoin
	// }
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types2.AttributeCoinToSell, sdk.NewCoin(msg.CoinToSell.Denom, amountToSell).String()),
		sdk.NewAttribute(types2.AttributeCoinToBuy, sdk.NewCoin(msg.MinCoinToBuy.Denom, amountToBuy).String()),
		sdk.NewAttribute(types2.AttributeAmountInBaseCoin, amountInBaseCoin.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

////////////////////////////////////////////////////////////////
// Redeem check handler
////////////////////////////////////////////////////////////////

func handleMsgRedeemCheck(ctx sdk.Context, k Keeper, msg types2.MsgRedeemCheck) (*sdk.Result, error) {
	// Decode provided check from base58 format to raw bytes
	checkBytes := base58.Decode(msg.Check)
	if len(checkBytes) == 0 {
		msgError := "unable to decode check from base58"
		return nil, sdkerrors.New(types2.DefaultCodespace, types2.InvalidCheck, msgError)
	}

	// Parse provided check from raw bytes to ensure it is valid
	check, err := types2.ParseCheck(checkBytes)
	if err != nil {
		msgError := fmt.Sprintf("unable to parse check: %s", err.Error())
		return nil, sdkerrors.New(types2.DefaultCodespace, types2.InvalidCheck, msgError)
	}

	// Decode provided proof from base64 format to raw bytes
	proof, err := base64.StdEncoding.DecodeString(msg.Proof)
	if err != nil {
		msgError := "unable to decode proof from base64"
		return nil, sdkerrors.New(types2.DefaultCodespace, types2.InvalidProof, msgError)
	}

	// Recover issuer address from check signature
	issuer, err := check.Sender()
	if err != nil {
		errMsg := fmt.Sprintf("unable to recover check issuer address: %s", err.Error())
		return nil, sdkerrors.New(types2.DefaultCodespace, types2.InvalidCheck, errMsg)
	}

	account := k.AccountKeeper.GetAccount(ctx, issuer)
	balance := k.BankKeeper.GetAllBalances(ctx, account.GetAddress()).AmountOf(strings.ToLower(check.Coin))

	// Retrieve the coin specified in the check
	coin, err := k.GetCoin(ctx, check.Coin)
	if err != nil {
		return nil, types2.ErrCoinDoesNotExist(check.Coin)
	}
	if coin.Symbol != check.Coin {
		return nil, types2.ErrRetrievedAnotherCoin(check.Coin, coin.Symbol)
	}

	feeCoin := cliUtils.GetBaseCoin()
	commission := helpers.UnitToPip(sdk.NewIntFromUint64(30))

	// Ensure that check issuer account holds enough coins
	amount := sdk.NewIntFromBigInt(check.Amount)
	if balance.LT(amount) {
		return nil, types2.ErrInsufficientFunds(check.Amount.String()+coin.Symbol, balance.String()+coin.Symbol)
	}

	if feeCoin == strings.ToLower(check.Coin) {
		if balance.LT(amount.Add(commission)) {
			return nil, types2.ErrInsufficientFunds(amount.String()+coin.Symbol, balance.String()+coin.Symbol)
		}
	}

	// Ensure the proper chain ID is specified in the check
	if check.ChainID != ctx.ChainID() {
		errMsg := fmt.Sprintf("wanted chain ID %s, but check is issued for chain with ID %s", ctx.ChainID(), check.ChainID)
		return nil, sdkerrors.New(types2.DefaultCodespace, types2.InvalidChainID, errMsg)
	}

	// Ensure nonce length
	if len(check.Nonce) > 16 {
		errMsg := "nonce is too big (should be up to 16 bytes)"
		return nil, sdkerrors.New(types2.DefaultCodespace, types2.InvalidNonce, errMsg)
	}

	// Check block number
	if check.DueBlock < uint64(ctx.BlockHeight()) {
		errMsg := fmt.Sprintf("check was expired at block %d", check.DueBlock)
		return nil, sdkerrors.New(types2.DefaultCodespace, types2.CheckExpired, errMsg)
	}

	// Ensure check is not redeemed yet
	if k.IsCheckRedeemed(ctx, check) {
		errMsg := "check was redeemed already"
		return nil, sdkerrors.New(types2.DefaultCodespace, types2.CheckRedeemed, errMsg)
	}

	// Recover public key from check lock
	publicKeyA, err := check.LockPubKey()
	if err != nil {
		msgError := fmt.Sprintf("unable to recover lock public key from check: %s", err.Error())
		return nil, sdkerrors.New(types2.DefaultCodespace, types2.InvalidCheck, msgError)
	}

	// Prepare bytes used to recover public key from provided proof
	senderAddressHash := make([]byte, 32)
	hw := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hw, []interface{}{
		msg.Sender,
	})
	if err != nil {
		msgError := fmt.Sprintf("unable to RLP encode check sender address: %s", err.Error())
		return nil, sdkerrors.New(types2.DefaultCodespace, types2.InvalidCheck, msgError)
	}
	hw.Sum(senderAddressHash[:0])

	// Recover public key from provided proof
	publicKeyB, err := crypto.Ecrecover(senderAddressHash[:], proof)

	// Compare both public keys to ensure provided proof is correct
	if !bytes.Equal(publicKeyA, publicKeyB) {
		msgError := fmt.Sprintf("provided proof is invalid %s", err.Error())
		return nil, sdkerrors.New(types2.DefaultCodespace, types2.InvalidProof, msgError)
	}

	// Set check redeemed
	k.SetCheckRedeemed(ctx, check)

	senderaddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	// Update accounts balances
	err = k.UpdateBalance(ctx, feeCoin, commission.Neg(), issuer)
	if err != nil {
		return nil, types2.ErrInsufficientFundsToPayCommission(commission.String())
	}
	err = k.BankKeeper.SendCoins(ctx, issuer, senderaddr, sdk.Coins{sdk.NewCoin(coin.Symbol, amount)})
	if err != nil {
		return nil, sdkerrors.New(types2.DefaultCodespace, 6, err.Error())
	}

	// Emit event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types2.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		sdk.NewAttribute(types2.AttributeIssuer, issuer.String()),
		sdk.NewAttribute(types2.AttributeCoin, sdk.NewCoin(check.Coin, sdk.NewIntFromBigInt(check.Amount)).String()),
		sdk.NewAttribute(types2.AttributeNonce, new(big.Int).SetBytes(check.Nonce).String()),
		sdk.NewAttribute(types2.AttributeDueBlock, strconv.FormatUint(check.DueBlock, 10)),
		sdk.NewAttribute(types2.AttributeCommissionRedeemCheck, sdk.NewCoin(feeCoin, commission).String()),
	))

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
