package rest

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	// this line is used by starport scaffolding # 1
)

const (
	MethodGet  = "GET"
	MethodPOST = "post"
)

// RegisterRoutes registers blog-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
}

// nolint
func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("feed/"+types.QueryListFeed, listFeedHandler(clientCtx)).Methods(MethodGet)
}

// nolint
func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	r.Handle("feed/", createFeedHandler(clientCtx)).Methods(MethodPOST)
}
