package utils

import (
	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Check if coin exists
func ExistsCoin(cliCtx client.CLIContext, symbol string) (bool, error) {
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, types.QueryGetCoin, symbol), nil)
	if err == nil {
		return res != nil, nil
	} else {
		return false, err
	}
}

// Return coin instance from State
func GetCoin(cliCtx client.CLIContext, symbol string) (types.Coin, error) {
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.ModuleName, types.QueryGetCoin, symbol), nil)
	coin := types.Coin{}
	if err = cliCtx.Codec.UnmarshalJSON(res, &coin); err != nil {
		return coin, err
	}
	return coin, err
}

// Calculate amountToSell and amountToBuy for BuyCoin TX
// In CLI part amountToSell is maxAmountToSell
func BuyCoinCalculateAmounts(coinToBuy types.Coin, coinToSell types.Coin, amountToBuy sdk.Int, amountToSell sdk.Int) (amountBuy sdk.Int, amountSell sdk.Int, err sdk.Error) {
	amountBuyInBaseCoin := formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountToBuy)
	amountSellRequired := formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountBuyInBaseCoin)

	var amountSellInBaseCoin sdk.Int

	if amountSellRequired.GT(amountToSell) {
		amountSell = amountToSell
		amountSellInBaseCoin = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountSell)
	} else {
		amountSell = amountSellRequired
		amountSellInBaseCoin = amountBuyInBaseCoin
	}
	amountBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountSellInBaseCoin)

	if coinToBuy.Volume.Add(amountToBuy).GT(coinToBuy.LimitVolume) {
		return sdk.Int{}, sdk.Int{}, sdk.NewError(types.DefaultCodespace, types.TxBreaksVolumeLimit, "Tx breaks LimitVolume rule")
	}
	coinToSellMinReserve := formulas.GetReserveLimitFromCRR(coinToSell.CRR)
	if coinToSell.Reserve.Sub(amountSellInBaseCoin).LT(coinToSellMinReserve) {
		return sdk.Int{}, sdk.Int{}, sdk.NewError(types.DefaultCodespace, types.TxBreaksMinReserveLimit, "Tx breaks MinReserveLimit rule")
	}
	return amountBuy, amountSell, nil
}

// Calculate amountToSell and amountToBuy for SellCoin TX
// In CLI part amountToBuy is minAmountToBuy
func SellCoinCalculateAmounts(coinToBuy types.Coin, coinToSell types.Coin, amountToBuy sdk.Int, amountToSell sdk.Int) (amountBuy sdk.Int, amountSell sdk.Int, err sdk.Error) {
	amountSellInBaseCoin := formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountToSell)

	amountBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountSellInBaseCoin)

	if coinToBuy.Volume.Add(amountBuy).GT(coinToBuy.LimitVolume) {
		return sdk.Int{}, sdk.Int{}, sdk.NewError(types.DefaultCodespace, types.TxBreaksVolumeLimit, "Tx breaks LimitVolume rule")
	}

	coinToSellMinReserve := formulas.GetReserveLimitFromCRR(coinToSell.CRR)
	if coinToSell.Reserve.Sub(amountSellInBaseCoin).LT(coinToSellMinReserve) {
		return sdk.Int{}, sdk.Int{}, sdk.NewError(types.DefaultCodespace, types.TxBreaksMinReserveLimit, "Tx breaks MinReserveLimit rule")
	}

	// Limit minAmountToBuy in CLI
	if amountBuy.LT(amountToBuy) {
		return sdk.Int{}, sdk.Int{}, sdk.NewError(types.DefaultCodespace, types.AmountBuyIsTooSmall, "Amount you will receive less than minimum")
	}

	return amountBuy, amountToSell, nil
}
