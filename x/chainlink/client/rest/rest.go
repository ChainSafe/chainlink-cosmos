// SPDX-License-Identifier: MIT

package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/gorilla/mux"
)

const (
	MethodGet  = "GET"
	MethodPUT  = "PUT"
	MethodPOST = "POST"
)

// RegisterRoutes registers blog-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
}

func registerQueryRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := clientrest.WithHTTPDeprecationHeaders(rtr)

	r.HandleFunc("/chainlink/feed/data/round/{roundId}/{feedId}", listRoundFeedDataHandler(clientCtx)).Methods(MethodGet) // query feed data by roundId and feedId
	r.HandleFunc("/chainlink/feed/data/latest/{feedId}", listLatestFeedDataHandler(clientCtx)).Methods(MethodGet)         // query the latest feed data by feedId
}

func registerTxHandlers(clientCtx client.Context, rtr *mux.Router) {
	r := clientrest.WithHTTPDeprecationHeaders(rtr)

	r.HandleFunc("/txs/encode", authrest.EncodeTxRequestHandlerFn(clientCtx)).Methods(MethodPOST)
	r.HandleFunc("/txs/decode", authrest.DecodeTxRequestHandlerFn(clientCtx)).Methods(MethodPOST)
	r.HandleFunc("/txs", authrest.BroadcastTxRequest(clientCtx)).Methods(MethodPOST)

	r.Handle("/chainlink/feed/data", submitFeedDataHandler(clientCtx)).Methods(MethodPUT)
}

func listRoundFeedDataHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		roundId := vars["roundId"]
		roundIdInt, err := strconv.ParseInt(roundId, 10, 16)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("roundId is invalid").Error())
			return
		}
		feedId := vars["feedId"]

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d/%s", types.FeedDataStoreKey, types.QueryRoundFeedData, roundIdInt, feedId), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}
func listLatestFeedDataHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		feedId := vars["feedId"]

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.FeedDataStoreKey, types.QueryLatestFeedData, feedId), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func submitFeedDataHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createFeedRequest
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

		msg := types.NewMsgFeedData(submitter, req.FeedId, req.FeedData, req.Signatures)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
