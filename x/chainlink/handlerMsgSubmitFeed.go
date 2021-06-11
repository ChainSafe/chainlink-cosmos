package chainlink

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/keeper"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func handlerMsgSubmitFeedData(ctx sdk.Context, k keeper.Keeper, feedData *types.MsgFeedData) (*sdk.Result, error) {
	k.SubmitFeedData(ctx, *feedData)

	return &sdk.Result{
		Data:   nil,
		Log:    "",
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}
