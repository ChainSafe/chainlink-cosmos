package keeper

import (
	"context"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.MsgServer = Keeper{}

// SubmitFeedData implements the tx/SubmitFeedData gRPC method
func (k Keeper) SubmitFeedData(c context.Context, msg *types.MsgFeedData) (*types.MsgFeedDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	k.SetFeedData(ctx, msg)

	// TODO: how to return txHash and height here?
	return &types.MsgFeedDataResponse{
		Height: uint64(ctx.BlockHeight()),
		TxHash: string(ctx.TxBytes()),
	}, nil
}
