package keeper

import (
	"context"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllFeedData(c context.Context, req *types.QueryAllFeedDataRequest) (*types.QueryAllFeedDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var feeds []*types.MsgFeed
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	feedStore := prefix.NewStore(store, types.KeyPrefix(types.FeedKey))

	pageRes, err := query.Paginate(feedStore, req.Pagination, func(key []byte, value []byte) error {
		var feed types.MsgFeed
		if err := k.cdc.UnmarshalBinaryBare(value, &feed); err != nil {
			return err
		}

		feeds = append(feeds, &feed)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllFeedDataResponse{
		FeedData:   feeds,
		Pagination: pageRes,
	}, nil
}
