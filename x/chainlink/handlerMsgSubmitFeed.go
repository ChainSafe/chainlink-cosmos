package chainlink

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/keeper"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
