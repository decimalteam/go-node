package rest

import (
	types2 "bitbucket.org/decimalteam/go-node/x/nft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"net/http"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/gorilla/mux"
)

func registerTxRoutes(cliCtx client.Context, r *mux.Router,
	cdc *codec.LegacyAmino, queryRoute string) {
	// Transfer an NFT to an address
	r.HandleFunc(
		"/nfts/transfer",
		transferNFTHandler(cdc, cliCtx),
	).Methods("POST")

	// Update an NFT metadata
	r.HandleFunc(
		"/nfts/collection/{denom}/nft/{id}/metadata",
		editNFTMetadataHandler(cdc, cliCtx),
	).Methods("PUT")

	// Mint an NFT
	r.HandleFunc(
		"/nfts/mint",
		mintNFTHandler(cdc, cliCtx),
	).Methods("POST")

	// Burn an NFT
	r.HandleFunc(
		"/nfts/collection/{denom}/nft/{id}/burn",
		burnNFTHandler(cdc, cliCtx),
	).Methods("PUT")
}

type transferNFTReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Denom     string       `json:"denom"`
	ID        string       `json:"id"`
	Recipient string       `json:"recipient"`
	Quantity  string       `json:"quantity"`
}

func transferNFTHandler(cdc *codec.LegacyAmino, cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req transferNFTReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		recipient, err := sdk.AccAddressFromBech32(req.Recipient)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		quantity, ok := sdk.NewIntFromString(req.Quantity)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid quantity")
			return
		}

		// create the message
		msg := types2.NewMsgTransferNFT(fromAddr, recipient, req.Denom, req.ID, quantity)

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, &msg)
	}
}

type editNFTMetadataReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Denom    string       `json:"denom"`
	ID       string       `json:"id"`
	TokenURI string       `json:"tokenURI"`
}

func editNFTMetadataHandler(cdc *codec.LegacyAmino, cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req editNFTMetadataReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types2.NewMsgEditNFTMetadata(fromAddr, req.ID, req.Denom, req.TokenURI)

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, &msg)
	}
}

type mintNFTReq struct {
	BaseReq   rest.BaseReq   `json:"base_req"`
	Recipient sdk.AccAddress `json:"recipient"`
	Denom     string         `json:"denom"`
	ID        string         `json:"id"`
	TokenURI  string         `json:"tokenURI"`
	Quantity  string         `json:"quantity"`
}

func mintNFTHandler(cdc *codec.LegacyAmino, cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req mintNFTReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		quantity, ok := sdk.NewIntFromString(req.Quantity)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid quantity")
			return
		}

		// create the message
		msg := types2.NewMsgMintNFT(fromAddr, req.Recipient, req.ID, req.Denom, req.TokenURI, quantity, sdk.NewInt(1), false)

		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, &msg)
	}
}

type burnNFTReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Denom    string       `json:"denom"`
	ID       string       `json:"id"`
	Quantity string       `json:"quantity"`
}

func burnNFTHandler(cdc *codec.LegacyAmino, cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req burnNFTReq
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		quantity, ok := sdk.NewIntFromString(req.Quantity)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid quantity")
			return
		}

		// create the message
		msg := types2.NewMsgBurnNFT(fromAddr, req.ID, req.Denom, quantity)
		tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, &msg)
	}
}
