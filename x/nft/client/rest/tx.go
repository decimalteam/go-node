package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/gorilla/mux"

	"bitbucket.org/decimalteam/go-node/x/nft/internal/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router,
	cdc *codec.Codec, queryRoute string) {
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

	//Update Reserv NFT
	r.HandleFunc(
		"/nfts/collection/{denom}/nft/{id}/updateReserve",
		updateReserveNFTHandler(cdc, cliCtx),
	).Methods("PUT")
}

type transferNFTReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Denom       string       `json:"denom"`
	ID          string       `json:"id"`
	Recipient   string       `json:"recipient"`
	SubTokenIDs []string     `json:"subTokenIDs"`
}

func transferNFTHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
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

		subTokenIDs := make([]int64, len(req.SubTokenIDs))
		for i, d := range req.SubTokenIDs {
			subTokenID, err := strconv.ParseInt(d, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid subTokenID")
				return
			}
			subTokenIDs[i] = subTokenID
		}

		// create the message
		msg := types.NewMsgTransferNFT(fromAddr, recipient, req.Denom, req.ID, subTokenIDs)

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type editNFTMetadataReq struct {
	BaseReq  rest.BaseReq `json:"base_req"`
	Denom    string       `json:"denom"`
	ID       string       `json:"id"`
	TokenURI string       `json:"tokenURI"`
}

func editNFTMetadataHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
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
		msg := types.NewMsgEditNFTMetadata(fromAddr, req.ID, req.Denom, req.TokenURI)

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
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

func mintNFTHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
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
		msg := types.NewMsgMintNFT(fromAddr, req.Recipient, req.ID, req.Denom, req.TokenURI, quantity, sdk.NewInt(1), false)

		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type burnNFTReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	Denom       string       `json:"denom"`
	ID          string       `json:"id"`
	SubTokenIDs []string     `json:"subTokenIDs"`
}

func burnNFTHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
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

		subTokenIDs := make([]int64, len(req.SubTokenIDs))
		for i, d := range req.SubTokenIDs {
			subTokenID, err := strconv.ParseInt(d, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid subTokenID")
				return
			}
			subTokenIDs[i] = subTokenID
		}

		// create the message
		msg := types.NewMsgBurnNFT(fromAddr, req.ID, req.Denom, subTokenIDs)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}

type MsgUpdateReserveNFTq struct {
	BaseReq      rest.BaseReq `json:"base_req"`
	ID           string       `json:"id"`
	Denom        string       `json:"denom"`
	SubTokenIDs  []string     `json:"sub_token_ids"`
	NewReserveNFT string       `json:"new_reserve"`
}

func updateReserveNFTHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MsgUpdateReserveNFTq
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

		subTokenIDs := make([]int64, len(req.SubTokenIDs))
		for i, d := range req.SubTokenIDs {
			subTokenID, err := strconv.ParseInt(d, 10, 64)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid subTokenID")
				return
			}
			subTokenIDs[i] = subTokenID
		}

		newReserve, ok := sdk.NewIntFromString(req.NewReserveNFT)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid quantity")
			return
		}
		// create the message
		msg := types.NewMsgUpdateReserveNFT(fromAddr, req.ID, req.Denom, subTokenIDs, newReserve)
		fmt.Println(req.Denom)
		fmt.Println(msg.Denom)
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
