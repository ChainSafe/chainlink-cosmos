package ante

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func feedRewardSchemaStrategyChecker(strategy string) error {
	if strategy != "" {
		_, ok := types.FeedRewardStrategyConvertor[strategy]
		if !ok {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "invalid feed reward strategy")
		}
	}

	return nil
}
