// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/gorilla/mux"
)

const (
	MethodGet  = "GET"
	MethodPUT  = "PUT"
	MethodPOST = "POST"
)

// RegisterRoutes registers blog-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, rtr *mux.Router) {
	r := clientrest.WithHTTPDeprecationHeaders(rtr)

	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)
}
