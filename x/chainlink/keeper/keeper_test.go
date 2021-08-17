// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package keeper

import (
	"fmt"
	"strconv"
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
	feedInfoStoreKey := sdk.NewKVStoreKey(types.FeedInfoStoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(feedDataStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(roundStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(moduleOwnerStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(feedInfoStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	// TODO: do i need to replace nil -> bankKeeper? not quite sure if that can be exposed from this level
	keeper := NewKeeper(codec.NewProtoCodec(registry), nil, feedDataStoreKey, roundStoreKey, moduleOwnerStoreKey, feedInfoStoreKey, memStoreKey)

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
			roundStore.Set(types.GetRoundIdKey(tc.feedId), i64tob(roundId-1))

			feedData := types.MsgFeedData{
				FeedId:    tc.feedId,
				Submitter: []byte(fmt.Sprintf("%s/%d", tc.feedId, roundId)),
			}

			_, _, err := k.SetFeedData(ctx, &feedData)
			require.NoError(t, err)
		}
	}

	// Retrieve key
	for _, tc := range testCases {
		testName := fmt.Sprintf("feed:%s,round:%v", tc.feedId, tc.roundIds)
		t.Run(testName, func(t *testing.T) {
			prefixKey := types.GetFeedDataKey(tc.feedId, "")
			//fmt.Println("[DEBUG] search for key", string(prefixKey))

			iterator := sdk.KVStorePrefixIterator(feedStore, prefixKey)

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

func TestKeeper_SetFeedData(t *testing.T) {
	k, ctx := setupKeeper(t)
	roundStore := ctx.KVStore(k.roundStoreKey)
	feedDateStore := ctx.KVStore(k.feedDataStoreKey)

	testCases := []struct {
		feedId  string
		roundId uint64
	}{
		{feedId: "feed1", roundId: 100},
		{feedId: "feed1", roundId: 200},
		{feedId: "feed2", roundId: 300},
		{feedId: "feed2", roundId: 400},
	}

	// Add all feed cases to store and try retrieve them
	for _, tc := range testCases {
		testName := fmt.Sprintf("feed:%s,round:%d", tc.feedId, tc.roundId)
		t.Run(testName, func(t *testing.T) {
			// force set roundId-1 for SetFeedData
			roundStore.Set(types.GetRoundIdKey(tc.feedId), i64tob(tc.roundId-1))

			msgFeedData := types.MsgFeedData{
				FeedId: tc.feedId,
			}

			_, _, err := k.SetFeedData(ctx, &msgFeedData)
			require.NoError(t, err)

			roundId := roundStore.Get(types.GetRoundIdKey(tc.feedId))
			require.Equal(t, i64tob(tc.roundId), roundId)

			var feedData types.OCRFeedDataInStore
			value := feedDateStore.Get(types.GetFeedDataKey(tc.feedId, strconv.FormatUint(tc.roundId, 10)))
			err = k.cdc.UnmarshalBinaryBare(value, &feedData)
			require.NoError(t, err)
			require.Equal(t, tc.feedId, feedData.GetFeedData().GetFeedId())
		})
	}
}

func TestKeeper_GetRoundFeedDataByFilter(t *testing.T) {
	k, ctx := setupKeeper(t)
	roundStore := ctx.KVStore(k.roundStoreKey)

	testCases := []struct {
		feedId    string
		roundId   uint64
		feedData  []byte
		submitter []byte
		insert    bool
	}{
		{feedId: "feed1", roundId: 100, feedData: []byte{'a', 'b', 'c'}, submitter: []byte("addressMock1"), insert: true},
		{feedId: "feed1", roundId: 200, feedData: []byte{'d', 'e', 'f'}, submitter: []byte("addressMock2"), insert: true},
		{feedId: "feed1", roundId: 300, feedData: []byte{'g', 'h', 'i'}, submitter: []byte("addressMock3"), insert: false},
		{feedId: "feed2", roundId: 400, feedData: []byte{'j', 'k', 'l'}, submitter: []byte("addressMock4"), insert: false},
	}

	// Add all feed cases to store
	for _, tc := range testCases {
		if !tc.insert {
			continue
		}
		// force set roundId-1 for SetFeedData
		roundStore.Set(types.GetRoundIdKey(tc.feedId), i64tob(tc.roundId-1))

		msgFeedData := types.MsgFeedData{
			FeedId:    tc.feedId,
			FeedData:  tc.feedData,
			Submitter: tc.submitter,
		}

		_, _, err := k.SetFeedData(ctx, &msgFeedData)
		require.NoError(t, err)
	}

	// Retrieve feed data
	for _, tc := range testCases {
		testName := fmt.Sprintf("feed:%s,round:%d,inserted:%t", tc.feedId, tc.roundId, tc.insert)
		t.Run(testName, func(t *testing.T) {
			resp, err := k.GetRoundFeedDataByFilter(ctx, &types.GetRoundDataRequest{
				FeedId:  tc.feedId,
				RoundId: tc.roundId,
			})

			require.NoError(t, err)

			roundData := resp.GetRoundData()

			if tc.insert {
				require.Equal(t, 1, len(roundData))
				require.Equal(t, strconv.FormatUint(tc.roundId, 10), string(roundData[0].GetFeedData().Context))
				require.Equal(t, tc.feedId, roundData[0].FeedId)
				require.Equal(t, tc.submitter, roundData[0].FeedData.Oracles)

				observations := roundData[0].GetFeedData().GetObservations()
				for i := 0; i < len(tc.feedData); i++ {
					require.Equal(t, tc.feedData[i], observations[i].Data[0])
				}
			} else {
				require.Equal(t, 0, len(roundData))
			}
		})
	}
}

func TestKeeper_GetLatestRoundFeedDataByFilter(t *testing.T) {
	k, ctx := setupKeeper(t)

	roundStore := ctx.KVStore(k.roundStoreKey)

	testCases := []struct {
		feedId    string
		roundId   uint64
		expected  uint64
		feedData  []byte
		submitter []byte
		insert    bool
	}{
		{feedId: "feed1", roundId: 100, expected: 100, feedData: []byte{'a', 'b', 'c'}, submitter: []byte("addressMock1"), insert: true},
		{feedId: "feed1", roundId: 200, expected: 200, feedData: []byte{'d', 'e', 'f'}, submitter: []byte("addressMock2"), insert: true},
		{feedId: "feed1", roundId: 300, expected: 200, feedData: []byte{'g', 'h', 'i'}, submitter: []byte("addressMock3"), insert: false},
		{feedId: "feed2", roundId: 400, expected: 000, feedData: []byte{'j', 'k', 'l'}, submitter: []byte("addressMock4"), insert: false},
		{feedId: "feed3", roundId: 500, expected: 500, feedData: []byte{'m', 'n', 'o'}, submitter: []byte("addressMock5"), insert: true},
	}

	// Add all feed cases to store and try retrieve the latest round
	for _, tc := range testCases {
		testName := fmt.Sprintf("feed:%s,round:%d,inserted:%t", tc.feedId, tc.roundId, tc.insert)
		t.Run(testName, func(t *testing.T) {
			if tc.insert {
				// force set roundId-1 for SetFeedData
				roundStore.Set(types.GetRoundIdKey(tc.feedId), i64tob(tc.roundId-1))

				msgFeedData := types.MsgFeedData{
					FeedId:    tc.feedId,
					FeedData:  tc.feedData,
					Submitter: tc.submitter,
				}

				_, _, err := k.SetFeedData(ctx, &msgFeedData)
				require.NoError(t, err)
			}

			resp, err := k.GetLatestRoundFeedDataByFilter(ctx, &types.GetLatestRoundDataRequest{
				FeedId: tc.feedId,
			})
			require.NoError(t, err)

			roundData := resp.GetRoundData()

			// if roundId is expected
			if tc.expected > 0 {
				require.Equal(t, 1, len(roundData))
				require.Equal(t, strconv.FormatUint(tc.expected, 10), string(roundData[0].GetFeedData().Context))
				require.Equal(t, tc.feedId, roundData[0].FeedId)
			} else {
				require.Equal(t, 0, len(roundData))
			}
		})
	}
}

func TestKeeper_GetLatestRoundId(t *testing.T) {
	k, ctx := setupKeeper(t)
	roundStore := ctx.KVStore(k.roundStoreKey)

	testCases := []struct {
		name    string
		feedId  string
		roundId uint64
		insert  bool
	}{
		{feedId: "feed1", roundId: 1, insert: true},
		{feedId: "feed1", roundId: 2, insert: true},
		{feedId: "feed2", roundId: 3, insert: true},
		{feedId: "feed2", roundId: 4, insert: true},
		{roundId: 4, insert: false},                        // get latest global roundId
		{feedId: "nonExisting", roundId: 0, insert: false}, // get non-existing roundId (should return 0)
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("feed:%s,round:%d", tc.feedId, tc.roundId)
		t.Run(testName, func(t *testing.T) {
			if tc.insert {
				roundStore.Set(types.GetRoundIdKey(tc.feedId), i64tob(tc.roundId))
			}

			latestRoundId := k.GetLatestRoundId(ctx, tc.feedId)
			require.Equal(t, tc.roundId, latestRoundId)
		})
	}
}

func TestKeeper_SetModuleOwner(t *testing.T) {
	t.Skip("TODO")
}

func TestKeeper_RemoveModuleOwner(t *testing.T) {
	t.Skip("TODO")
}

func TestKeeper_GetModuleOwnerList(t *testing.T) {
	t.Skip("TODO")
}
