package chainlink

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/keeper"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func handlerMsgCreateFeed(ctx sdk.Context, k keeper.Keeper, feed *types.MsgFeed) (*sdk.Result, error) {
	k.CreateFeed(ctx, *feed)

	return &sdk.Result{
		Data:   nil,
		Log:    "",
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}
