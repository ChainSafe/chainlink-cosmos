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
		case *types.MsgAddFeedProvider:
			return handlerMsgAddFeedProvider(ctx, k, msg)
		case *types.MsgRemoveFeedProvider:
			return handlerMsgRemoveFeedProvider(ctx, k, msg)
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

func handlerMsgAddFeedProvider(ctx sdk.Context, k keeper.Keeper, msgAddFeedProvider *types.MsgAddFeedProvider) (*sdk.Result, error) {
	msgResult, err := k.AddFeedProviderTx(sdk.WrapSDKContext(ctx), msgAddFeedProvider)
	if err != nil {
		return nil, err
	}
	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func handlerMsgRemoveFeedProvider(ctx sdk.Context, k keeper.Keeper, msgRemoveFeedProvider *types.MsgRemoveFeedProvider) (*sdk.Result, error) {
	msgResult, err := k.RemoveFeedProviderTx(sdk.WrapSDKContext(ctx), msgRemoveFeedProvider)
	if err != nil {
		return nil, err
	}
	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}
	return result, nil
}
