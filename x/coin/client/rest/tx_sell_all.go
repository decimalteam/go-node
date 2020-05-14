package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	decsdk "bitbucket.org/decimalteam/go-node/utils/types"
	"bitbucket.org/decimalteam/go-node/x/auth/client/utils"
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
)

type CoinSellAllReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	CoinToSell  string       `json:"coin_to_sell" yaml:"coin_to_sell"`
	CoinToBuy   string       `json:"coin_to_buy" yaml:"coin_to_buy"`
	AmountToBuy sdk.Int      `json:"amount_to_buy" yaml:"amount_to_buy"`
}

func CoinSellAllRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CoinSellAllReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq

		//addr, err := decsdk.AccAddressFromPrefixedHex(baseReq.From)
		//if err != nil {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		//	return
		//}
		var coinToSellSymbol = req.CoinToSell
		var coinToBuySymbol = req.CoinToBuy
		var amountToBuy = req.AmountToBuy

		// Check if coin to buy exists
		coinToBuy, _ := cliUtils.GetCoin(cliCtx, coinToBuySymbol)
		if coinToBuy.Symbol != coinToBuySymbol {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Coin to buy with symbol %s does not exist", coinToBuySymbol))
			return
		}
		// Check if coin to sell exists
		coinToSell, _ := cliUtils.GetCoin(cliCtx, coinToSellSymbol)
		if coinToSell.Symbol != coinToSellSymbol {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Coin to sell with symbol %s does not exist", coinToSellSymbol))
			return
		}
		// TODO: Validate limits and check if sufficient balance (formulas)
		//acc, _ := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
		//balance := acc.GetCoins()
		//amountToSell := balance.AmountOf(strings.ToLower(coinToSellSymbol))
		//amountToBuy := formulas.CalculateSaleReturn()

		//_, _, err = cliUtils.SellCoinCalculateAmounts(coinToBuy, coinToSell, minAmountToBuy, amountToSell)
		//if err != nil {
		//	return err
		//}
		// Do basic validating
		msg := types.NewMsgSellAllCoin(decsdk.AccAddress(cliCtx.GetFromAddress()), coinToBuySymbol, coinToSellSymbol, amountToBuy)
		//err = msg.ValidateBasic()
		//if err != nil {
		//	return err
		//}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
