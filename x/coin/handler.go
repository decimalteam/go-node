package coin

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/utils/helpers"
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

const (
	SendFee        int64 = 10
	SellFee        int64 = 100
	BuyFee         int64 = 100
	RedeemCheckFee int64 = 30
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
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCreateCoin:
			return handleMsgCreateCoin(ctx, k, msg)
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
	var coin = types.Coin{
		Title:       msg.Title,
		CRR:         msg.ConstantReserveRatio,
		Symbol:      msg.Symbol,
		Reserve:     msg.InitialReserve,
		LimitVolume: msg.LimitVolume,
		Volume:      msg.InitialVolume,
	}

	commission, feeCoin, err := k.GetCommission(ctx, helpers.BipToPip(getCreateCoinCommission(msg.Symbol)))
	if err != nil {
		return nil, types.ErrCalculateCommission(err)
	}

	acc := k.AccountKeeper.GetAccount(ctx, msg.Sender)
	balance := acc.GetCoins()
	if balance.AmountOf(strings.ToLower(cliUtils.GetBaseCoin())).LT(msg.InitialReserve) {
		return nil, sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinReserve, "Not enough coin to reserve")
	}

	if balance.AmountOf(strings.ToLower(cliUtils.GetBaseCoin())).LT(commission) {
		return nil, types.ErrorInsufficientCoinToPayCommission(commission.String())
	}

	if balance.AmountOf(strings.ToLower(cliUtils.GetBaseCoin())).LT(commission.Add(msg.InitialReserve)) {
		return nil, types.ErrorInsufficientFunds(commission.Add(msg.InitialReserve).String())
	}

	err = k.UpdateBalance(ctx, strings.ToLower(cliUtils.GetBaseCoin()), msg.InitialReserve.Neg(), msg.Sender)
	if err != nil {
		return nil, types.ErrorUpdateBalance(err)
	}

	k.SetCoin(ctx, coin)
	err = k.UpdateBalance(ctx, strings.ToLower(coin.Symbol), msg.InitialVolume, msg.Sender)
	if err != nil {
		return nil, types.ErrorUpdateBalance(err)
	}

	err = k.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), msg.Sender)
	if err != nil {
		return nil, types.ErrorUpdateBalance(err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeSymbol, msg.Symbol),
			sdk.NewAttribute(types.AttributeTitle, msg.Title),
			sdk.NewAttribute(types.AttributeCRR, strconv.FormatUint(uint64(msg.ConstantReserveRatio), 10)),
			sdk.NewAttribute(types.AttributeInitVolume, msg.InitialVolume.String()),
			sdk.NewAttribute(types.AttributeInitReserve, msg.InitialReserve.String()),
			sdk.NewAttribute(types.AttributeLimitVolume, msg.LimitVolume.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

////////////////////////////////////////////////////////////////
// Transfer coins handlers
////////////////////////////////////////////////////////////////

func handleMsgSendCoin(ctx sdk.Context, k Keeper, msg types.MsgSendCoin) (*sdk.Result, error) {
	commission, feeCoin, err := k.GetCommission(ctx, helpers.UnitToPip(SendFee))
	if err != nil {
		return nil, types.ErrCalculateCommission(err)
	}

	acc := k.AccountKeeper.GetAccount(ctx, msg.Sender)
	balance := acc.GetCoins()

	if balance.AmountOf(strings.ToLower(cliUtils.GetBaseCoin())).LT(commission) {
		return nil, types.ErrorInsufficientCoinToPayCommission(commission.String())
	}

	err = k.BankKeeper.SendCoins(ctx, msg.Sender, msg.Receiver, sdk.Coins{msg.Coin})
	if err != nil {
		return nil, sdkerrors.New(types.DefaultCodespace, 6, err.Error())
	}

	err = k.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), msg.Sender)
	if err != nil {
		return nil, types.ErrorUpdateBalance(err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeCoin, msg.Coin.Denom),
			sdk.NewAttribute(types.AttributeAmount, msg.Coin.Amount.String()),
			sdk.NewAttribute(types.AttributeReceiver, msg.Receiver.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgMultiSendCoin(ctx sdk.Context, k Keeper, msg types.MsgMultiSendCoin) (*sdk.Result, error) {
	commission, feeCoin, err := k.GetCommission(ctx, helpers.UnitToPip(SendFee+int64((len(msg.Sends)-1)*5)))
	if err != nil {
		return nil, types.ErrCalculateCommission(err)
	}

	acc := k.AccountKeeper.GetAccount(ctx, msg.Sender)
	balance := acc.GetCoins()

	if balance.AmountOf(strings.ToLower(cliUtils.GetBaseCoin())).LT(commission) {
		return nil, types.ErrorInsufficientCoinToPayCommission(commission.String())
	}

	for i := range msg.Sends {
		_ = k.BankKeeper.SendCoins(ctx, msg.Sender, msg.Sends[i].Receiver, sdk.Coins{msg.Sends[i].Coin})
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
				sdk.NewAttribute(types.AttributeCoin, msg.Sends[i].Coin.Denom),
				sdk.NewAttribute(types.AttributeAmount, msg.Sends[i].Coin.Amount.String()),
				sdk.NewAttribute(types.AttributeReceiver, msg.Sends[i].Receiver.String()),
			),
		})
	}

	err = k.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), msg.Sender)
	if err != nil {
		return nil, types.ErrorUpdateBalance(err)
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

////////////////////////////////////////////////////////////////
// Trading coins handlers
////////////////////////////////////////////////////////////////

func handleMsgBuyCoin(ctx sdk.Context, k Keeper, msg types.MsgBuyCoin) (*sdk.Result, error) {
	commission, feeCoin, err := k.GetCommission(ctx, helpers.UnitToPip(BuyFee))
	if err != nil {
		return nil, types.ErrCalculateCommission(err)
	}

	// Retrieve buyer account and it's balance of selling coins
	account := k.AccountKeeper.GetAccount(ctx, msg.Sender)
	balance := account.GetCoins().AmountOf(strings.ToLower(msg.MaxCoinToSell.Denom))

	if balance.LT(commission) {
		return nil, types.ErrorInsufficientCoinToPayCommission(commission.String())
	}

	// Retrieve the coin requested to buy
	coinToBuy, err := k.GetCoin(ctx, msg.CoinToBuy.Denom)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to buy: %v", msg.CoinToBuy, err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}
	if coinToBuy.Symbol != msg.CoinToBuy.Denom {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to buy: retrieved coin %s instead", msg.CoinToBuy, coinToBuy.Symbol)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}

	// Retrieve the coin requested to sell
	coinToSell, err := k.GetCoin(ctx, msg.MaxCoinToSell.Denom)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to sell: %v", msg.MaxCoinToSell, err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}
	if coinToSell.Symbol != msg.MaxCoinToSell.Denom {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to sell: retrieved coin %s instead", msg.MaxCoinToSell, coinToSell.Symbol)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToSellNotExists, errMsg)
	}

	// Ensure supply limit of the coin to buy does not overflow
	if !coinToBuy.IsBase() {
		if coinToBuy.Volume.Add(msg.CoinToBuy.Amount).GT(coinToBuy.LimitVolume) {
			errMsg := fmt.Sprintf(
				"Wanted to buy %f %s, but this operation will overflow coin supply limit",
				floatFromInt(msg.CoinToBuy.Amount), coinToBuy.Symbol,
			)
			return nil, sdkerrors.New(types.DefaultCodespace, types.TxBreaksVolumeLimit, errMsg)
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
		amountToSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountToBuy)
		amountInBaseCoin = amountToBuy
	default:
		// Buyer buys custom coin for custom coin
		amountInBaseCoin = formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountToBuy)
		amountToSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountInBaseCoin)
	}

	// Ensure maximum amount of coins to sell (price guard)
	if amountToSell.GT(msg.MaxCoinToSell.Amount) {
		errMsg := fmt.Sprintf(
			"Wanted to sell maximum %f %s, but required to spend %f %s at the moment",
			floatFromInt(msg.MaxCoinToSell.Amount), coinToSell.Symbol, floatFromInt(amountToSell), coinToSell.Symbol,
		)
		return nil, sdkerrors.New(types.DefaultCodespace, types.MaximumValueToSellReached, errMsg)
	}

	// Ensure reserve of the coin to sell does not underflow
	if !coinToSell.IsBase() {
		// TODO: Compare with some minimal reserve value?
		if coinToSell.Reserve.Sub(amountInBaseCoin).Sign() < 0 {
			errMsg := fmt.Sprintf(
				"Wanted to sell %f %s, but this operation will underflow coin reserve",
				floatFromInt(amountInBaseCoin), coinToSell.Symbol,
			)
			return nil, sdkerrors.New(types.DefaultCodespace, types.TxBreaksMinReserveLimit, errMsg)
		}
	}

	// Ensure that buyer account holds enough coins to sell
	if balance.LT(amountToSell) {
		errMsg := fmt.Sprintf(
			"Wanted to sell %f %s, but available only %f %s at the moment",
			floatFromInt(balance), coinToSell.Symbol, floatFromInt(amountToSell), coinToSell.Symbol,
		)
		return nil, sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, errMsg)
	}

	if balance.LT(amountToSell.Add(commission)) {
		return nil, types.ErrorInsufficientFunds(amountToSell.Add(commission).String())
	}

	// Update buyer account balances
	err = k.UpdateBalance(ctx, msg.MaxCoinToSell.Denom, amountToSell.Neg(), msg.Sender)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to update balance of account %s: %v", account.GetAddress(), err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}
	err = k.UpdateBalance(ctx, msg.CoinToBuy.Denom, amountToBuy, msg.Sender)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to update balance of account %s: %v", account.GetAddress(), err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}

	err = k.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), msg.Sender)
	if err != nil {
		return nil, types.ErrorUpdateBalance(err)
	}

	// Update coins
	if !coinToSell.IsBase() {
		k.UpdateCoin(ctx, coinToSell, coinToSell.Reserve.Sub(amountInBaseCoin), coinToSell.Volume.Sub(amountToSell))
	}
	if !coinToBuy.IsBase() {
		k.UpdateCoin(ctx, coinToBuy, coinToBuy.Reserve.Add(amountInBaseCoin), coinToBuy.Volume.Add(amountToBuy))
	}

	// Emit transaction events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeCoinToBuy, msg.CoinToBuy.Denom),
			sdk.NewAttribute(types.AttributeAmountToBuy, amountToBuy.String()),
			sdk.NewAttribute(types.AttributeCoinToSell, msg.MaxCoinToSell.Denom),
			sdk.NewAttribute(types.AttributeAmountToSell, amountToSell.String()),
			sdk.NewAttribute(types.AttributeAmountInBaseCoin, amountInBaseCoin.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgSellCoin(ctx sdk.Context, k Keeper, msg types.MsgSellCoin, sellAll bool) (*sdk.Result, error) {
	commission, feeCoin, err := k.GetCommission(ctx, helpers.UnitToPip(SellFee))
	if err != nil {
		return nil, types.ErrCalculateCommission(err)
	}

	// Retrieve seller account and it's balance of selling coins
	account := k.AccountKeeper.GetAccount(ctx, msg.Sender)
	balance := account.GetCoins().AmountOf(strings.ToLower(msg.CoinToSell.Denom))

	if balance.LT(commission) {
		return nil, types.ErrorInsufficientCoinToPayCommission(commission.String())
	}

	// Fill amount to sell in case of MsgSellAll
	if sellAll {
		msg.CoinToSell.Amount = balance
	}

	// Retrieve the coin requested to sell
	coinToSell, err := k.GetCoin(ctx, msg.CoinToSell.Denom)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to sell: %v", msg.CoinToSell, err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToSellNotExists, errMsg)
	}
	if coinToSell.Symbol != msg.CoinToSell.Denom {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to sell: retrieved coin %s instead", msg.CoinToSell, coinToSell.Symbol)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToSellNotExists, errMsg)
	}

	// Retrieve the coin requested to buy
	coinToBuy, err := k.GetCoin(ctx, msg.MinCoinToBuy.Denom)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to buy: %v", msg.MinCoinToBuy, err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}
	if coinToBuy.Symbol != msg.MinCoinToBuy.Denom {
		errMsg := fmt.Sprintf("Unable to retrieve coin %s requested to buy: retrieved coin %s instead", msg.MinCoinToBuy, coinToBuy.Symbol)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CoinToBuyNotExists, errMsg)
	}

	// Ensure that seller account holds enough coins to sell
	if balance.LT(msg.CoinToSell.Amount) {
		errMsg := fmt.Sprintf(
			"Wanted to sell %f %s, but available only %f %s at the moment",
			floatFromInt(balance), coinToSell.Symbol, floatFromInt(msg.CoinToSell.Amount), coinToSell.Symbol,
		)
		return nil, sdkerrors.New(types.DefaultCodespace, types.InsufficientCoinToSell, errMsg)
	}

	if balance.LT(msg.CoinToSell.Amount.Add(commission)) {
		return nil, types.ErrorInsufficientFunds(msg.CoinToSell.Amount.Add(commission).String())
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
		errMsg := fmt.Sprintf(
			"Wanted to buy minimum %f %s, but expected to receive %f %s at the moment",
			floatFromInt(msg.MinCoinToBuy.Amount), coinToBuy.Symbol, floatFromInt(amountToBuy), coinToBuy.Symbol,
		)
		return nil, sdkerrors.New(types.DefaultCodespace, types.MinimumValueToBuyReached, errMsg)
	}

	// Ensure reserve of the coin to sell does not underflow
	if !coinToSell.IsBase() {
		// TODO: Compare with some minimal reserve value?
		if coinToSell.Reserve.Sub(amountInBaseCoin).Sign() < 0 {
			errMsg := fmt.Sprintf(
				"Wanted to sell %f %s, but this operation will underflow coin reserve",
				floatFromInt(amountInBaseCoin), coinToSell.Symbol,
			)
			return nil, sdkerrors.New(types.DefaultCodespace, types.TxBreaksMinReserveLimit, errMsg)
		}
	}

	// Ensure supply limit of the coin to buy does not overflow
	if !coinToBuy.IsBase() {
		if coinToBuy.Volume.Add(amountToBuy).GT(coinToBuy.LimitVolume) {
			errMsg := fmt.Sprintf(
				"Wanted to buy %f %s, but this operation will overflow coin supply limit",
				floatFromInt(amountToBuy), coinToBuy.Symbol,
			)
			return nil, sdkerrors.New(types.DefaultCodespace, types.TxBreaksVolumeLimit, errMsg)
		}
	}

	// Update seller account balances
	err = k.UpdateBalance(ctx, msg.CoinToSell.Denom, amountToSell.Neg(), msg.Sender)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to update balance of account %s: %v", account.GetAddress(), err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}
	err = k.UpdateBalance(ctx, msg.MinCoinToBuy.Denom, amountToBuy, msg.Sender)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to update balance of account %s: %v", account.GetAddress(), err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}

	err = k.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), msg.Sender)
	if err != nil {
		return nil, types.ErrorUpdateBalance(err)
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
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(types.AttributeCoinToSell, msg.CoinToSell.Denom),
			sdk.NewAttribute(types.AttributeAmountToSell, amountToSell.String()),
			sdk.NewAttribute(types.AttributeCoinToBuy, msg.MinCoinToBuy.Denom),
			sdk.NewAttribute(types.AttributeAmountToBuy, amountToBuy.String()),
			sdk.NewAttribute(types.AttributeAmountInBaseCoin, amountInBaseCoin.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

////////////////////////////////////////////////////////////////
// Redeem check handler
////////////////////////////////////////////////////////////////

func handleMsgRedeemCheck(ctx sdk.Context, k Keeper, msg types.MsgRedeemCheck) (*sdk.Result, error) {
	commission, feeCoin, err := k.GetCommission(ctx, helpers.UnitToPip(RedeemCheckFee))
	if err != nil {
		return nil, types.ErrCalculateCommission(err)
	}

	// Retrieve seller account and it's balance of selling coins
	account := k.AccountKeeper.GetAccount(ctx, msg.Sender)
	balanceFee := account.GetCoins().AmountOf(strings.ToLower(feeCoin))

	// Decode provided check from base64 format to raw bytes
	checkBytes, err := base64.StdEncoding.DecodeString(msg.Check)
	if err != nil {
		msgError := "unable to decode check from base64"
		return nil, sdkerrors.New(types.DefaultCodespace, types.InvalidCheck, msgError)
	}

	// Parse provided check from raw bytes to ensure it is valid
	check, err := types.ParseCheck(checkBytes)
	if err != nil {
		msgError := fmt.Sprintf("unable to parse check: %s", err.Error())
		return nil, sdkerrors.New(types.DefaultCodespace, types.InvalidCheck, msgError)
	}

	// Decode provided proof from base64 format to raw bytes
	proof, err := base64.StdEncoding.DecodeString(msg.Proof)
	if err != nil {
		msgError := "unable to decode proof from base64"
		return nil, sdkerrors.New(types.DefaultCodespace, types.InvalidProof, msgError)
	}

	// Recover issuer address from check signature
	issuer, err := check.Sender()
	if err != nil {
		errMsg := fmt.Sprintf("unable to recover check issuer address: %s", err.Error)
		return nil, sdkerrors.New(types.DefaultCodespace, types.InvalidCheck, errMsg)
	}

	balance := account.GetCoins().AmountOf(strings.ToLower(check.Coin))

	// Retrieve the coin specified in the check
	coin, err := k.GetCoin(ctx, check.Coin)
	if err != nil {
		errMsg := fmt.Sprintf("unable to retrieve coin %s requested to sell: %v", check.Coin, err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.InvalidCoinSymbol, errMsg)
	}
	if coin.Symbol != check.Coin {
		errMsg := fmt.Sprintf("unable to retrieve coin %s requested to sell: retrieved coin %s instead", check.Coin, coin.Symbol)
		return nil, sdkerrors.New(types.DefaultCodespace, types.InvalidCoinSymbol, errMsg)
	}

	// Ensure that check issuer account holds enough coins
	amount := sdk.NewIntFromBigInt(check.Amount)
	if balance.LT(amount) {
		errMsg := fmt.Sprintf(
			"wanted to send %f %s, but available only %f %s at the moment",
			floatFromInt(balance), coin.Symbol, floatFromInt(amount), coin.Symbol,
		)
		return nil, sdkerrors.New(types.DefaultCodespace, types.InvalidAmount, errMsg)
	}

	if balanceFee.LT(commission) {
		return nil, types.ErrorInsufficientCoinToPayCommission(commission.String())
	}

	if feeCoin == check.Coin && balance.LT(amount.Add(commission)) {
		return nil, types.ErrorInsufficientFunds(amount.Add(commission).String())
	}

	// Ensure the proper chain ID is specified in the check
	if check.ChainID != ctx.ChainID() {
		errMsg := fmt.Sprintf("wanted chain ID %s, but check is issued for chain with ID %s", ctx.ChainID(), check.ChainID)
		return nil, sdkerrors.New(types.DefaultCodespace, types.InvalidChainID, errMsg)
	}

	// Ensure nonce length
	if len(check.Nonce) > 16 {
		errMsg := "nonce is too big (should be up to 16 bytes)"
		return nil, sdkerrors.New(types.DefaultCodespace, types.InvalidNonce, errMsg)
	}

	// Check block number
	if check.DueBlock < uint64(ctx.BlockHeight()) {
		errMsg := fmt.Sprintf("check was expired at block %d", check.DueBlock)
		return nil, sdkerrors.New(types.DefaultCodespace, types.CheckExpired, errMsg)
	}

	// Ensure check is not redeemed yet
	if k.IsCheckRedeemed(ctx, check) {
		errMsg := "check was redeemed already"
		return nil, sdkerrors.New(types.DefaultCodespace, types.CheckRedeemed, errMsg)
	}

	// Recover public key from check lock
	publicKeyA, err := check.LockPubKey()
	if err != nil {
		msgError := fmt.Sprintf("unable to recover lock public key from check: %s", err.Error())
		return nil, sdkerrors.New(types.DefaultCodespace, types.InvalidCheck, msgError)
	}

	// Prepare bytes used to recover public key from provided proof
	senderAddressHash := make([]byte, 32)
	hw := sha3.NewLegacyKeccak256()
	err = rlp.Encode(hw, []interface{}{
		msg.Sender,
	})
	if err != nil {
		msgError := fmt.Sprintf("unable to RLP encode check sender address: %s", err.Error())
		return nil, sdkerrors.New(types.DefaultCodespace, types.InvalidCheck, msgError)
	}
	hw.Sum(senderAddressHash[:0])

	// Recover public key from provided proof
	publicKeyB, err := crypto.Ecrecover(senderAddressHash[:], proof)

	// Compare both public keys to ensure provided proof is correct
	if !bytes.Equal(publicKeyA, publicKeyB) {
		msgError := fmt.Sprintf("provided proof is invalid", err.Error())
		return nil, sdkerrors.New(types.DefaultCodespace, types.InvalidProof, msgError)
	}

	// Set check redeemed
	k.SetCheckRedeemed(ctx, check)

	// Update accounts balances
	err = k.UpdateBalance(ctx, coin.Symbol, amount.Neg(), issuer)
	if err != nil {
		errMsg := fmt.Sprintf("unable to update balance of check issuer account %s: %v", issuer, err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}
	err = k.UpdateBalance(ctx, coin.Symbol, amount, msg.Sender)
	if err != nil {
		errMsg := fmt.Sprintf("unable to update balance of check redeemer account %s: %v", msg.Sender, err)
		return nil, sdkerrors.New(types.DefaultCodespace, types.UpdateBalanceError, errMsg)
	}

	err = k.UpdateBalance(ctx, strings.ToLower(feeCoin), commission.Neg(), msg.Sender)
	if err != nil {
		return nil, types.ErrorUpdateBalance(err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		sdk.NewAttribute(types.AttributeIssuer, issuer.String()),
		sdk.NewAttribute(types.AttributeCoin, check.Coin),
		sdk.NewAttribute(types.AttributeAmount, check.Amount.String()),
		sdk.NewAttribute(types.AttributeCheckNonce, new(big.Int).SetBytes(check.Nonce).String()),
		sdk.NewAttribute(types.AttributeDueBlock, strconv.FormatUint(check.DueBlock, 10)),
	))

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
