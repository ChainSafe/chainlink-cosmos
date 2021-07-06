package keeper

import (
	"context"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ types.MsgServer = Keeper{}

// SubmitFeedDataTx implements the tx/SubmitFeedDataTx gRPC method
func (k Keeper) SubmitFeedDataTx(c context.Context, msg *types.MsgFeedData) (*types.MsgResponse, error) {
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

// AddModuleOwnerTx implements the tx/AddModuleOwnerTx gRPC method
func (k Keeper) AddModuleOwnerTx(c context.Context, msg *types.MsgModuleOwner) (*types.MsgResponse, error) {
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

// ModuleOwnershipTransferTx implements the tx/ModuleOwnershipTransferTx gRPC method
func (k Keeper) ModuleOwnershipTransferTx(c context.Context, msg *types.MsgModuleOwnershipTransfer) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_, _ = k.RemoveModuleOwner(ctx, msg)

	transferMsg := &types.MsgModuleOwner{
		Address:         msg.GetNewModuleOwnerAddress(),
		PubKey:          msg.GetNewModuleOwnerPubKey(),
		AssignerAddress: msg.GetAssignerAddress(),
	}
	height, txHash := k.SetModuleOwner(ctx, transferMsg)

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, "incorrect height found")
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}
