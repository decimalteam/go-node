package utils

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
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
func BuyCoinCalculateAmounts(coinToBuy types.Coin, coinToSell types.Coin, wantsBuy sdk.Int, wantsSell sdk.Int) (amountBuy sdk.Int, amountSell sdk.Int, err sdkerrors.Error) {
	var amountSellInBaseCoin sdk.Int
	var amountBuyInBaseCoin sdk.Int
	if coinToSell.IsBase() {
		amountBuyInBaseCoin = formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, wantsBuy)
		if amountBuyInBaseCoin.LT(wantsSell) {
			amountBuy = wantsBuy
			amountSell = amountBuyInBaseCoin
			amountSellInBaseCoin = amountSell
		} else {
			amountBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, wantsSell)
			amountSell = wantsSell
			amountSellInBaseCoin = amountSell
		}
	} else if coinToBuy.IsBase() {
		amountSell = formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, wantsBuy)
		if amountSell.LTE(wantsSell) {
			amountBuy = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountSell)
			amountBuyInBaseCoin = amountBuy
			amountSellInBaseCoin = amountBuy
		} else {
			amountBuy = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, wantsSell)
			amountBuyInBaseCoin = amountBuy
			amountSellInBaseCoin = amountBuy
		}
	} else {
		amountBuyInBaseCoin := formulas.CalculatePurchaseAmount(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, wantsBuy)
		amountSellRequired := formulas.CalculateSaleAmount(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountBuyInBaseCoin)

		if amountSellRequired.GT(wantsSell) {
			amountSell = wantsSell
			amountSellInBaseCoin = formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountSell)
		} else {
			amountSell = amountSellRequired
			amountSellInBaseCoin = amountBuyInBaseCoin
		}
		amountBuy = formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountSellInBaseCoin)
	}

	if coinToBuy.Volume.Add(wantsBuy).GT(coinToBuy.LimitVolume) && !coinToBuy.IsBase() {
		return sdk.Int{}, sdk.Int{}, sdkerrors.New(types.DefaultCodespace, types.TxBreaksVolumeLimit, "Tx breaks LimitVolume rule")
	}

	coinToSellMinReserve := formulas.GetReserveLimitFromCRR(coinToSell.CRR)
	if coinToSell.Reserve.Sub(amountSellInBaseCoin).LT(coinToSellMinReserve) && !coinToSell.IsBase() {
		return sdk.Int{}, sdk.Int{}, sdkerrors.New(types.DefaultCodespace, types.TxBreaksMinReserveLimit, "Tx breaks MinReserveLimit rule")
	}
	return amountBuy, amountSell, nil
}

// Calculate amountToSell and amountToBuy for SellCoin TX
// In CLI part amountToBuy is minAmountToBuy
func SellCoinCalculateAmounts(coinToBuy types.Coin, coinToSell types.Coin, wantsBuy sdk.Int, wantsSell sdk.Int) (amountBuy sdk.Int, amountSell sdk.Int, err sdkerrors.Error) {
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
		return sdk.Int{}, sdk.Int{}, sdkerrors.New(types.DefaultCodespace, types.TxBreaksVolumeLimit, "Tx breaks LimitVolume rule")
	}

	coinToSellMinReserve := formulas.GetReserveLimitFromCRR(coinToSell.CRR)
	if coinToSell.Reserve.Sub(amountSellInBase).LT(coinToSellMinReserve) && !coinToSell.IsBase() {
		return sdk.Int{}, sdk.Int{}, sdkerrors.New(types.DefaultCodespace, types.TxBreaksMinReserveLimit, "Tx breaks MinReserveLimit rule")
	}

	// Limit minAmountToBuy in CLI
	if amountBuy.LT(wantsBuy) {
		return sdk.Int{}, sdk.Int{}, sdkerrors.New(types.DefaultCodespace, types.AmountBuyIsTooSmall, "Amount you will receive less than minimum")
	}

	return amountBuy, wantsSell, nil
}
