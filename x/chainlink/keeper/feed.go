package keeper

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CreateFeed(ctx sdk.Context, feed types.MsgFeed) {
	// TODO: add more complex feed validation here such as verify against other modules

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FeedKey))
	f := k.cdc.MustMarshalBinaryBare(&feed)
	store.Set(types.KeyPrefix(types.FeedKey), f)
}

func (k Keeper) GetAllFeed(ctx sdk.Context) (msgs []types.MsgFeed) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.FeedKey))
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.FeedKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var msg types.MsgFeed
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &msg)
		msgs = append(msgs, msg)
	}

	return
}
