// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc("/chainlink/legacy/feed/data/round/{roundId}/{feedId}", listRoundFeedDataHandler(clientCtx)).Methods(MethodGet) // query feed data by roundId and feedId
	r.HandleFunc("/chainlink/legacy/feed/data/latest/{feedId}", listLatestFeedDataHandler(clientCtx)).Methods(MethodGet)         // query the latest feed data by feedId
	r.HandleFunc("/chainlink/legacy/module/owner", getModuleOwner(clientCtx)).Methods(MethodGet)                                 // query the module owners
	r.HandleFunc("/chainlink/legacy/module/feed/{feedId}", getFeedInfo(clientCtx)).Methods(MethodGet)                            // query the feed info by feedId
	r.HandleFunc("/chainlink/legacy/module/account/{accountAddress}", getAccountInfo(clientCtx)).Methods(MethodGet)              // query the chainlink account
	r.HandleFunc("/chainlink/legacy/module/feed/reward/strategy", getFeedRewardAvailStrategy(clientCtx)).Methods(MethodGet)      // query the available feed reward strategies
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

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%d/%s", types.QuerierRoute, types.QueryRoundFeedData, roundIdInt, feedId), nil)
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

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryLatestFeedData, feedId), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getModuleOwner(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryModuleOwner), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getFeedInfo(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		feedId := vars["feedId"]

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryFeedInfo, feedId), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getAccountInfo(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accAddr := vars["accountAddress"]

		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryAccountInfo, accAddr), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func getFeedRewardAvailStrategy(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		availStrategies := make([]string, 0, len(types.FeedRewardStrategyConvertor))
		for name := range types.FeedRewardStrategyConvertor {
			availStrategies = append(availStrategies, name)
		}

		_, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryFeedRewardStrategy), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, availStrategies)
	}
}
