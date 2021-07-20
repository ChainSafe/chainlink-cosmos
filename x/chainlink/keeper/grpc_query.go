package keeper

import (
	"context"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = Keeper{}

// GetRoundData implements the Query/GetRoundData gRPC method
func (s Keeper) GetRoundData(c context.Context, req *types.GetRoundDataRequest) (*types.GetRoundDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return s.GetRoundFeedDataByFilter(ctx, req)
}

// LatestRoundData implements the Query/LatestRoundData gRPC method
func (s Keeper) LatestRoundData(c context.Context, req *types.GetLatestRoundDataRequest) (*types.GetLatestRoundDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return s.GetLatestRoundFeedDataByFilter(ctx, req)
}

// GetAllModuleOwner implements the Query/GetAllModuleOwner gRPC method
func (s Keeper) GetAllModuleOwner(c context.Context, _ *types.GetModuleOwnerRequest) (*types.GetModuleOwnerResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return s.GetModuleOwnerList(ctx), nil
}

func (s Keeper) GetFeedByFeedId(c context.Context, req *types.GetFeedByIdRequest) (*types.GetFeedByIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return s.GetFeed(ctx, req.FeedId), nil
}
