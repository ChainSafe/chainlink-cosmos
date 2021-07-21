// Copyright 2021 ChainSafe Systems
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

	DataProviderSetChangeTypeAdd          = "Add"
	DataProviderSetChangeTypeRemove       = "Remove"
	FeedParamChangeTypeSubmissionCount    = "SubmissionCount"
	FeedParamChangeTypeHeartbeat          = "Heartbeat"
	FeedParamChangeTypeDeviationThreshold = "DeviationThreshold"
	FeedParamChangeTypeRewardSchema       = "RewardSchema"
)

// SubmitFeedDataTx implements the tx/SubmitFeedDataTx gRPC method
func (k Keeper) SubmitFeedDataTx(c context.Context, msg *types.MsgFeedData) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	height, txHash, err := k.SetFeedData(ctx, msg)
	if err != nil {
		return nil, err
	}
	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
	}

	// reward distribution
	feed := k.GetFeed(ctx, msg.FeedId)
	feedReward := feed.GetFeed().FeedReward

	dataProviders := feed.GetFeed().DataProviders

	err = k.DistributeReward(ctx, msg, dataProviders, feedReward)
	if err != nil {
		return nil, err
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

	// emit ModuleOwnershipTransfer event
	err := types.EmitEvent(&types.MsgModuleOwnershipTransferEvent{
		NewModuleOwnerAddr: msg.GetNewModuleOwnerAddress(),
		Signer:             msg.GetAssignerAddress(),
	}, ctx.EventManager())
	if err != nil {
		return nil, err
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

	// emit NewFeed event
	err := types.EmitEvent(&types.MsgNewFeedEvent{
		FeedId:        msg.GetFeedId(),
		DataProviders: msg.GetDataProviders(),
		FeedOwner:     msg.GetFeedOwner(),
	}, ctx.EventManager())
	if err != nil {
		return nil, err
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

	// emit DataProviderSetChange event
	err = types.EmitEvent(&types.MsgDataProviderSetChangeEvent{
		FeedId:           msg.GetFeedId(),
		ChangeType:       DataProviderSetChangeTypeAdd,
		DataProviderAddr: msg.GetDataProvider().GetAddress(),
		Signer:           msg.GetSigner(),
	}, ctx.EventManager())
	if err != nil {
		return nil, err
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

	// emit DataProviderSetChange event
	err = types.EmitEvent(&types.MsgDataProviderSetChangeEvent{
		FeedId:           msg.GetFeedId(),
		ChangeType:       DataProviderSetChangeTypeRemove,
		DataProviderAddr: msg.GetAddress(),
		Signer:           msg.GetSigner(),
	}, ctx.EventManager())
	if err != nil {
		return nil, err
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

	// emit FeedParameterChange event
	err = types.EmitEvent(&types.MsgFeedParameterChangeEvent{
		FeedId:            msg.GetFeedId(),
		ChangeType:        FeedParamChangeTypeSubmissionCount,
		NewParameterValue: msg.GetSubmissionCount(),
		Signer:            msg.GetSigner(),
	}, ctx.EventManager())
	if err != nil {
		return nil, err
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

	// emit FeedParameterChange event
	err = types.EmitEvent(&types.MsgFeedParameterChangeEvent{
		FeedId:            msg.GetFeedId(),
		ChangeType:        FeedParamChangeTypeHeartbeat,
		NewParameterValue: msg.GetHeartbeatTrigger(),
		Signer:            msg.GetSigner(),
	}, ctx.EventManager())
	if err != nil {
		return nil, err
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

	// emit FeedParameterChange event
	err = types.EmitEvent(&types.MsgFeedParameterChangeEvent{
		FeedId:            msg.GetFeedId(),
		ChangeType:        FeedParamChangeTypeDeviationThreshold,
		NewParameterValue: msg.GetDeviationThresholdTrigger(),
		Signer:            msg.GetSigner(),
	}, ctx.EventManager())
	if err != nil {
		return nil, err
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}

func (k Keeper) SetFeedRewardTx(c context.Context, msg *types.MsgSetFeedReward) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := k.SetFeedReward(ctx, msg)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
	}

	// emit FeedParameterChange event
	err = types.EmitEvent(&types.MsgFeedParameterChangeEvent{
		FeedId:            msg.GetFeedId(),
		ChangeType:        FeedParamChangeTypeRewardSchema,
		NewParameterValue: msg.GetFeedReward(),
		Signer:            msg.GetSigner(),
	}, ctx.EventManager())
	if err != nil {
		return nil, err
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

	// emit FeedOwnershipTransfer event
	err = types.EmitEvent(&types.MsgFeedOwnershipTransferEvent{
		FeedId:           msg.GetFeedId(),
		NewFeedOwnerAddr: msg.GetNewFeedOwnerAddress(),
		Signer:           msg.GetSigner(),
	}, ctx.EventManager())
	if err != nil {
		return nil, err
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}
