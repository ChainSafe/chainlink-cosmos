package keeper

import (
	"context"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ types.MsgServer = Keeper{}

// SubmitFeedData implements the tx/SubmitFeedData gRPC method
func (k Keeper) SubmitFeedData(c context.Context, msg *types.MsgFeedData) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	height, txHash := k.SetFeedData(ctx, msg)

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, "incorrect height found")
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}

// AddModuleOwner implements the tx/AddModuleOwner gRPC method
func (k Keeper) AddModuleOwner(c context.Context, msg *types.ModuleOwner) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	height, txHash := k.SetModuleOwner(ctx, msg)

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, "incorrect height found")
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}
