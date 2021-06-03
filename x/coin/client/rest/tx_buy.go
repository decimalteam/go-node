package rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	coinTypes "bitbucket.org/decimalteam/go-node/x/coin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CoinBuyReq struct {
	BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
	CoinToSell   string       `json:"coin_to_sell" yaml:"coin_to_sell"`
	AmountToSell string       `json:"amount_to_sell" yaml:"amount_to_sell"`
	CoinToBuy    string       `json:"coin_to_buy" yaml:"coin_to_buy"`
}

func CoinBuyRequestHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CoinBuyReq

		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq

		addr, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		var coinToBuySymbol = req.CoinToBuy

		var coinToSellSymbol = req.CoinToSell
		var amountToSell, _ = sdk.NewIntFromString(req.AmountToSell)

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

		valueSell := formulas.CalculatePurchaseReturn(coinToBuy.Volume, coinToBuy.Reserve, coinToBuy.CRR, amountToSell)

		// Do basic validating
		msg := coinTypes.NewMsgBuyCoin(addr, sdk.NewCoin(coinToBuySymbol, valueSell), sdk.NewCoin(coinToSellSymbol, amountToSell))
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Get account balance
		acc, _ := cliUtils.GetAccount(cliCtx, addr)
		balance, _ := cliUtils.GetAccountCoins(cliCtx, acc.GetAddress())

		if balance.AmountOf(strings.ToLower(coinToSellSymbol)).LT(amountToSell) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Not enough coin to sell")
			return
		}
		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, []sdk.Msg{&msg}...)
	}
}
