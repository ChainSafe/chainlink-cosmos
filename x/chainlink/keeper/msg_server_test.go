package keeper

import (
	"context"
	"testing"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// nolint
func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	// nolint
	keeper, ctx := setupKeeper(t)
	return NewMsgServerImpl(*keeper), sdk.WrapSDKContext(ctx)
}
