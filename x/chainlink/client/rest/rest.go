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

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("feed/data/{feedId}"+types.QueryFeedData, listFeedDataHandler(clientCtx)).Methods(MethodGet)
}

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	r.Handle("feed/data/", createFeedHandler(clientCtx)).Methods(MethodPOST)
}
