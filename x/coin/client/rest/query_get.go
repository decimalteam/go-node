package rest

import (
	"bitbucket.org/decimalteam/go-node/x/coin/internal/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func queryGetCoinHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		symbol := mux.Vars(r)["symbol"]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/coin/%s/%s", types.QueryGetCoin, symbol), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, "{\"msg\":\"Coin not found\"}")
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
