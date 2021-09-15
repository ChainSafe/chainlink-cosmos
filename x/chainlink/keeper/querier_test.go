// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package keeper

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQuerier_GetRoundFeedData(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	amino := codec.NewLegacyAmino()
	querier := NewQuerier(*keeper, amino)
	roundStore := ctx.KVStore(keeper.roundStoreKey)

	testCases := []struct {
		feedId          string
		roundId         uint64
		feedData        [][]byte
		submitter       sdk.AccAddress
		signature       [][]byte
		isFeedDataValid bool
		insert          bool
	}{
		{feedId: "feed1", roundId: 100, feedData: [][]byte{[]byte("a"), []byte("b"), []byte("c")}, submitter: sdk.AccAddress("addressMock1"), signature: [][]byte{{'a', 'b'}, {'c', 'd'}}, isFeedDataValid: true, insert: true},
		{feedId: "feed1", roundId: 200, feedData: [][]byte{[]byte("a"), []byte("b"), []byte("c")}, submitter: sdk.AccAddress("addressMock2"), signature: [][]byte{{'e', 'f'}, {'g', 'h'}}, isFeedDataValid: false, insert: true},
		{feedId: "feed1", roundId: 300, feedData: [][]byte{[]byte("a"), []byte("b"), []byte("c")}, submitter: sdk.AccAddress("addressMock3"), signature: [][]byte{{'i', 'j'}, {'k', 'l'}}, isFeedDataValid: true, insert: false},
		{feedId: "feed2", roundId: 400, feedData: [][]byte{[]byte("a"), []byte("b"), []byte("c")}, submitter: sdk.AccAddress("addressMock4"), signature: [][]byte{{'m', 'n'}, {'o', 'p'}}, isFeedDataValid: false, insert: false},
	}

	// Add all feed cases to store
	for _, tc := range testCases {
		if !tc.insert {
			continue
		}
		// force set roundId-1 for SetFeedData
		roundStore.Set(types.GetRoundIdKey(tc.feedId), i64tob(tc.roundId-1))

		msgFeedData := types.MsgFeedData{
			FeedId:                        tc.feedId,
			ObservationFeedData:           tc.feedData,
			Submitter:                     tc.submitter,
			ObservationFeedDataSignatures: tc.signature,
			IsFeedDataValid:               tc.isFeedDataValid,
		}

		_, _, err := keeper.SetFeedData(ctx, &msgFeedData)
		require.NoError(t, err)
	}

	// Add all feed cases to store and try retrieve them
	for _, tc := range testCases {
		testName := fmt.Sprintf("feed:%s,round:%d", tc.feedId, tc.roundId)
		t.Run(testName, func(t *testing.T) {
			result, err := querier(ctx, []string{types.QueryRoundFeedData, strconv.FormatUint(tc.roundId, 10), tc.feedId}, abci.RequestQuery{})
			require.NoError(t, err)

			var roundDataResponse types.GetRoundDataResponse
			err = amino.UnmarshalJSON(result, &roundDataResponse)
			require.NoError(t, err)

			if tc.insert {
				require.Equal(t, 1, len(roundDataResponse.GetRoundData()))
				require.Equal(t, tc.feedId, roundDataResponse.GetRoundData()[0].GetFeedId())
				require.Equal(t, strconv.FormatUint(tc.roundId, 10), string(roundDataResponse.GetRoundData()[0].GetFeedData().GetContext()))
				require.Equal(t, tc.submitter.Bytes(), roundDataResponse.GetRoundData()[0].GetFeedData().GetOracles())

				// TODO if tc.isFeedDataValid is true, check if event is emitted with correct signature

				observations := roundDataResponse.GetRoundData()[0].GetFeedData().GetObservations()
				for i := 0; i < len(tc.feedData); i++ {
					require.Equal(t, string(tc.feedData[i]), string(observations[i].Data[0]))
				}
			} else {
				require.Equal(t, 0, len(roundDataResponse.GetRoundData()))
			}
		})
	}
}

func TestQuerier_LatestRoundFeedData(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	amino := codec.NewLegacyAmino()
	querier := NewQuerier(*keeper, amino)
	roundStore := ctx.KVStore(keeper.roundStoreKey)

	testCases := []struct {
		feedId          string
		roundId         uint64
		expected        uint64
		feedData        [][]byte
		submitter       []byte
		signature       [][]byte
		isFeedDataValid bool
		insert          bool
	}{
		{feedId: "feed1", roundId: 100, expected: 100, feedData: [][]byte{[]byte("a"), []byte("b"), []byte("c")}, submitter: sdk.AccAddress("addressMock1"), signature: [][]byte{{'a', 'b'}, {'c', 'd'}}, isFeedDataValid: true, insert: true},
		{feedId: "feed1", roundId: 200, expected: 200, feedData: [][]byte{[]byte("a"), []byte("b"), []byte("c")}, submitter: sdk.AccAddress("addressMock2"), signature: [][]byte{{'e', 'f'}, {'g', 'h'}}, isFeedDataValid: false, insert: true},
		{feedId: "feed1", roundId: 300, expected: 200, feedData: [][]byte{[]byte("a"), []byte("b"), []byte("c")}, submitter: sdk.AccAddress("addressMock3"), signature: [][]byte{{'i', 'j'}, {'k', 'l'}}, isFeedDataValid: true, insert: false},
		{feedId: "feed2", roundId: 400, expected: 000, feedData: [][]byte{[]byte("a"), []byte("b"), []byte("c")}, submitter: sdk.AccAddress("addressMock4"), signature: [][]byte{{'m', 'n'}, {'o', 'p'}}, isFeedDataValid: false, insert: false},
		{feedId: "feed3", roundId: 500, expected: 500, feedData: [][]byte{[]byte("a"), []byte("b"), []byte("c")}, submitter: sdk.AccAddress("addressMock5"), signature: [][]byte{{'q', 'r'}, {'s', 't'}}, isFeedDataValid: true, insert: true},
	}

	// Add all feed cases to store and try retrieve the latest round
	for _, tc := range testCases {
		testName := fmt.Sprintf("feed:%s,round:%d,inserted:%t", tc.feedId, tc.roundId, tc.insert)
		t.Run(testName, func(t *testing.T) {
			if tc.insert {
				// force set roundId-1 for SetFeedData
				roundStore.Set(types.GetRoundIdKey(tc.feedId), i64tob(tc.roundId-1))

				msgFeedData := types.MsgFeedData{
					FeedId:                        tc.feedId,
					ObservationFeedData:           tc.feedData,
					Submitter:                     tc.submitter,
					ObservationFeedDataSignatures: tc.signature,
					IsFeedDataValid:               tc.isFeedDataValid,
				}

				_, _, err := keeper.SetFeedData(ctx, &msgFeedData)
				require.NoError(t, err)
			}

			result, err := querier(ctx, []string{types.QueryLatestFeedData, tc.feedId}, abci.RequestQuery{})
			require.NoError(t, err)

			var roundDataResponse types.GetRoundDataResponse
			err = amino.UnmarshalJSON(result, &roundDataResponse)
			require.NoError(t, err)

			// if roundId is expected
			if tc.expected > 0 {
				require.Equal(t, 1, len(roundDataResponse.GetRoundData()))
				require.Equal(t, tc.feedId, roundDataResponse.GetRoundData()[0].GetFeedId())
				require.Equal(t, strconv.FormatUint(tc.expected, 10), string(roundDataResponse.GetRoundData()[0].GetFeedData().GetContext()))

				// TODO if tc.isFeedDataValid is true, check if event is emitted with correct signature
			} else {
				require.Equal(t, 0, len(roundDataResponse.GetRoundData()))
			}
		})
	}
}

func TestQuerier_GetFeedInfo(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	amino := codec.NewLegacyAmino()
	querier := NewQuerier(*keeper, amino)

	testCases := []*types.MsgFeed{
		{
			FeedId:    "feed1",
			FeedOwner: GenerateAccount(),
			DataProviders: []*types.DataProvider{
				{Address: GenerateAccount()},
			},
			SubmissionCount:           1,
			HeartbeatTrigger:          2,
			DeviationThresholdTrigger: 3,
			FeedReward: &types.FeedRewardSchema{
				Amount:   4,
				Strategy: "none",
			},
			Desc:               "desc test",
			ModuleOwnerAddress: GenerateAccount(),
		},
		{
			FeedId:    "feed1",
			FeedOwner: GenerateAccount(),
			DataProviders: []*types.DataProvider{
				{Address: GenerateAccount()},
				{Address: GenerateAccount()},
			},
			SubmissionCount:           10,
			HeartbeatTrigger:          20,
			DeviationThresholdTrigger: 30,
			FeedReward: &types.FeedRewardSchema{
				Amount:   40,
				Strategy: "abc",
			},
			Desc:               "desc test 2",
			ModuleOwnerAddress: GenerateAccount(),
		},
		{
			FeedId:    "feed2",
			FeedOwner: GenerateAccount(),
			DataProviders: []*types.DataProvider{
				{Address: GenerateAccount()},
				{Address: GenerateAccount()},
			},
			SubmissionCount:           100,
			HeartbeatTrigger:          200,
			DeviationThresholdTrigger: 300,
			FeedReward: &types.FeedRewardSchema{
				Amount:   400,
				Strategy: "xyz",
			},
			Desc:               "desc test 3",
			ModuleOwnerAddress: GenerateAccount(),
		},
	}

	// Add feed to store and try retrieve it
	for _, tc := range testCases {
		testName := fmt.Sprintf("feed:%s,desc:%s", tc.FeedId, tc.Desc)
		t.Run(testName, func(t *testing.T) {
			keeper.SetFeed(ctx, tc)

			result, err := querier(ctx, []string{types.QueryFeedInfo, tc.FeedId}, abci.RequestQuery{})
			require.NoError(t, err)

			var feedInfo types.GetFeedByIdResponse
			err = amino.UnmarshalJSON(result, &feedInfo)
			require.NoError(t, err)

			require.Equal(t, tc.FeedId, feedInfo.GetFeed().GetFeedId())
			require.Equal(t, tc.FeedOwner, feedInfo.GetFeed().GetFeedOwner())
			require.Equal(t, tc.DataProviders, feedInfo.GetFeed().GetDataProviders())
			require.Equal(t, tc.SubmissionCount, feedInfo.GetFeed().GetSubmissionCount())
			require.Equal(t, tc.HeartbeatTrigger, feedInfo.GetFeed().GetHeartbeatTrigger())
			require.Equal(t, tc.DeviationThresholdTrigger, feedInfo.GetFeed().GetDeviationThresholdTrigger())
			require.Equal(t, tc.FeedReward, feedInfo.GetFeed().GetFeedReward())
			require.Equal(t, tc.Desc, feedInfo.GetFeed().GetDesc())
		})
	}
}

func TestQuerier_GetModuleOwners(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	amino := codec.NewLegacyAmino()
	querier := NewQuerier(*keeper, amino)

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
			keeper.SetModuleOwner(ctx, tc.moduleOwner)

			result, err := querier(ctx, []string{types.QueryModuleOwner}, abci.RequestQuery{})
			require.NoError(t, err)

			var moduleOwner types.GetModuleOwnerResponse
			err = amino.UnmarshalJSON(result, &moduleOwner)
			require.NoError(t, err)

			require.Equal(t, len(tc.expected), len(moduleOwner.GetModuleOwner()))
		})
	}
}

func TestQuerier_Fail(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	amino := codec.NewLegacyAmino()
	querier := NewQuerier(*keeper, amino)

	testCases := []struct {
		name           string
		path           []string
		expectedErr    error
		expectedResult []byte
	}{
		{
			name:           "QueryFeedInfo: missing feed id",
			path:           []string{types.QueryFeedInfo},
			expectedErr:    sdkerrors.ErrInvalidRequest,
			expectedResult: nil,
		},
		{
			name:           "QueryFeedInfo: unknown feed id",
			path:           []string{types.QueryFeedInfo, "unknownFeedId"},
			expectedErr:    sdkerrors.ErrKeyNotFound,
			expectedResult: nil,
		},
		{
			name:           "QueryRoundFeedData: missing round id and feed id",
			path:           []string{types.QueryRoundFeedData},
			expectedErr:    sdkerrors.ErrInvalidRequest,
			expectedResult: nil,
		},
		{
			name:           "QueryRoundFeedData: non-string round id",
			path:           []string{types.QueryRoundFeedData, "roundIdString", "feedId"},
			expectedErr:    strconv.ErrSyntax,
			expectedResult: nil,
		},
		{
			name:           "QueryRoundFeedData: unknown round id and feed id",
			path:           []string{types.QueryRoundFeedData, "999", "unknownFeedId"},
			expectedErr:    nil, // return no error, just empty pagination result
			expectedResult: []byte("{\n  \"pagination\": {}\n}"),
		},
		{
			name:           "QueryLatestFeedData: missing feed id",
			path:           []string{types.QueryLatestFeedData},
			expectedErr:    sdkerrors.ErrInvalidRequest,
			expectedResult: nil,
		},
		{
			name:           "QueryLatestFeedData: unknown feed id",
			path:           []string{types.QueryLatestFeedData, "unknownFeedId"},
			expectedErr:    nil, // return no error, just empty result
			expectedResult: []byte("{}"),
		},
		{
			name:           "QueryModuleOwner: no module owner set",
			path:           []string{types.QueryModuleOwner},
			expectedErr:    nil,
			expectedResult: []byte("{}"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := querier(ctx, tc.path, abci.RequestQuery{})
			require.ErrorIs(t, err, tc.expectedErr)
			require.Equal(t, tc.expectedResult, result)
		})
	}
}
