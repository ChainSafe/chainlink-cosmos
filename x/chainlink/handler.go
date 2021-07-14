package chainlink

import (
	"fmt"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/keeper"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgFeedData:
			return handlerMsgSubmitFeedData(ctx, k, msg)
		case *types.MsgModuleOwner:
			return handlerMsgAddModuleOwner(ctx, k, msg)
		case *types.MsgModuleOwnershipTransfer:
			return handlerMsgModuleOwnershipTransfer(ctx, k, msg)
		case *types.MsgFeed:
			return handlerMsgAddNewFeed(ctx, k, msg)
		case *types.MsgAddDataProvider:
			return handlerMsgAddDataProvider(ctx, k, msg)
		case *types.MsgRemoveDataProvider:
			return handlerMsgRemoveDataProvider(ctx, k, msg)
		case *types.MsgSetSubmissionCount:
			return handlerMsgSetSubmissionCount(ctx, k, msg)
		case *types.MsgSetHeartbeatTrigger:
			return handlerMsgSetHeartbeatTrigger(ctx, k, msg)
		case *types.MsgSetDeviationThresholdTrigger:
			return handlerMsgSetDeviationThresholdTrigger(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handlerMsgSubmitFeedData(ctx sdk.Context, k keeper.Keeper, feedData *types.MsgFeedData) (*sdk.Result, error) {
	msgResult, err := k.SubmitFeedDataTx(sdk.WrapSDKContext(ctx), feedData)
	if err != nil {
		return nil, err
	}
	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func handlerMsgAddModuleOwner(ctx sdk.Context, k keeper.Keeper, moduleOwner *types.MsgModuleOwner) (*sdk.Result, error) {
	msgResult, err := k.AddModuleOwnerTx(sdk.WrapSDKContext(ctx), moduleOwner)
	if err != nil {
		return nil, err
	}
	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func handlerMsgModuleOwnershipTransfer(ctx sdk.Context, k keeper.Keeper, moduleOwner *types.MsgModuleOwnershipTransfer) (*sdk.Result, error) {
	msgResult, err := k.ModuleOwnershipTransferTx(sdk.WrapSDKContext(ctx), moduleOwner)
	if err != nil {
		return nil, err
	}
	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func handlerMsgAddNewFeed(ctx sdk.Context, k keeper.Keeper, newFeed *types.MsgFeed) (*sdk.Result, error) {
	msgResult, err := k.AddFeedTx(sdk.WrapSDKContext(ctx), newFeed)
	if err != nil {
		return nil, err
	}
	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func handlerMsgAddDataProvider(ctx sdk.Context, k keeper.Keeper, msgAddDataProvider *types.MsgAddDataProvider) (*sdk.Result, error) {
	msgResult, err := k.AddDataProviderTx(sdk.WrapSDKContext(ctx), msgAddDataProvider)
	if err != nil {
		return nil, err
	}
	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func handlerMsgRemoveDataProvider(ctx sdk.Context, k keeper.Keeper, msgRemoveDataProvider *types.MsgRemoveDataProvider) (*sdk.Result, error) {
	msgResult, err := k.RemoveDataProviderTx(sdk.WrapSDKContext(ctx), msgRemoveDataProvider)
	if err != nil {
		return nil, err
	}
	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func handlerMsgSetSubmissionCount(ctx sdk.Context, k keeper.Keeper, msgSetSubmissionCount *types.MsgSetSubmissionCount) (*sdk.Result, error) {
	msgResult, err := k.SetSubmissionCountTx(sdk.WrapSDKContext(ctx), msgSetSubmissionCount)
	if err != nil {
		return nil, err
	}
	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func handlerMsgSetHeartbeatTrigger(ctx sdk.Context, k keeper.Keeper, msgSetHeartbeatTrigger *types.MsgSetHeartbeatTrigger) (*sdk.Result, error) {
	msgResult, err := k.SetHeartbeatTriggerTx(sdk.WrapSDKContext(ctx), msgSetHeartbeatTrigger)
	if err != nil {
		return nil, err
	}
	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func handlerMsgSetDeviationThresholdTrigger(ctx sdk.Context, k keeper.Keeper, msgSetDeviationThresholdTrigger *types.MsgSetDeviationThresholdTrigger) (*sdk.Result, error) {
	msgResult, err := k.SetDeviationThresholdTriggerTx(sdk.WrapSDKContext(ctx), msgSetDeviationThresholdTrigger)
	if err != nil {
		return nil, err
	}
	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}
