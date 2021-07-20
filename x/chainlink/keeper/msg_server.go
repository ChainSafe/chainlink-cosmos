package keeper

import (
	"context"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ErrIncorrectHeightFound = "incorrect height found"
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
	height, txHash := s.SetFeedData(ctx, msg)

	if height == 0 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidHeight, ErrIncorrectHeightFound)
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

	return &types.MsgResponse{
		Height: uint64(height),
		TxHash: string(txHash),
	}, nil
}
