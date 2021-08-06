// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package keeper

import (
	"context"
	"time"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ts "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	ErrIncorrectHeightFound   = "incorrect height found"
	ErrInsufficientSignatures = "incorrect number of signatures provided"
	ErrHeartBeatTrigger       = "heartbeat interval has not passed"

	DataProviderSetChangeTypeAdd          = "Add"
	DataProviderSetChangeTypeRemove       = "Remove"
	FeedParamChangeTypeSubmissionCount    = "SubmissionCount"
	FeedParamChangeTypeHeartbeat          = "Heartbeat"
	FeedParamChangeTypeDeviationThreshold = "DeviationThreshold"
	FeedParamChangeTypeRewardSchema       = "RewardSchema"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// SubmitFeedDataTx implements the tx/SubmitFeedDataTx gRPC method
func (s msgServer) SubmitFeedDataTx(c context.Context, msg *types.MsgFeedData) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	height, blockTime, txHash, err := s.SetFeedData(ctx, msg)
	if err != nil {
		return nil, err
	}
	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
	}

	feed := s.GetFeed(ctx, msg.FeedId)

	// check if data passes proper checks
	// submissions count check - should have at least the minimum number of signatures present
	feedSubmissionCount := feed.GetFeed().SubmissionCount
	if uint32(len(msg.Signatures)) < feedSubmissionCount {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, ErrInsufficientSignatures)
	}

	// hearbeat trigger - the interval set in which has to be passed to be valid
	feedHeartbeatTrigger := feed.GetFeed().HeartbeatTrigger
	if blockTime.Before((feed.Feed.LastUpdate.AsTime()).Add(time.Duration(feedHeartbeatTrigger))) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, ErrHeartBeatTrigger)
	}

	// deviation threshold trigger - the fraction of deviation must be met in order to trigger a new round and submit the feed data
	// TODO the feed data submitted is an arbitrary byte array, which could be deserialized to a OCRABIENCODED struct, but there is no price data that can be compared. this will need to be discussed.
	// TODO doesnt seem as though there is a PRICE variable sent in the OCRv2

	// update the lastupdate timestamp of the feed
	msgFeed := feed.Feed
	msgFeed.LastUpdate = ts.New(blockTime)
	h, b := s.SetFeed(ctx, msgFeed)

	// distribute rewards
	// TODO will need to deconstruct each individual signature and make sure that the DP's key are equal.
	// TODO will need to figure out a scheme to correspond the chainlink signature (pubkey) to the cosmos account.
	// TODO only pay out the DP who provides their signature.
	feedReward := feed.GetFeed().FeedReward
	dataProviders := feed.GetFeed().DataProviders
	err = s.DistributeReward(ctx, msg, dataProviders, feedReward)
	if err != nil {
		return nil, err
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}

// AddModuleOwnerTx implements the tx/AddModuleOwnerTx gRPC method
func (s msgServer) AddModuleOwnerTx(c context.Context, msg *types.MsgModuleOwner) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	height, txHash := s.SetModuleOwner(ctx, msg)

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
	}

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}

// ModuleOwnershipTransferTx implements the tx/ModuleOwnershipTransferTx gRPC method
func (s msgServer) ModuleOwnershipTransferTx(c context.Context, msg *types.MsgModuleOwnershipTransfer) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	_, _ = s.RemoveModuleOwner(ctx, msg)

	transferMsg := &types.MsgModuleOwner{
		Address:         msg.GetNewModuleOwnerAddress(),
		PubKey:          msg.GetNewModuleOwnerPubKey(),
		AssignerAddress: msg.GetAssignerAddress(),
	}
	height, txHash := s.SetModuleOwner(ctx, transferMsg)

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
func (s msgServer) AddFeedTx(c context.Context, msg *types.MsgFeed) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash := s.SetFeed(ctx, msg)

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
func (s msgServer) AddDataProviderTx(c context.Context, msg *types.MsgAddDataProvider) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := s.AddDataProvider(ctx, msg)

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
func (s msgServer) RemoveDataProviderTx(c context.Context, msg *types.MsgRemoveDataProvider) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := s.RemoveDataProvider(ctx, msg)

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

func (s msgServer) SetSubmissionCountTx(c context.Context, msg *types.MsgSetSubmissionCount) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := s.SetSubmissionCount(ctx, msg)

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

func (s msgServer) SetHeartbeatTriggerTx(c context.Context, msg *types.MsgSetHeartbeatTrigger) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := s.SetHeartbeatTrigger(ctx, msg)

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

func (s msgServer) SetDeviationThresholdTriggerTx(c context.Context, msg *types.MsgSetDeviationThresholdTrigger) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := s.SetDeviationThresholdTrigger(ctx, msg)

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

func (s msgServer) SetFeedRewardTx(c context.Context, msg *types.MsgSetFeedReward) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := s.SetFeedReward(ctx, msg)

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

func (s msgServer) FeedOwnershipTransferTx(c context.Context, msg *types.MsgFeedOwnershipTransfer) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := s.FeedOwnershipTransfer(ctx, msg)

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

func (s msgServer) RequestNewRoundTx(c context.Context, msg *types.MsgRequestNewRound) (*types.MsgResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	height, txHash, err := s.RequestNewRound(ctx, msg)

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
