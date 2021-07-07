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

func TestFeedKeyStructure(t *testing.T) {
	k, ctx := setupKeeper(t)
	roundStore := ctx.KVStore(k.roundStoreKey)
	feedStore := ctx.KVStore(k.feedDataStoreKey)

	testCases := []struct {
		feedId   string
		roundIds []uint64
	}{
		{feedId: "test1", roundIds: []uint64{1, 11, 111, 1111}},
		{feedId: "test11", roundIds: []uint64{1, 11, 111, 1111}},
		{feedId: "test111", roundIds: []uint64{1, 11, 111, 1111}},
		{feedId: "test1111", roundIds: []uint64{1, 11, 111, 1111}},
	}

	// Add all feed cases to store
	for _, tc := range testCases {
		for _, roundId := range tc.roundIds {
			// force set roundId-1 for SetFeedData
			roundStore.Set(types.KeyPrefix(types.RoundIdKey+"/"+tc.feedId), i64tob(roundId-1))

			feedData := types.MsgFeedData{
				FeedId:    tc.feedId,
				Submitter: []byte(fmt.Sprintf("%s/%d", tc.feedId, roundId)),
			}

			k.SetFeedData(ctx, &feedData)
		}
	}

	// Retrieve key
	for _, tc := range testCases {
		testName := fmt.Sprintf("%s,%v", tc.feedId, tc.roundIds)
		t.Run(testName, func(t *testing.T) {
			prefixKey := types.FeedDataKey + "/" + tc.feedId + "/"
			//fmt.Println("[DEBUG] search for key", prefixKey)

			iterator := sdk.KVStorePrefixIterator(feedStore, types.KeyPrefix(prefixKey))

			defer iterator.Close()

			for ; iterator.Valid(); iterator.Next() {
				var feedData types.OCRFeedDataInStore
				k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &feedData)
				//fmt.Println("[DEBUG] found key", string(iterator.Key()), feedData.FeedData.FeedId, feedData.RoundId)

				require.Equal(t, tc.feedId, feedData.GetFeedData().GetFeedId())
				require.Equal(t, []byte(fmt.Sprintf("%s/%d", tc.feedId, feedData.GetRoundId())), feedData.GetFeedData().GetSubmitter().Bytes())
				require.Contains(t, tc.roundIds, feedData.GetRoundId())
			}
		})
	}
}
