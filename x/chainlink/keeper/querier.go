// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package keeper

import (
	"strconv"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	abci "github.com/tendermint/tendermint/abci/types"
)

// Legacy querier for backwards compatibility

const defaultPageLimit = 10

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		var (
			res []byte
			err error
		)

		switch path[0] {
		case types.QueryRoundFeedData:
			return getRoundFeedData(ctx, path, k, legacyQuerierCdc)
		case types.QueryLatestFeedData:
			return latestRoundFeedData(ctx, path, k, legacyQuerierCdc)
		case types.QueryModuleOwner:
			return getModuleOwners(ctx, path, k, legacyQuerierCdc)
		case types.QueryFeedInfo:
			return getFeedInfo(ctx, path, k, legacyQuerierCdc)
		case types.QueryAccountInfo:
			return getAccountInfo(ctx, path, k, legacyQuerierCdc)
		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}

func getRoundFeedData(ctx sdk.Context, path []string, keeper Keeper, legacQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	if len(path) < 3 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 3 parameters is required")
	}
	roundId, err := strconv.ParseUint(path[1], 10, 64)
	if err != nil {
		return nil, err
	}
	feedId := path[2]

	req := &types.GetRoundDataRequest{
		FeedId:     feedId,
		RoundId:    roundId,
		Pagination: &query.PageRequest{Limit: defaultPageLimit},
	}
	msgs, err := keeper.GetRoundFeedDataByFilter(ctx, req)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacQuerierCdc, msgs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func latestRoundFeedData(ctx sdk.Context, path []string, keeper Keeper, legacQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}
	feedId := path[1]

	req := &types.GetLatestRoundDataRequest{
		FeedId: feedId,
	}
	msgs, err := keeper.GetLatestRoundFeedDataByFilter(ctx, req)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacQuerierCdc, msgs)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func getModuleOwners(ctx sdk.Context, path []string, keeper Keeper, legacQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	if len(path) < 1 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 1 parameters is required")
	}

	resp := keeper.GetModuleOwnerList(ctx)

	bz, err := codec.MarshalJSONIndent(legacQuerierCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func getFeedInfo(ctx sdk.Context, path []string, keeper Keeper, legacQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}
	feedId := path[1]

	resp := keeper.GetFeed(ctx, feedId)
	if resp.Feed == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "No feed found")
	}

	bz, err := codec.MarshalJSONIndent(legacQuerierCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func getAccountInfo(ctx sdk.Context, path []string, keeper Keeper, legacQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
	}

	accAddrString := path[1]
	accAddr, err := sdk.AccAddressFromBech32(accAddrString)
	if err != nil {
		return nil, err
	}

	req := &types.GetAccountRequest{AccountAddress: accAddr}
	resp := keeper.GetAccount(ctx, req)

	bz, err := codec.MarshalJSONIndent(legacQuerierCdc, resp)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
