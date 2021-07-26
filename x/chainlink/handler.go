// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package chainlink

import (
	"fmt"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/keeper"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgFeedData:
			res, err := msgServer.SubmitFeedDataTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgModuleOwner:
			res, err := msgServer.AddModuleOwnerTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgModuleOwnershipTransfer:
			res, err := msgServer.ModuleOwnershipTransferTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgFeed:
			res, err := msgServer.AddFeedTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgAddDataProvider:
			res, err := msgServer.AddDataProviderTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRemoveDataProvider:
			res, err := msgServer.RemoveDataProviderTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSetSubmissionCount:
			res, err := msgServer.SetSubmissionCountTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSetHeartbeatTrigger:
			res, err := msgServer.SetHeartbeatTriggerTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSetDeviationThresholdTrigger:
			res, err := msgServer.SetDeviationThresholdTriggerTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSetFeedReward:
			res, err := msgServer.SetFeedRewardTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgFeedOwnershipTransfer:
			res, err := msgServer.FeedOwnershipTransferTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgRequestNewRound:
			res, err := msgServer.RequestNewRoundTx(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
