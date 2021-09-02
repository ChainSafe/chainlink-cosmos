// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package keeper

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/ocr/utils"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
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
	accountStoreKey := sdk.NewKVStoreKey(types.AccountStoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(feedDataStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(roundStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(moduleOwnerStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(feedInfoStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(accountStoreKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	// TODO: do i need to replace nil -> bankKeeper? not quite sure if that can be exposed from this level
	keeper := NewKeeper(codec.NewProtoCodec(registry), nil, feedDataStoreKey, roundStoreKey, moduleOwnerStoreKey, feedInfoStoreKey, accountStoreKey, memStoreKey)

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
				FeedData:  utils.MustGenerateFakeABIReport(roundId, []int64{100, 101}),
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
				FeedId:   tc.feedId,
				FeedData: utils.MustGenerateFakeABIReport(tc.roundId, []int64{100, 101}),
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
		obs       []int64
		submitter []byte
		insert    bool
	}{
		{feedId: "feed1", roundId: 100, obs: []int64{100, 101, 102}, submitter: []byte("addressMock1"), insert: true},
		{feedId: "feed1", roundId: 200, obs: []int64{103, 104, 105}, submitter: []byte("addressMock2"), insert: true},
		{feedId: "feed1", roundId: 300, obs: []int64{106, 107, 108}, submitter: []byte("addressMock3"), insert: false},
		{feedId: "feed2", roundId: 400, obs: []int64{109, 110, 111}, submitter: []byte("addressMock4"), insert: false},
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
			FeedData:  utils.MustGenerateFakeABIReport(tc.roundId, tc.obs),
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
				require.EqualValues(t, tc.roundId, roundData[0].GetFeedData().GetContext().GetRound())
				require.Equal(t, tc.feedId, roundData[0].GetFeedId())

				observations := roundData[0].GetFeedData().GetReport().GetAttributedObservations()
				for i := 0; i < len(tc.obs); i++ {
					require.EqualValues(t, tc.obs[i], observations[i].GetObservation().GetValue()[0])
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
		obs       []int64
		submitter []byte
		insert    bool
	}{
		{feedId: "feed1", roundId: 100, expected: 100, obs: []int64{100, 101, 102}, submitter: []byte("addressMock1"), insert: true},
		{feedId: "feed1", roundId: 200, expected: 200, obs: []int64{103, 104, 105}, submitter: []byte("addressMock2"), insert: true},
		{feedId: "feed1", roundId: 300, expected: 200, obs: []int64{106, 107, 108}, submitter: []byte("addressMock3"), insert: false},
		{feedId: "feed2", roundId: 400, expected: 000, obs: []int64{109, 110, 111}, submitter: []byte("addressMock4"), insert: false},
		{feedId: "feed3", roundId: 500, expected: 500, obs: []int64{112, 113, 114}, submitter: []byte("addressMock5"), insert: true},
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
					FeedData:  utils.MustGenerateFakeABIReport(tc.roundId, tc.obs),
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
				require.EqualValues(t, tc.expected, roundData[0].GetFeedData().GetContext().GetRound())
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
	k, ctx := setupKeeper(t)
	moduleStore := ctx.KVStore(k.moduleOwnerStoreKey)

	_, pubKey, acc := testdata.KeyTestPubAddr()
	cosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pubKey)
	require.NoError(t, err)

	// store module owner
	k.SetModuleOwner(ctx, &types.MsgModuleOwner{
		Address: acc,
		PubKey:  []byte(cosmosPubKey),
	})

	// try to retrieve module owner
	data := moduleStore.Get(types.GetModuleOwnerKey(acc.String()))
	var moduleOwner types.MsgModuleOwner
	err = k.cdc.UnmarshalBinaryBare(data, &moduleOwner)
	require.NoError(t, err)

	require.EqualValues(t, acc, moduleOwner.GetAddress())
	require.EqualValues(t, cosmosPubKey, moduleOwner.GetPubKey())
}

func TestKeeper_RemoveModuleOwner(t *testing.T) {
	k, ctx := setupKeeper(t)
	moduleStore := ctx.KVStore(k.moduleOwnerStoreKey)

	_, pubKey, acc := testdata.KeyTestPubAddr()
	cosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pubKey)
	require.NoError(t, err)

	// store module owner
	k.SetModuleOwner(ctx, &types.MsgModuleOwner{
		Address: acc,
		PubKey:  []byte(cosmosPubKey),
	})

	// delete module owner
	k.RemoveModuleOwner(ctx, &types.MsgModuleOwnershipTransfer{
		AssignerAddress: acc,
	})

	// check if module owner doesn't exist anymore
	data := moduleStore.Get(types.GetModuleOwnerKey(acc.String()))
	require.Equal(t, []byte(nil), data)
}

func TestKeeper_GetModuleOwnerList(t *testing.T) {
	k, ctx := setupKeeper(t)

	_, pubKey1, acc1 := testdata.KeyTestPubAddr()
	cosmosPubKey1, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pubKey1)
	require.NoError(t, err)
	owner1 := &types.MsgModuleOwner{
		Address: acc1,
		PubKey:  []byte(cosmosPubKey1),
	}

	_, pubKey2, acc2 := testdata.KeyTestPubAddr()
	cosmosPubKey2, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pubKey2)
	require.NoError(t, err)
	owner2 := &types.MsgModuleOwner{
		Address: acc2,
		PubKey:  []byte(cosmosPubKey2),
	}

	testCases := []struct {
		test        string
		moduleOwner *types.MsgModuleOwner
		expected    []*types.MsgModuleOwner
	}{
		{
			test:        "owner 1",
			moduleOwner: owner1,
			expected:    []*types.MsgModuleOwner{owner1},
		},
		{
			test:        "owner 2 and previous one",
			moduleOwner: owner2,
			expected:    []*types.MsgModuleOwner{owner1, owner2},
		},
	}

	// Set module owner and try retrieve it
	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			k.SetModuleOwner(ctx, tc.moduleOwner)

			moduleOwner := k.GetModuleOwnerList(ctx)
			require.Equal(t, len(tc.expected), len(moduleOwner.GetModuleOwner()))
		})
	}
}

func TestKeeper_SetAndGetFeed(t *testing.T) {
	k, ctx := setupKeeper(t)

	feedToInsert := types.MsgFeed{
		FeedId:    "feed1",
		FeedOwner: GenerateAccount(),
		DataProviders: types.DataProviders{
			{Address: GenerateAccount()},
			{Address: GenerateAccount()},
		},
		SubmissionCount:           1,
		HeartbeatTrigger:          2,
		DeviationThresholdTrigger: 3,
		FeedReward: &types.FeedRewardSchema{
			Amount:   4,
			Strategy: "none",
		},
		Desc:               "desc1",
		ModuleOwnerAddress: GenerateAccount(),
	}

	k.SetFeed(ctx, &feedToInsert)
	result := k.GetFeed(ctx, feedToInsert.GetFeedId())

	require.Equal(t, feedToInsert.GetFeedId(), result.GetFeed().GetFeedId())
	require.Equal(t, feedToInsert.GetFeedOwner(), result.GetFeed().GetFeedOwner())
	require.EqualValues(t, feedToInsert.GetDataProviders(), result.GetFeed().GetDataProviders())
	require.Equal(t, feedToInsert.GetSubmissionCount(), result.GetFeed().GetSubmissionCount())
	require.Equal(t, feedToInsert.GetHeartbeatTrigger(), result.GetFeed().GetHeartbeatTrigger())
	require.Equal(t, feedToInsert.GetDeviationThresholdTrigger(), result.GetFeed().GetDeviationThresholdTrigger())
	require.Equal(t, feedToInsert.GetFeedReward(), result.GetFeed().GetFeedReward())
	require.Equal(t, feedToInsert.GetDesc(), result.GetFeed().GetDesc())
	require.Equal(t, feedToInsert.GetModuleOwnerAddress(), result.GetFeed().GetModuleOwnerAddress())
}

func TestKeeper_AddDataProvider(t *testing.T) {
	k, ctx := setupKeeper(t)

	dataProvider1 := GenerateAccount()
	dataProvider2 := GenerateAccount()

	feedToInsert := types.MsgFeed{
		FeedId: "feed1",
		DataProviders: types.DataProviders{
			{Address: dataProvider1},
		},
	}

	// store feed with 1 data provider
	k.SetFeed(ctx, &feedToInsert)

	// check if only 1 data provider is set
	result := k.GetFeed(ctx, feedToInsert.GetFeedId())
	require.Equal(t, feedToInsert.GetFeedId(), result.GetFeed().GetFeedId())
	require.EqualValues(t, 1, len(result.GetFeed().GetDataProviders()))
	require.EqualValues(t, feedToInsert.GetDataProviders(), result.GetFeed().GetDataProviders())

	// add new data provider
	_, _, err := k.AddDataProvider(ctx, &types.MsgAddDataProvider{
		FeedId: feedToInsert.GetFeedId(),
		DataProvider: &types.DataProvider{
			Address: dataProvider2,
		},
	})
	require.NoError(t, err)

	// check if 2 data provider are present
	result = k.GetFeed(ctx, feedToInsert.GetFeedId())
	require.Equal(t, feedToInsert.GetFeedId(), result.GetFeed().GetFeedId())
	require.EqualValues(t, 2, len(result.GetFeed().GetDataProviders()))
	require.EqualValues(t, types.DataProviders{
		{Address: dataProvider1},
		{Address: dataProvider2},
	}, result.GetFeed().GetDataProviders())
}

func TestKeeper_RemoveDataProvider(t *testing.T) {
	k, ctx := setupKeeper(t)

	dataProvider1 := GenerateAccount()
	dataProvider2 := GenerateAccount()

	feedToInsert := types.MsgFeed{
		FeedId: "feed1",
		DataProviders: types.DataProviders{
			{Address: dataProvider1},
			{Address: dataProvider2},
		},
	}

	// store feed with 2 data provider
	k.SetFeed(ctx, &feedToInsert)

	// check if 2 data provider are set
	result := k.GetFeed(ctx, feedToInsert.GetFeedId())
	require.Equal(t, feedToInsert.GetFeedId(), result.GetFeed().GetFeedId())
	require.EqualValues(t, 2, len(result.GetFeed().GetDataProviders()))
	require.EqualValues(t, feedToInsert.GetDataProviders(), result.GetFeed().GetDataProviders())

	// remove data provider #1
	_, _, err := k.RemoveDataProvider(ctx, &types.MsgRemoveDataProvider{
		FeedId:  feedToInsert.GetFeedId(),
		Address: dataProvider1,
	})
	require.NoError(t, err)

	// check if only data provider #2 is set
	result = k.GetFeed(ctx, feedToInsert.GetFeedId())
	require.Equal(t, feedToInsert.GetFeedId(), result.GetFeed().GetFeedId())
	require.EqualValues(t, 1, len(result.GetFeed().GetDataProviders()))
	require.EqualValues(t, types.DataProviders{
		{Address: dataProvider2},
	}, result.GetFeed().GetDataProviders())
}

func TestKeeper_ModifyFeedInfo(t *testing.T) {
	k, ctx := setupKeeper(t)

	type param struct {
		submissionCount           uint32
		heartbeatTrigger          uint32
		deviationThresholdTrigger uint32
		feedReward                *types.FeedRewardSchema
	}

	testCases := []struct {
		feedId string
		insert param
		modify param
	}{
		{
			feedId: "feed1",
			insert: param{submissionCount: 10, heartbeatTrigger: 20, deviationThresholdTrigger: 30, feedReward: &types.FeedRewardSchema{Amount: 4, Strategy: "none"}},
			modify: param{submissionCount: 11, heartbeatTrigger: 22, deviationThresholdTrigger: 33, feedReward: &types.FeedRewardSchema{Amount: 44, Strategy: "abc"}},
		},
		{
			feedId: "feed2",
			insert: param{submissionCount: 100, heartbeatTrigger: 200, deviationThresholdTrigger: 300, feedReward: &types.FeedRewardSchema{Amount: 400, Strategy: "none"}},
			modify: param{submissionCount: 101, heartbeatTrigger: 202, deviationThresholdTrigger: 303, feedReward: &types.FeedRewardSchema{Amount: 404, Strategy: "xyz"}},
		},
	}

	for _, tc := range testCases {
		testName := fmt.Sprintf("feed:%s", tc.feedId)
		t.Run(testName, func(t *testing.T) {
			k.SetFeed(ctx, &types.MsgFeed{
				FeedId:                    tc.feedId,
				SubmissionCount:           tc.insert.submissionCount,
				HeartbeatTrigger:          tc.insert.heartbeatTrigger,
				DeviationThresholdTrigger: tc.insert.deviationThresholdTrigger,
				FeedReward:                tc.insert.feedReward,
			})

			_, _, err := k.SetSubmissionCount(ctx, &types.MsgSetSubmissionCount{
				FeedId:          tc.feedId,
				SubmissionCount: tc.modify.submissionCount,
			})
			require.NoError(t, err)

			_, _, err = k.SetHeartbeatTrigger(ctx, &types.MsgSetHeartbeatTrigger{
				FeedId:           tc.feedId,
				HeartbeatTrigger: tc.modify.heartbeatTrigger,
			})
			require.NoError(t, err)

			_, _, err = k.SetDeviationThresholdTrigger(ctx, &types.MsgSetDeviationThresholdTrigger{
				FeedId:                    tc.feedId,
				DeviationThresholdTrigger: tc.modify.deviationThresholdTrigger,
			})
			require.NoError(t, err)

			_, _, err = k.SetFeedReward(ctx, &types.MsgSetFeedReward{
				FeedId:     tc.feedId,
				FeedReward: tc.modify.feedReward,
			})
			require.NoError(t, err)

			result := k.GetFeed(ctx, tc.feedId)
			require.Equal(t, tc.feedId, result.GetFeed().GetFeedId())
			require.Equal(t, tc.modify.submissionCount, result.GetFeed().GetSubmissionCount())
			require.Equal(t, tc.modify.heartbeatTrigger, result.GetFeed().GetHeartbeatTrigger())
			require.Equal(t, tc.modify.deviationThresholdTrigger, result.GetFeed().GetDeviationThresholdTrigger())
			require.Equal(t, tc.modify.feedReward, result.GetFeed().GetFeedReward())
		})
	}
}

func TestKeeper_DistributeReward(t *testing.T) {
	t.Skip("TODO")
}

func TestKeeper_FeedOwnershipTransfer(t *testing.T) {
	k, ctx := setupKeeper(t)

	oldFeedOwner := GenerateAccount()
	newFeedOwner := GenerateAccount()

	feedToInsert := types.MsgFeed{
		FeedId:    "feed1",
		FeedOwner: oldFeedOwner,
	}

	// store feed with old owner
	k.SetFeed(ctx, &feedToInsert)

	// check if old owner is set
	result := k.GetFeed(ctx, feedToInsert.GetFeedId())
	require.Equal(t, feedToInsert.GetFeedId(), result.GetFeed().GetFeedId())
	require.Equal(t, oldFeedOwner, result.GetFeed().GetFeedOwner())

	// transfer ownership to new owner
	_, _, err := k.FeedOwnershipTransfer(ctx, &types.MsgFeedOwnershipTransfer{
		FeedId:              feedToInsert.GetFeedId(),
		NewFeedOwnerAddress: newFeedOwner,
	})
	require.NoError(t, err)

	// check if new owner is set
	result = k.GetFeed(ctx, feedToInsert.GetFeedId())
	require.Equal(t, feedToInsert.GetFeedId(), result.GetFeed().GetFeedId())
	require.Equal(t, newFeedOwner, result.GetFeed().GetFeedOwner())
}
