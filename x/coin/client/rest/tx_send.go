package rest

import (
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"net/http"
	"strings"
)

type CoinSendReq struct {
	BaseReq  rest.BaseReq `json:"base_req" yaml:"base_req"`
	Coin     string       `json:"coin" yaml:"coin"`
	Amount   string       `json:"amount" yaml:"amount"`
	Receiver string       `json:"receiver" yaml:"receiver"`
}

func CoinSendRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CoinSendReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
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
		msg := types.NewMsgSendCoin(addr, coin, amount, receiver)
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
		acc, _ := cliUtils.GetAccount(cliCtx, cliCtx.GetFromAddress())
		balance := acc.GetCoins()
		if balance.AmountOf(strings.ToLower(coin)).LT(amount) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Not enough coin to send")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
