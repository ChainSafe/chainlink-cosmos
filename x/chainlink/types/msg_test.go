// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type MsgFeedDataTestSuite struct {
	suite.Suite
	submitter  sdk.AccAddress
	feedId     string
	feedData   []byte
	signatures [][]byte
}

// TODO: replace this method and import the one from util_test.go after merged.
func GenerateAccount() sdk.AccAddress {
	_, _, addr := testdata.KeyTestPubAddr()
	return addr
}

func TestMsgFeedDataTestSuite(t *testing.T) {
	suite.Run(t, new(MsgFeedDataTestSuite))
}

func (ts *MsgFeedDataTestSuite) SetupTest() {
	ts.submitter = GenerateAccount()
	ts.feedId = "testfeed"
	ts.feedData = []byte("feedData")
	ts.signatures = [][]byte{[]byte("signatures")}
}

func (ts *MsgFeedDataTestSuite) TestMsgFeedDataConstructor() {
	msg := NewMsgFeedData(
		ts.submitter,
		ts.feedId,
		ts.feedData,
		ts.signatures,
	)

	bz := ModuleCdc.MustMarshalJSON(msg)
	signedBytes := sdk.MustSortJSON(bz)

	ts.Require().Equal(msg.Route(), RouterKey)
	ts.Require().Equal(msg.Type(), SubmitFeedData)
	ts.Require().Equal(msg.GetSigners(), []sdk.AccAddress{ts.submitter})
	ts.Require().Equal(msg.GetSignBytes(), signedBytes)
}

func (ts *MsgFeedDataTestSuite) TestMsgFeedDataValidateBasic() {
	testCases := []struct {
		description string
		submitter   sdk.AccAddress
		feedId      string
		feedData    []byte
		signatures  [][]byte
		expPass     bool
	}{
		{description: "passing test - all valid values", submitter: ts.submitter, feedId: ts.feedId, feedData: ts.feedData, signatures: ts.signatures, expPass: true},
		{description: "failing test - empty submitter", submitter: nil, feedId: ts.feedId, feedData: ts.feedData, signatures: ts.signatures, expPass: false},
		{description: "failing test - empty feedId", submitter: ts.submitter, feedId: "", feedData: ts.feedData, signatures: ts.signatures, expPass: false},
		{description: "failing test - invalid feedId format", submitter: ts.submitter, feedId: "BAD/FEED/ID", feedData: ts.feedData, signatures: ts.signatures, expPass: false},
		{description: "failing test - empty feedData", submitter: ts.submitter, feedId: ts.feedId, feedData: nil, signatures: ts.signatures, expPass: false},
		{description: "failing test - empty signatures", submitter: ts.submitter, feedId: ts.feedId, feedData: ts.feedData, signatures: [][]byte{}, expPass: false},
	}

	for i, tc := range testCases {
		msg := NewMsgFeedData(
			tc.submitter,
			tc.feedId,
			tc.feedData,
			tc.signatures,
		)
		err := msg.ValidateBasic()

		if tc.expPass {
			ts.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.description)
		} else {
			ts.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.description)
		}
	}
}
