// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package keeper

import (
	"testing"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var (
	testfeedid     = "testfeed"
	testfeedData   = []byte("feedData")
	testRoundID    = uint64(310)
	testsignatures = [][]byte{[]byte("signatures")}
	testContext    = &types.ReportContext{}
	testReport     = &types.AttestedReportMany{}

	testNum      = uint64(310)
	testNumBytes = []uint8([]byte{0x36, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
)

func GenerateAccount() sdk.AccAddress {
	_, _, addr := testdata.KeyTestPubAddr()
	return addr
}

func TestFeedDataFilter(t *testing.T) {
	feedData := types.OCRFeedDataInStore{
		FeedData: &types.MsgFeedData{FeedId: testfeedid, Submitter: GenerateAccount(), FeedData: testfeedData, Signatures: testsignatures},
		RoundId:  testRoundID,
		Report:   &types.OffChainReport{Context: testContext, Report: testReport},
	}

	expRoundData := &types.RoundData{
		FeedId:   feedData.GetFeedData().GetFeedId(),
		FeedData: feedData.GetReport(),
	}

	require.Equal(t, feedDataFilter(testfeedid, testRoundID, feedData), expRoundData)
}

func TestI64tob(t *testing.T) {
	require.Equal(t, testNumBytes, i64tob(testNum))
}

func TestBtoi64(t *testing.T) {
	require.Equal(t, testNum, btoi64(testNumBytes))
}
