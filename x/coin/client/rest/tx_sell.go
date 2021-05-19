package rest

import (
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/decimalteam/go-node/utils/formulas"
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

type CoinSellReq struct {
	BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
	CoinToSell   string       `json:"coin_to_sell" yaml:"coin_to_sell"`
	AmountToSell string       `json:"amount_to_sell" yaml:"amount_to_sell"`
	CoinToBuy    string       `json:"coin_to_buy" yaml:"coin_to_buy"`
}

func CoinSellRequestHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CoinSellReq

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
		var coinToSellSymbol = req.CoinToSell
		var amountToSell, _ = sdk.NewIntFromString(req.AmountToSell)

		var coinToBuySymbol = req.CoinToBuy

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

		valueSell := formulas.CalculateSaleReturn(coinToSell.Volume, coinToSell.Reserve, coinToSell.CRR, amountToSell)

		// Do basic validating
		msg := types.NewMsgSellCoin(addr, sdk.NewCoin(coinToSellSymbol, valueSell), sdk.NewCoin(coinToBuySymbol, amountToSell))
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Get account balance
		acc, _ := cliUtils.GetAccount(cliCtx, addr)
		balance := acc.GetCoins()
		if balance.AmountOf(strings.ToLower(coinToSellSymbol)).LT(amountToSell) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Not enough coin to sell")
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
