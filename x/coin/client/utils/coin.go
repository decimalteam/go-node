package utils

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"fmt"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"strings"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	clientctx "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Check if coin exists
func ExistsCoin(clientCtx clientctx.Context, symbol string) (bool, error) {
	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types2.ModuleName, types2.QueryGetCoin, symbol), nil)
	if err == nil {
		return res != nil, nil
	} else {
		return false, err
	}
}

// Return coin instance from State
func GetCoin(clientCtx clientctx.Context, symbol string) (types2.Coin, error) {
	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types2.ModuleName, types2.QueryGetCoin, symbol), nil)
	coin := types2.Coin{}
	if err = clientCtx.LegacyAmino.UnmarshalJSON(res, &coin); err != nil {
		return coin, err
	}
	return coin, err
}

func GetAccountCoins(clientCtx clientctx.Context, addr sdk.AccAddress) (sdk.Coins, error) {
	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("%s/%s/%s", bankTypes.ModuleName, "balances", addr), nil)
	coins := sdk.Coins{}

	if err != nil {
		return coins, err
	}

	if err = clientCtx.LegacyAmino.UnmarshalJSON(res, &coins); err != nil {
		return coins, err
	}

	return coins, nil
}

// Calculate amountToSell and amountToBuy for BuyCoin TX
// In CLI part amountToSell is maxAmountToSell
func BuyCoinCalculateAmounts(ctx sdk.Context, coinToBuy types2.Coin, coinToSell types2.Coin, wantsBuy sdk.Int, wantsSell sdk.Int) (amountBuy sdk.Int, amountSell sdk.Int, err *sdkerrors.Error) {
	var amountSellInBaseCoin sdk.Int
	var amountBuyInBaseCoin sdk.Int
	if coinToSell.IsBase() {
		fmt.Printf("####### [BuyCoinCalculateAmounts] Coin to sell is base\n")
		amountBuyInBaseCoin = formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, wantsBuy)
		if amountBuyInBaseCoin.LTE(wantsSell) {
			fmt.Printf("####### [BuyCoinCalculateAmounts] Coin to buy in base coin is less than limit: %s <= %s\n", amountBuyInBaseCoin, wantsSell)
			amountBuy = wantsBuy
			amountSell = amountBuyInBaseCoin
			amountSellInBaseCoin = amountSell
		} else {
			fmt.Printf("####### [BuyCoinCalculateAmounts] Amount to buy in base coin is greater than limit: %s > %s\n", amountBuyInBaseCoin, wantsSell)
			amountBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, wantsSell)
			amountSell = wantsSell
			amountSellInBaseCoin = amountSell
		}
	} else if coinToBuy.IsBase() {
		fmt.Printf("####### [BuyCoinCalculateAmounts] Coin to buy is base\n")
		amountSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, wantsBuy)
		if amountSell.LTE(wantsSell) {
			fmt.Printf("####### [BuyCoinCalculateAmounts] Amount to sell is less than limit: %s <= %s\n", amountSell, wantsSell)
			amountBuy = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountSell)
			amountBuyInBaseCoin = amountBuy
			amountSellInBaseCoin = amountBuy
		} else {
			fmt.Printf("####### [BuyCoinCalculateAmounts] Amount to sell is greater than limit: %s > %s\n", amountSell, wantsSell)
			amountBuy = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, wantsSell)
			amountBuyInBaseCoin = amountBuy
			amountSellInBaseCoin = amountBuy
		}
	} else {
		fmt.Printf("####### [BuyCoinCalculateAmounts] Coins to buy and sell are both custom\n")
		amountBuyInBaseCoin := formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, wantsBuy)
		amountSellRequired := formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountBuyInBaseCoin)

		if amountSellRequired.GT(wantsSell) {
			fmt.Printf("####### [BuyCoinCalculateAmounts] Amount to sell required is greater than limit: %s <= %s\n", amountSellRequired, wantsSell)
			amountSell = wantsSell
			amountSellInBaseCoin = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountSell)
		} else {
			fmt.Printf("####### [BuyCoinCalculateAmounts] Amount to sell required is less than limit: %s <= %s\n", amountSellRequired, wantsSell)
			amountSell = amountSellRequired
			amountSellInBaseCoin = amountBuyInBaseCoin
		}
		amountBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountSellInBaseCoin)
	}

	if coinToBuy.Volume.Add(wantsBuy).GT(coinToBuy.LimitVolume) && !coinToBuy.IsBase() {
		return sdk.Int{}, sdk.Int{}, types2.ErrTxBreaksVolumeLimit(coinToBuy.Volume.Add(wantsBuy).String(), coinToBuy.LimitVolume.String())
	}

	coinToSellMinReserve := formulas.GetReserveLimitFromCRR(coinToSell.CRR)
	if coinToSell.Reserve.Sub(amountSellInBaseCoin).LT(coinToSellMinReserve) && !coinToSell.IsBase() {
		return sdk.Int{}, sdk.Int{}, types2.ErrTxBreaksMinReserveRule(ctx, coinToSell.Reserve.Sub(amountSellInBaseCoin).String())
	}
	return amountBuy, amountSell, nil
}

// Calculate amountToSell and amountToBuy for SellCoin TX
// In CLI part amountToBuy is minAmountToBuy
func SellCoinCalculateAmounts(ctx sdk.Context, coinToBuy types2.Coin, coinToSell types2.Coin, wantsBuy sdk.Int, wantsSell sdk.Int) (amountBuy sdk.Int, amountSell sdk.Int, err *sdkerrors.Error) {
	var amountSellInBase sdk.Int

	if coinToSell.IsBase() {
		amountSellInBase = wantsSell
		amountBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, wantsSell)
	} else if coinToBuy.IsBase() {
		amountSellInBase = formulas.CalculateSaleReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, wantsSell)
		amountBuy = amountSellInBase
		return formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, wantsSell), wantsSell, nil
	} else {
		amountSellInBaseCoin := formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, wantsSell)
		amountBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountSellInBaseCoin)
	}

	if coinToBuy.Volume.Add(amountBuy).GT(coinToBuy.LimitVolume) && !coinToBuy.IsBase() {
		return sdk.Int{}, sdk.Int{}, types2.ErrTxBreaksVolumeLimit(coinToBuy.Volume.Add(amountBuy).String(), coinToBuy.LimitVolume.String())
	}

	coinToSellMinReserve := formulas.GetReserveLimitFromCRR(coinToSell.CRR)
	if coinToSell.Reserve.Sub(amountSellInBase).LT(coinToSellMinReserve) && !coinToSell.IsBase() {
		return sdk.Int{}, sdk.Int{}, types2.ErrTxBreaksMinReserveRule(ctx, coinToSell.Reserve.Sub(amountSellInBase).String())
	}

	// Limit minAmountToBuy in CLI
	if amountBuy.LT(wantsBuy) {
		return sdk.Int{}, sdk.Int{}, types2.ErrMinimumValueToBuyReached(amountBuy.String(), wantsBuy.String())
	}

	return amountBuy, wantsSell, nil
}

func GetBaseCoin() string {
	if strings.HasPrefix(config.ChainID, "decimal-testnet") {
		return config.SymbolTestBaseCoin
	} else {
		return config.SymbolBaseCoin
	}
}
