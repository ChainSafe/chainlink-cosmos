package keeper

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"strconv"

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
		default:
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}

		return res, err
	}
}

func getRoundFeedData(ctx sdk.Context, path []string, keeper Keeper, legacQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	if len(path) < 2 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			"Insufficient parameters, at least 2 parameters is required")
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
