package keeper

import (
	"github.com/ChainSafe/chainlink/x/chainlink/types"
)

var _ types.QueryServer = Keeper{}
