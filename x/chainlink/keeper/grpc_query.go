package keeper

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
)

var _ types.QueryServer = Keeper{}
