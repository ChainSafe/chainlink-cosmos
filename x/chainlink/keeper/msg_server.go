// SPDX-License-Identifier: MIT

package keeper

import (
	"context"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ types.MsgServer = Keeper{}

const (
	ErrIncorrectHeightFound = "incorrect height found"
)

// SubmitFeedDataTx implements the tx/SubmitFeedDataTx gRPC method
func (k Keeper) SubmitFeedDataTx(c context.Context, msg *types.MsgFeedData) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	height, txHash := k.SetFeedData(ctx, msg)

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
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
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
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
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}

// AddFeedTx implements the tx/AddFeedTx gRPC method
func (k Keeper) AddFeedTx(c context.Context, msg *types.MsgFeed) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash := k.SetFeed(ctx, msg)

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}

// AddDataProviderTx implements the tx/AddDataProvider gRPC method
func (k Keeper) AddDataProviderTx(c context.Context, msg *types.MsgAddDataProvider) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := k.AddDataProvider(ctx, msg)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}

// RemoveDataProviderTx implements the tx/RemoveDataProvider gRPC method
func (k Keeper) RemoveDataProviderTx(c context.Context, msg *types.MsgRemoveDataProvider) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := k.RemoveDataProvider(ctx, msg)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}

func (k Keeper) SetSubmissionCountTx(c context.Context, msg *types.MsgSetSubmissionCount) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := k.SetSubmissionCount(ctx, msg)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}

func (k Keeper) SetHeartbeatTriggerTx(c context.Context, msg *types.MsgSetHeartbeatTrigger) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := k.SetHeartbeatTrigger(ctx, msg)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}

func (k Keeper) SetDeviationThresholdTriggerTx(c context.Context, msg *types.MsgSetDeviationThresholdTrigger) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := k.SetDeviationThresholdTrigger(ctx, msg)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}

func (k Keeper) FeedOwnershipTransferTx(c context.Context, msg *types.MsgFeedOwnershipTransfer) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := k.FeedOwnershipTransfer(ctx, msg)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}
