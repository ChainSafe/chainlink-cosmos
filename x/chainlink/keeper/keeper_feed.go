package keeper

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) SubmitFeedData(ctx sdk.Context, feedData types.MsgFeedData) {
	// TODO: add more complex feed validation here such as verify against other modules

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FeedDataKey))
	f := k.cdc.MustMarshalBinaryBare(&feedData)
	store.Set(types.KeyPrefix(types.FeedDataKey), f)
}

func (k Keeper) GetFeedData(ctx sdk.Context) (feedData []types.MsgFeedData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FeedDataKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.FeedDataKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var feed types.MsgFeedData
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &feed)

		// TODO: filter by feedId
		feedData = append(feedData, feed)
	}

	return
}
