package rest

import (
	types2 "bitbucket.org/decimalteam/go-node/x/coin/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(ctx client.Context, r *mux.Router) {
	r.HandleFunc("/coins", getCoinsHandlerFunc(ctx)).Methods("GET")
	r.HandleFunc("/coin/{symbol}", getCoinHandlerFunc(ctx)).Methods("GET")
}

func getCoinsHandlerFunc(ctx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, ctx, r)
		if !ok {
			return
		}
		res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/coin/%s", types2.QueryListCoins), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, ctx, res)
	}
}

func getCoinHandlerFunc(ctx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, ctx, r)
		if !ok {
			return
		}
		symbol := mux.Vars(r)["symbol"]

		res, _, err := ctx.QueryWithData(fmt.Sprintf("custom/coin/%s/%s", types2.QueryGetCoin, symbol), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, "{\"msg\":\"Coin not found\"}")
			return
		}
		rest.PostProcessResponse(w, ctx, res)
	}
}
