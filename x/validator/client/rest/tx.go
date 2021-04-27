package rest

import (
	"bitbucket.org/decimalteam/go-node/x/validator/internal/types"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/ripemd160"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/validator/recover_address",
		recoverAddress(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/validator/delegators/{delegatorAddr}/delegations",
		postDelegationsHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/validator/delegators/{delegatorAddr}/unbonding_delegations",
		postUnbondingDelegationsHandlerFn(cliCtx),
	).Methods("POST")
}

// HashLength represents fixed hash length.
const HashLength = 32

type (
	// Hash represents the 32 byte Keccak256 hash of arbitrary data.
	Hash [HashLength]byte

	// DelegateRequest defines the properties of a delegation request's body.
	DelegateRequest struct {
		BaseReq          rest.BaseReq   `json:"base_req" yaml:"base_req"`
		DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"` // in bech32
		ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"` // in bech32
		Amount           sdk.Coin       `json:"amount" yaml:"amount"`
	}

	// UndelegateRequest defines the properties of a undelegate request's body.
	UndelegateRequest struct {
		BaseReq          rest.BaseReq   `json:"base_req" yaml:"base_req"`
		DelegatorAddress sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"` // in bech32
		ValidatorAddress sdk.ValAddress `json:"validator_address" yaml:"validator_address"` // in bech32
		Amount           sdk.Coin       `json:"amount" yaml:"amount"`
	}

	RecoverAddressRequest struct {
		Hash    Hash         `json:"hash" yaml:"hash"`
		V       big.Int      `json:"v" yaml:"v"`
		R       big.Int      `json:"r" yaml:"r"`
		S       big.Int      `json:"s" yaml:"s"`
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	}
)

func postDelegationsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DelegateRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgDelegate(sdk.ValAddress(req.DelegatorAddress), sdk.AccAddress(req.ValidatorAddress), req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, req.DelegatorAddress) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own delegator address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postUnbondingDelegationsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UndelegateRequest

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgUnbond(sdk.ValAddress(req.DelegatorAddress), sdk.AccAddress(req.ValidatorAddress), req.Amount)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, req.DelegatorAddress) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own delegator address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func recoverAddress(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(1)
		response := []string{}
		var addr sdk.AccAddress
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, 4006, err.Error())
		}

		fmt.Println(2)
		targets := []RecoverAddressRequest{}

		err = json.Unmarshal([]byte(reqBody), &targets)
		if err != nil {
			rest.WriteErrorResponse(w, 4007, err.Error())
		}
		var basReq rest.BaseReq

		for _, msg := range targets {
			basReq = msg.BaseReq
			sighash := msg.Hash
			Vb := msg.V
			R := msg.R
			S := msg.S

			if Vb.BitLen() > 8 {
				rest.WriteErrorResponse(w, 4001, err.Error())
			}
			V := byte(Vb.Uint64() - 27)
			if !crypto.ValidateSignatureValues(V, &R, &S, true) {
				rest.WriteErrorResponse(w, 4002, err.Error())
			}
			// encode the snature in uncompressed format
			r, s := R.Bytes(), S.Bytes()
			sig := make([]byte, 65)
			copy(sig[32-len(r):32], r)
			copy(sig[64-len(s):64], s)
			sig[64] = V
			// recover the public key from the snature
			pub, err := crypto.Ecrecover(sighash[:], sig)
			if err != nil {
				rest.WriteErrorResponse(w, 4003, err.Error())
			}
			if len(pub) == 0 || pub[0] != 4 {
				rest.WriteErrorResponse(w, 4004, err.Error())
			}
			pub2, err := crypto.UnmarshalPubkey(pub)
			if err != nil {
				rest.WriteErrorResponse(w, 4005, err.Error())
			}
			pub3 := crypto.CompressPubkey(pub2)
			hasherSHA256 := sha256.New()
			hasherSHA256.Write(pub3)
			sha := hasherSHA256.Sum(nil)
			hasherRIPEMD160 := ripemd160.New()
			hasherRIPEMD160.Write(sha)
			addresses := sdk.AccAddress(hasherRIPEMD160.Sum(nil))
			addr = addresses
			response = append(response, addresses.String())
		}
		msg := types.NewMsgRecoveredAddress(addr)
		fmt.Println(3)
		utils.WriteGenerateStdTxResponse(w, cliCtx, basReq, []sdk.Msg{msg})
	}
}
