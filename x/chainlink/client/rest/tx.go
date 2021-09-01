// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package rest

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/gorilla/mux"
	"net/http"
)

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/txs/encode", authrest.EncodeTxRequestHandlerFn(clientCtx)).Methods(MethodPOST)
	r.HandleFunc("/txs/decode", authrest.DecodeTxRequestHandlerFn(clientCtx)).Methods(MethodPOST)
	r.HandleFunc("/txs", authrest.BroadcastTxRequest(clientCtx)).Methods(MethodPOST)

	r.HandleFunc("/chainlink/feed/data", NewFeedDataRequestHandler(clientCtx)).Methods(MethodPUT)
}

type FeedDataRequest struct {
	BaseReq       rest.BaseReq `json:"baseReq"`
	FeedId        string       `json:"feedId"`
	FeedData      [][]byte     `json:"feedData"`
	Signatures    [][]byte     `json:"signature"`
	CosmosPubKeys [][]byte     `json:"cosmosPubKeys"`
}

func NewFeedDataRequestHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req FeedDataRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		submitter, err := sdk.AccAddressFromBech32(baseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgFeedData(submitter, req.FeedId, req.FeedData, req.Signatures, req.CosmosPubKeys)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
