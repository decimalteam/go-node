package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/coin/{symbol}",
		queryGetCoinHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/coins",
		queryGetCoinListHandlerFn(cliCtx),
	).Methods("GET")
}
