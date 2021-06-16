package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
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

	r.HandleFunc("/chainlink/feed/data/{feedId}", listFeedDataByFeedIdHandler(clientCtx)).Methods(MethodGet) // query feed data by feedId
}

func registerTxHandlers(clientCtx client.Context, rtr *mux.Router) {
	r := clientrest.WithHTTPDeprecationHeaders(rtr)

	r.HandleFunc("/txs/encode", authrest.EncodeTxRequestHandlerFn(clientCtx)).Methods(MethodPOST)
	r.HandleFunc("/txs/decode", authrest.DecodeTxRequestHandlerFn(clientCtx)).Methods(MethodPOST)
	r.HandleFunc("/txs", authrest.BroadcastTxRequest(clientCtx)).Methods(MethodPOST)

	r.Handle("/chainlink/feed/data", submitFeedDataHandler(clientCtx)).Methods(MethodPUT)
}
