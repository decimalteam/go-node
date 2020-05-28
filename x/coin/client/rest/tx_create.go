package rest

import (
	cliUtils "bitbucket.org/decimalteam/go-node/x/coin/client/utils"
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"bitbucket.org/decimalteam/go-node/x/validator/client/utils/rest"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"net/http"
	"strconv"
)

type CoinCreateReq struct {
	BaseReq              rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title                string       `json:"title" yaml:"title"`                                   // Full coin title (Bitcoin)
	ConstantReserveRatio string       `json:"constant_reserve_ratio" yaml:"constant_reserve_ratio"` // between 10 and 100
	Symbol               string       `json:"symbol" yaml:"symbol"`                                 // Short coin title (BTC)
	InitialVolume        string       `json:"initial_volume" yaml:"initial_volume"`
	InitialReserve       string       `json:"initial_reserve" yaml:"initial_reserve"`
	LimitVolume          string       `json:"limit_volume" yaml:"limit_volume"` // How many coins can be issued
}

func CoinCreateRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CoinCreateReq

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
		var title = req.Title
		var symbol = req.Symbol
		crr, err := strconv.ParseUint(req.ConstantReserveRatio, 10, 8)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Failed to convert CRR to uint")
			return
		}
		var initReserve, _ = sdk.NewIntFromString(req.InitialReserve)
		var initVolume, _ = sdk.NewIntFromString(req.LimitVolume)
		var limitVolume, _ = sdk.NewIntFromString(req.LimitVolume)

		msg := types.NewMsgCreateCoin(title, uint(crr), symbol, initVolume, initReserve, limitVolume, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// Check if coin does not exist yet
		coinExists, _ := cliUtils.ExistsCoin(cliCtx, symbol)
		if coinExists {
			rest.WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Coin with symbol %s already exists", symbol))
			return
		}

		rest.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
