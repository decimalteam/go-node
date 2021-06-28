package rest

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"net/http"
	"strings"

	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

type CoinSendReq struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req"`
	Coin     string       `json:"coin" yaml:"coin"`
	Amount   string       `json:"amount" yaml:"amount"`
	Receiver string       `json:"receiver" yaml:"receiver"`
}

func CoinSendRequestHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CoinSendReq

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
		coin := req.Coin
		amount, _ := sdk.NewIntFromString(req.Amount)
		receiver, err := sdk.AccAddressFromBech32(req.Receiver)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types2.NewMsgSendCoin(addr, sdk.NewCoin(coin, amount), receiver)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// Check if coin exists
		existsCoin, _ := cliUtils.ExistsCoin(cliCtx, coin)
		print(err)
		if !existsCoin {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Coin to sent with symbol %s does not exist", coin))
			return
		}

		// Check if enough balance
		acc, err := cliUtils.GetAccount(cliCtx, addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		balance, _ := cliUtils.GetAccountCoins(cliCtx, acc.GetAddress())
		if balance.AmountOf(strings.ToLower(coin)).LT(amount) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Not enough coin to send")
			return
		}

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, &msg)
	}
}
