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
		case *types.ModuleOwner:
			return handlerMsgAddModuleOwner(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handlerMsgSubmitFeedData(ctx sdk.Context, k keeper.Keeper, feedData *types.MsgFeedData) (*sdk.Result, error) {
	msgResult, err := k.SubmitFeedData(sdk.WrapSDKContext(ctx), feedData)
	if err != nil {
		return nil, err
	}

	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func handlerMsgAddModuleOwner(ctx sdk.Context, k keeper.Keeper, moduleOwner *types.ModuleOwner) (*sdk.Result, error) {
	fmt.Println("KKKEEEPPPERRRRRRR!!!", moduleOwner)
	msgResult, err := k.AddModuleOwner(sdk.WrapSDKContext(ctx), moduleOwner)
	if err != nil {
		return nil, err
	}
	result, err := sdk.WrapServiceResult(ctx, msgResult, err)
	if err != nil {
		return nil, err
	}

	return result, nil
}
