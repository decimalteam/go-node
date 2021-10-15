package coin

import (
	"bitbucket.org/decimalteam/go-node/utils/updates"
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
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
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

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("stacktrace from panic: %s \n%s\n", r, string(debug.Stack()))
			}
		}()
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateCoin:
			return handleMsgCreateCoin(ctx, k, msg)
		case types.MsgUpdateCoin:
			return handleMsgUpdateCoin(ctx, k, msg)
		case types.MsgSendCoin:
			return handleMsgSendCoin(ctx, k, msg)
		case types.MsgMultiSendCoin:
			return handleMsgMultiSendCoin(ctx, k, msg)
		case types.MsgBuyCoin:
			return handleMsgBuyCoin(ctx, k, msg)
		case types.MsgSellCoin:
			return handleMsgSellCoin(ctx, k, msg, false)
		case types.MsgSellAllCoin:
			msgSell := MsgSellCoin{
				Sender:       msg.Sender,
				CoinToSell:   msg.CoinToSell,
				MinCoinToBuy: msg.MinCoinToBuy,
			}
			return handleMsgSellCoin(ctx, k, msgSell, true)
		case types.MsgRedeemCheck:
			return handleMsgRedeemCheck(ctx, k, msg)
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

func handleMsgCreateCoin(ctx sdk.Context, k Keeper, msg types.MsgCreateCoin) (*sdk.Result, error) {
	if msg.InitialReserve.LT(MinCoinReserve(ctx)) {
		return nil, types.ErrInvalidCoinInitialReserve(MinCoinReserve(ctx).String())
	}

	var coin = types.Coin{
		Title:       msg.Title,
		CRR:         msg.ConstantReserveRatio,
		Symbol:      strings.ToLower(msg.Symbol),
		Reserve:     msg.InitialReserve,
		LimitVolume: msg.LimitVolume,
		Volume:      msg.InitialVolume,
	}

	_, err := k.GetCoin(ctx, strings.ToLower(msg.Symbol))
	if err == nil {
		return nil, types.ErrCoinAlreadyExist(msg.Symbol)
	}

	if ctx.BlockHeight() >= updates.Update6Block {
		coin.Creator = msg.Sender
	}

	if ctx.BlockHeight() >= updates.Update6Block {
		coin.Identity = msg.Identity
	}

	commission, feeCoin, err := k.GetCommission(ctx, helpers.BipToPip(getCreateCoinCommission(coin.Symbol)))
	if err != nil {
		return nil, types.ErrCalculateCommission(err.Error())
	}

	acc := k.AccountKeeper.GetAccount(ctx, msg.Sender)
	balance := acc.GetCoins()
	if balance.AmountOf(k.GetBaseCoin(ctx)).LT(msg.InitialReserve) {
		return nil, types.ErrInsufficientCoinReserve()
	}

	if balance.AmountOf(feeCoin).LT(commission) {
		return nil, types.ErrInsufficientFundsToPayCommission(commission.String())
	}

	if feeCoin == k.GetBaseCoin(ctx) {
		if balance.AmountOf(k.GetBaseCoin(ctx)).LT(commission.Add(msg.InitialReserve)) {
			return nil, types.ErrInsufficientFunds(commission.Add(msg.InitialReserve).String(), balance.AmountOf(k.GetBaseCoin(ctx)).String())
		}
	}

	err = k.UpdateBalance(ctx, k.GetBaseCoin(ctx), msg.InitialReserve.Neg(), msg.Sender)
	if err != nil {
		return nil, types.ErrUpdateBalance(msg.Sender.String(), err.Error())
	}

	k.SetCoin(ctx, coin)
	err = k.UpdateBalance(ctx, coin.Symbol, msg.InitialVolume, msg.Sender)
	if err != nil {
		return nil, types.ErrUpdateBalance(msg.Sender.String(), err.Error())
	}

	err = k.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), msg.Sender)
	if err != nil {
		return nil, types.ErrUpdateBalance(msg.Sender.String(), err.Error())
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		sdk.NewAttribute(types.AttributeSymbol, coin.Symbol),
		sdk.NewAttribute(types.AttributeTitle, coin.Title),
		sdk.NewAttribute(types.AttributeCRR, strconv.FormatUint(uint64(msg.ConstantReserveRatio), 10)),
		sdk.NewAttribute(types.AttributeInitVolume, msg.InitialVolume.String()),
		sdk.NewAttribute(types.AttributeInitReserve, msg.InitialReserve.String()),
		sdk.NewAttribute(types.AttributeLimitVolume, msg.LimitVolume.String()),
		sdk.NewAttribute(types.AttributeCommissionCreateCoin, sdk.NewCoin(strings.ToLower(feeCoin), commission).String()),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

////////////////////////////////////////////////////////////////
// Updating coin handler
////////////////////////////////////////////////////////////////

func handleMsgUpdateCoin(ctx sdk.Context, k Keeper, msg types.MsgUpdateCoin) (*sdk.Result, error) {
	coin, err := k.GetCoin(ctx, strings.ToLower(msg.Symbol))
	if err != nil {
		return nil, types.ErrCoinAlreadyExist(msg.Symbol)
	}

	if !coin.Creator.Equals(msg.Sender) {
		return nil, types.ErrUpdateOnlyForCreator()
	}

	if coin.Volume.GT(msg.LimitVolume) {
		return nil, types.ErrLimitVolumeBroken(coin.Volume.String(), msg.LimitVolume.String())
	}

	coin.LimitVolume = msg.LimitVolume
	coin.Identity = msg.Identity

	k.SetCoin(ctx, coin)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

////////////////////////////////////////////////////////////////
// Transfer coins handlers
////////////////////////////////////////////////////////////////

func handleMsgSendCoin(ctx sdk.Context, k Keeper, msg types.MsgSendCoin) (*sdk.Result, error) {
	err := k.BankKeeper.SendCoins(ctx, msg.Sender, msg.Receiver, sdk.Coins{msg.Coin})
	if err != nil {
		return nil, types.ErrInternal(err.Error())
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		sdk.NewAttribute(types.AttributeCoin, msg.Coin.String()),
		sdk.NewAttribute(types.AttributeReceiver, msg.Receiver.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgMultiSendCoin(ctx sdk.Context, k Keeper, msg types.MsgMultiSendCoin) (*sdk.Result, error) {
	for i := range msg.Sends {
		err := k.BankKeeper.SendCoins(ctx, msg.Sender, msg.Sends[i].Receiver, sdk.Coins{msg.Sends[i].Coin})
		if err != nil {
			return nil, types.ErrInternal(err.Error())
		}

		ctx.EventManager().EmitEvent(sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeCoin, msg.Sends[i].Coin.String()),
			sdk.NewAttribute(types.AttributeReceiver, msg.Sends[i].Receiver.String()),
		))
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

////////////////////////////////////////////////////////////////
// Trading coins handlers
////////////////////////////////////////////////////////////////

func handleMsgBuyCoin(ctx sdk.Context, k Keeper, msg types.MsgBuyCoin) (*sdk.Result, error) {
	// Retrieve buyer account and it's balance of selling coins
	account := k.AccountKeeper.GetAccount(ctx, msg.Sender)
	balance := account.GetCoins().AmountOf(strings.ToLower(msg.MaxCoinToSell.Denom))

	// Retrieve the coin requested to buy
	coinToBuy, err := k.GetCoin(ctx, msg.CoinToBuy.Denom)
	if err != nil {
		return nil, types.ErrCoinDoesNotExist(msg.CoinToBuy.Denom)
	}
	if coinToBuy.Symbol != msg.CoinToBuy.Denom {
		return nil, types.ErrRetrievedAnotherCoin(msg.CoinToBuy.Denom, coinToBuy.Symbol)
	}

	// Retrieve the coin requested to sell
	coinToSell, err := k.GetCoin(ctx, msg.MaxCoinToSell.Denom)
	if err != nil {
		return nil, types.ErrCoinDoesNotExist(msg.MaxCoinToSell.Denom)
	}
	if coinToSell.Symbol != msg.MaxCoinToSell.Denom {
		return nil, types.ErrRetrievedAnotherCoin(msg.MaxCoinToSell.Denom, coinToSell.Symbol)
	}

	// Ensure supply limit of the coin to buy does not overflow
	if !coinToBuy.IsBase() {
		if coinToBuy.Volume.Add(msg.CoinToBuy.Amount).GT(coinToBuy.LimitVolume) {
			return nil, types.ErrTxBreaksVolumeLimit(coinToBuy.Volume.Add(msg.CoinToBuy.Amount).String(), coinToBuy.LimitVolume.String())
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
			return nil, types.ErrInsufficientCoinReserve()
		}

		amountToSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountToBuy)
		amountInBaseCoin = amountToBuy
	default:
		// Buyer buys custom coin for custom coin

		amountInBaseCoin = formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountToBuy)
		if amountInBaseCoin.GT(coinToSell.Reserve) {
			return nil, types.ErrInsufficientCoinReserve()
		}

		amountToSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountInBaseCoin)
	}

	// Ensure maximum amount of coins to sell (price guard)
	if amountToSell.GT(msg.MaxCoinToSell.Amount) {
		return nil, types.ErrMaximumValueToSellReached(msg.MaxCoinToSell.Amount.String(), amountToSell.String())
	}

	// Ensure reserve of the coin to sell does not underflow
	if !coinToSell.IsBase() {
		if coinToSell.Reserve.Sub(amountInBaseCoin).LT(types.MinCoinReserve(ctx)) {
			return nil, types.ErrTxBreaksMinReserveRule(MinCoinReserve(ctx).String(), amountInBaseCoin.String())
		}
	}

	// Ensure that buyer account holds enough coins to sell
	if balance.LT(amountToSell) {
		return nil, types.ErrInsufficientFunds(amountToSell.String(), balance.String())
	}

	// Update buyer account balances
	err = k.UpdateBalance(ctx, msg.MaxCoinToSell.Denom, amountToSell.Neg(), msg.Sender)
	if err != nil {
		return nil, types.ErrUpdateBalance(account.GetAddress().String(), err.Error())
	}
	err = k.UpdateBalance(ctx, msg.CoinToBuy.Denom, amountToBuy, msg.Sender)
	if err != nil {
		return nil, types.ErrUpdateBalance(account.GetAddress().String(), err.Error())
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
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		sdk.NewAttribute(types.AttributeCoinToBuy, sdk.NewCoin(msg.CoinToBuy.Denom, amountToBuy).String()),
		sdk.NewAttribute(types.AttributeCoinToSell, sdk.NewCoin(msg.MaxCoinToSell.Denom, amountToSell).String()),
		sdk.NewAttribute(types.AttributeAmountInBaseCoin, amountInBaseCoin.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSellCoin(ctx sdk.Context, k Keeper, msg types.MsgSellCoin, sellAll bool) (*sdk.Result, error) {
	// Retrieve seller account and it's balance of selling coins
	account := k.AccountKeeper.GetAccount(ctx, msg.Sender)
	balance := account.GetCoins().AmountOf(strings.ToLower(msg.CoinToSell.Denom))

	// Fill amount to sell in case of MsgSellAll
	if sellAll {
		msg.CoinToSell.Amount = balance
	}

	// Retrieve the coin requested to sell
	coinToSell, err := k.GetCoin(ctx, msg.CoinToSell.Denom)
	if err != nil {
		return nil, types.ErrCoinDoesNotExist(msg.CoinToSell.Denom)
	}
	if coinToSell.Symbol != msg.CoinToSell.Denom {
		return nil, types.ErrRetrievedAnotherCoin(msg.CoinToSell.Denom, coinToSell.Symbol)
	}

	// Retrieve the coin requested to buy
	coinToBuy, err := k.GetCoin(ctx, msg.MinCoinToBuy.Denom)
	if err != nil {
		return nil, types.ErrCoinDoesNotExist(msg.MinCoinToBuy.Denom)
	}
	if coinToBuy.Symbol != msg.MinCoinToBuy.Denom {
		return nil, types.ErrRetrievedAnotherCoin(msg.MinCoinToBuy.Denom, coinToBuy.Symbol)
	}

	// Ensure that seller account holds enough coins to sell
	if balance.LT(msg.CoinToSell.Amount) {
		return nil, types.ErrInsufficientFunds(msg.CoinToSell.String(), balance.String())
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
		return nil, types.ErrMinimumValueToBuyReached(amountToBuy.String(), msg.MinCoinToBuy.Amount.String())
	}

	// Ensure reserve of the coin to sell does not underflow
	if !coinToSell.IsBase() {
		if coinToSell.Reserve.Sub(amountInBaseCoin).LT(types.MinCoinReserve(ctx)) {
			return nil, types.ErrTxBreaksMinReserveRule(MinCoinReserve(ctx).String(), amountInBaseCoin.String())
		}
	}

	// Ensure supply limit of the coin to buy does not overflow
	if !coinToBuy.IsBase() {
		if coinToBuy.Volume.Add(amountToBuy).GT(coinToBuy.LimitVolume) {
			return nil, types.ErrTxBreaksVolumeLimit(coinToBuy.Volume.Add(amountToBuy).String(), coinToBuy.LimitVolume.String())
		}
	}

	// Update seller account balances
	err = k.UpdateBalance(ctx, msg.CoinToSell.Denom, amountToSell.Neg(), msg.Sender)
	if err != nil {
		return nil, types.ErrUpdateBalance(msg.Sender.String(), err.Error())
	}
	err = k.UpdateBalance(ctx, msg.MinCoinToBuy.Denom, amountToBuy, msg.Sender)
	if err != nil {
		return nil, types.ErrUpdateBalance(msg.Sender.String(), err.Error())
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
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		sdk.NewAttribute(types.AttributeCoinToSell, sdk.NewCoin(msg.CoinToSell.Denom, amountToSell).String()),
		sdk.NewAttribute(types.AttributeCoinToBuy, sdk.NewCoin(msg.MinCoinToBuy.Denom, amountToBuy).String()),
		sdk.NewAttribute(types.AttributeAmountInBaseCoin, amountInBaseCoin.String()),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

////////////////////////////////////////////////////////////////
// Redeem check handler
////////////////////////////////////////////////////////////////

func handleMsgRedeemCheck(ctx sdk.Context, k Keeper, msg types.MsgRedeemCheck) (*sdk.Result, error) {
	// Decode provided check from base58 format to raw bytes
	checkBytes := base58.Decode(msg.Check)
	if len(checkBytes) == 0 {
		return nil, types.ErrUnableDecodeCheck(msg.Check)
	}

	// Parse provided check from raw bytes to ensure it is valid
	check, err := types.ParseCheck(checkBytes)
	if err != nil {
		return nil, types.ErrInvalidCheck(err.Error())
	}

	// Decode provided proof from base64 format to raw bytes
	proof, err := base64.StdEncoding.DecodeString(msg.Proof)
	if err != nil {
		return nil, types.ErrUnableDecodeProof()
	}

	// Recover issuer address from check signature
	issuer, err := check.Sender()
	if err != nil {
		return nil, types.ErrUnableRecoverAddress(err.Error())
	}

	account := k.AccountKeeper.GetAccount(ctx, issuer)
	balance := account.GetCoins().AmountOf(strings.ToLower(check.Coin))

	// Retrieve the coin specified in the check
	coin, err := k.GetCoin(ctx, check.Coin)
	if err != nil {
		return nil, types.ErrCoinDoesNotExist(check.Coin)
	}
	if coin.Symbol != check.Coin {
		return nil, types.ErrRetrievedAnotherCoin(check.Coin, coin.Symbol)
	}

	feeCoin := k.GetBaseCoin(ctx)
	commission := helpers.UnitToPip(sdk.NewIntFromUint64(30))

	// Ensure that check issuer account holds enough coins
	amount := sdk.NewIntFromBigInt(check.Amount)
	if balance.LT(amount) {
		return nil, types.ErrInsufficientFunds(check.Amount.String()+coin.Symbol, balance.String()+coin.Symbol)
	}

	if feeCoin == strings.ToLower(check.Coin) {
		if balance.LT(amount.Add(commission)) {
			return nil, types.ErrInsufficientFunds(amount.String()+coin.Symbol, balance.String()+coin.Symbol)
		}
	} else {
		if ctx.BlockHeight() >= 32000 {
			feeBalance := account.GetCoins().AmountOf(feeCoin)
			if feeBalance.LT(commission) {
				return nil, types.ErrInsufficientFunds(commission.String()+feeCoin, feeBalance.String()+feeCoin)
			}
		}
	}

	// Ensure the proper chain ID is specified in the check
	if check.ChainID != ctx.ChainID() {
		return nil, types.ErrInvalidChainID(ctx.ChainID(), check.ChainID)
	}

	// Ensure nonce length
	if len(check.Nonce) > 16 {
		return nil, types.ErrInvalidNonce()
	}

	// Check block number
	if check.DueBlock < uint64(ctx.BlockHeight()) {
		return nil, types.ErrCheckExpired(
			strconv.FormatInt(int64(check.DueBlock), 10))
	}

	// Ensure check is not redeemed yet
	if k.IsCheckRedeemed(ctx, check) {
		return nil, types.ErrCheckRedeemed()
	}

	// Recover public key from check lock
	publicKeyA, err := check.LockPubKey()
	if err != nil {
		return nil, types.ErrUnableRecoverLockPkey(err.Error())
	}

	// Prepare bytes used to recover public key from provided proof
	senderAddressHash := make([]byte, 32)
	hw := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hw, []interface{}{
		msg.Sender,
	})
	if err != nil {
		return nil, types.ErrUnableRPLEncodeCheck(err.Error())
	}
	hw.Sum(senderAddressHash[:0])

	// Recover public key from provided proof
	publicKeyB, err := crypto.Ecrecover(senderAddressHash[:], proof)

	// Compare both public keys to ensure provided proof is correct
	if !bytes.Equal(publicKeyA, publicKeyB) {
		return nil, types.ErrInvalidProof(err.Error())
	}

	// Set check redeemed
	k.SetCheckRedeemed(ctx, check)

	// Update accounts balances
	err = k.UpdateBalance(ctx, feeCoin, commission.Neg(), issuer)
	if err != nil {
		return nil, types.ErrInsufficientFundsToPayCommission(commission.String())
	}
	err = k.BankKeeper.SendCoins(ctx, issuer, msg.Sender, sdk.Coins{sdk.NewCoin(coin.Symbol, amount)})
	if err != nil {
		return nil, types.ErrInternal(err.Error())
	}

	// Emit event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		sdk.NewAttribute(types.AttributeIssuer, issuer.String()),
		sdk.NewAttribute(types.AttributeCoin, sdk.NewCoin(check.Coin, sdk.NewIntFromBigInt(check.Amount)).String()),
		sdk.NewAttribute(types.AttributeNonce, new(big.Int).SetBytes(check.Nonce).String()),
		sdk.NewAttribute(types.AttributeDueBlock, strconv.FormatUint(check.DueBlock, 10)),
		sdk.NewAttribute(types.AttributeCommissionRedeemCheck, sdk.NewCoin(feeCoin, commission).String()),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
