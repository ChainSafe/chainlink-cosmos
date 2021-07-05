package keeper

import (
	"fmt"
	"testing"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

// nolint
func setupKeeper(t testing.TB) (*Keeper, sdk.Context) {
	feedDataStoreKey := sdk.NewKVStoreKey(types.FeedDataStoreKey)
	roundStoreKey := sdk.NewKVStoreKey(types.RoundStoreKey)
	moduleOwnerStoreKey := sdk.NewKVStoreKey(types.ModuleOwnerStoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(feedDataStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(roundStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(moduleOwnerStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	keeper := NewKeeper(codec.NewProtoCodec(registry), feedDataStoreKey, roundStoreKey, moduleOwnerStoreKey, memStoreKey)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return keeper, ctx
}

// func setupKeeper(t testing.TB) (*Keeper, sdk.Context) {
// 	storeKey := sdk.NewKVStoreKey(types.StoreKey)
// 	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

// 	db := tmdb.NewMemDB()
// 	stateStore := store.NewCommitMultiStore(db)
// 	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
// 	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
// 	require.NoError(t, stateStore.LoadLatestVersion())

// 	registry := codectypes.NewInterfaceRegistry()
// 	keeper := NewKeeper(codec.NewProtoCodec(registry), storeKey, memStoreKey)

// 	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
// 	return keeper, ctx
// }

func TestFeedKeyStructure(t *testing.T) {
	k, ctx := setupKeeper(t)
	roundStore := ctx.KVStore(k.roundStoreKey)
	feedStore := ctx.KVStore(k.feedDataStoreKey)

	var tests = []struct {
		feedId  string
		roundId uint64
	}{
		{"test1", 1111},
		{"test11", 111},
		{"test111", 11},
		{"test1111", 1},
	}

	// Add all feed cases to store
	for _, tt := range tests {
		// force set roundId-1 for SetFeedData
		roundStore.Set(types.KeyPrefix(types.RoundIdKey+"/"+tt.feedId), i64tob(tt.roundId-1))

		feedData := types.MsgFeedData{
			FeedId: tt.feedId,
		}

		k.SetFeedData(ctx, &feedData)
	}

	// Retrieve key
	for _, tt := range tests {
		testName := fmt.Sprintf("%s,%d", tt.feedId, tt.roundId)
		t.Run(testName, func(t *testing.T) {
			prefixKey := types.FeedDataKey + "/" + tt.feedId + "/"
			//fmt.Println("[DEBUG] search for key", prefixKey)

			iterator := sdk.KVStorePrefixIterator(feedStore, types.KeyPrefix(prefixKey))

			defer iterator.Close()

			for ; iterator.Valid(); iterator.Next() {
				var feedData types.OCRFeedDataInStore
				k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &feedData)
				//fmt.Println("[DEBUG] found key", string(iterator.Key()), feedData.FeedData.FeedId, feedData.RoundId)

				if feedData.FeedData.FeedId != tt.feedId {
					t.Errorf("FeedId: got %s, want %s", feedData.FeedData.FeedId, tt.feedId)
				}

				if feedData.RoundId != tt.roundId {
					t.Errorf("RoundId: got %d, want %d", feedData.RoundId, tt.roundId)
				}
			}
		})
	}
}
