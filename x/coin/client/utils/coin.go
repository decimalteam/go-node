package utils

import (
	"fmt"
	"strings"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"

	clientctx "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Check if coin exists
func ExistsCoin(cliCtx clientctx.CLIContext, symbol string) (bool, error) {
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, types.QueryGetCoin, symbol), nil)
	if err == nil {
		return res != nil, nil
	} else {
		return false, err
	}
}

// Return coin instance from State
func GetCoin(cliCtx clientctx.CLIContext, symbol string) (types.Coin, error) {
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, types.QueryGetCoin, symbol), nil)
	coin := types.Coin{}
	if err = cliCtx.Codec.UnmarshalJSON(res, &coin); err != nil {
		return coin, err
	}
	return coin, err
}

// Calculate amountToSell and amountToBuy for BuyCoin TX
// In CLI part amountToSell is maxAmountToSell
func BuyCoinCalculateAmounts(ctx sdk.Context, coinToBuy types.Coin, coinToSell types.Coin, wantsBuy sdk.Int, wantsSell sdk.Int) (amountBuy sdk.Int, amountSell sdk.Int, err *sdkerrors.Error) {
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
		return sdk.Int{}, sdk.Int{}, types.ErrTxBreaksVolumeLimit(coinToBuy.Volume.Add(wantsBuy).String(), coinToBuy.LimitVolume.String())
	}

	coinToSellMinReserve := formulas.GetReserveLimitFromCRR(coinToSell.CRR)
	if coinToSell.Reserve.Sub(amountSellInBaseCoin).LT(coinToSellMinReserve) && !coinToSell.IsBase() {
		return sdk.Int{}, sdk.Int{}, types.ErrTxBreaksMinReserveRule(types.MinCoinReserve(ctx).String(), coinToSell.Reserve.Sub(amountSellInBaseCoin).String())
	}
	return amountBuy, amountSell, nil
}

// Calculate amountToSell and amountToBuy for SellCoin TX
// In CLI part amountToBuy is minAmountToBuy
func SellCoinCalculateAmounts(ctx sdk.Context, coinToBuy types.Coin, coinToSell types.Coin, wantsBuy sdk.Int, wantsSell sdk.Int) (amountBuy sdk.Int, amountSell sdk.Int, err *sdkerrors.Error) {
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
		return sdk.Int{}, sdk.Int{}, types.ErrTxBreaksVolumeLimit(coinToBuy.Volume.Add(amountBuy).String(), coinToBuy.LimitVolume.String())
	}

	coinToSellMinReserve := formulas.GetReserveLimitFromCRR(coinToSell.CRR)
	if coinToSell.Reserve.Sub(amountSellInBase).LT(coinToSellMinReserve) && !coinToSell.IsBase() {
		return sdk.Int{}, sdk.Int{}, types.ErrTxBreaksMinReserveRule(types.MinCoinReserve(ctx).String(), coinToSell.Reserve.Sub(amountSellInBase).String())
	}

	// Limit minAmountToBuy in CLI
	if amountBuy.LT(wantsBuy) {
		return sdk.Int{}, sdk.Int{}, types.ErrMinimumValueToBuyReached(amountBuy.String(), wantsBuy.String())
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
