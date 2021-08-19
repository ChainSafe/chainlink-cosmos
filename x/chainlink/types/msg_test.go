// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

// TODO: replace this method and import the one from util_test.go after merged.
func GenerateAccount() (types.PrivKey, string, sdk.AccAddress) {
	priv, pub, addr := testdata.KeyTestPubAddr()
	cosmosPubKey, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pub)
	if err != nil {
		panic(err)
	}
	return priv, cosmosPubKey, addr
}

type MsgFeedDataTestSuite struct {
	suite.Suite
	submitter  sdk.AccAddress
	feedId     string
	feedData   []byte
	signatures [][]byte
}

func TestMsgFeedDataTestSuite(t *testing.T) {
	suite.Run(t, new(MsgFeedDataTestSuite))
}

func (ts *MsgFeedDataTestSuite) SetupTest() {
	_, _, ts.submitter = GenerateAccount()
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
		{
			description: "MsgFeedDataTestSuite: passing case - all valid values",
			submitter:   ts.submitter,
			feedId:      ts.feedId,
			feedData:    ts.feedData,
			signatures:  ts.signatures,
			expPass:     true,
		},
		{
			description: "MsgFeedDataTestSuite: failing case - empty submitter",
			submitter:   nil,
			feedId:      ts.feedId,
			feedData:    ts.feedData,
			signatures:  ts.signatures,
			expPass:     false,
		},
		{
			description: "MsgFeedDataTestSuite: failing case - empty feedId",
			submitter:   ts.submitter,
			feedId:      "",
			feedData:    ts.feedData,
			signatures:  ts.signatures,
			expPass:     false,
		},
		{
			description: "MsgFeedDataTestSuite: failing case - invalid feedId format",
			submitter:   ts.submitter,
			feedId:      "BAD/FEED/ID",
			feedData:    ts.feedData,
			signatures:  ts.signatures,
			expPass:     false,
		},
		{
			description: "MsgFeedDataTestSuite: failing case - empty feedData",
			submitter:   ts.submitter,
			feedId:      ts.feedId,
			feedData:    nil,
			signatures:  ts.signatures,
			expPass:     false,
		},
		{
			description: "MsgFeedDataTestSuite: failing case - empty signatures",
			submitter:   ts.submitter,
			feedId:      ts.feedId,
			feedData:    ts.feedData,
			signatures:  [][]byte{},
			expPass:     false,
		},
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

type MsgModuleOwnerTestSuite struct {
	suite.Suite
	assignerAddress          sdk.AccAddress
	assignerPublicKey        []byte
	newModuleOwnerAddress    sdk.AccAddress
	newModuleOwnerPublicKey  []byte
	invalidModOwnerAddress   sdk.AccAddress
	invalidModOwnerPublicKey []byte
}

func TestMsgModuleOwnerTestSuite(t *testing.T) {
	suite.Run(t, new(MsgModuleOwnerTestSuite))
}

func (ts *MsgModuleOwnerTestSuite) SetupTest() {
	// assigner is a different account than the address + publicKey
	_, pubkey, addr := GenerateAccount()
	ts.assignerAddress = addr
	ts.assignerPublicKey = []byte(pubkey)

	_, pubkey, addr = GenerateAccount()
	ts.newModuleOwnerAddress = addr
	ts.newModuleOwnerPublicKey = []byte(pubkey)

	_, pubkey, addr = GenerateAccount()
	ts.invalidModOwnerAddress = addr
	ts.invalidModOwnerPublicKey = []byte(pubkey)
}

func (ts *MsgModuleOwnerTestSuite) TestMsgModuleOwnerConstructor() {
	msg := NewMsgModuleOwner(
		ts.assignerAddress,
		ts.newModuleOwnerAddress,
		ts.newModuleOwnerPublicKey,
	)

	bz := ModuleCdc.MustMarshalJSON(msg)
	signedBytes := sdk.MustSortJSON(bz)

	ts.Require().Equal(msg.Route(), RouterKey)
	ts.Require().Equal(msg.Type(), AddModuleOwner)
	ts.Require().Equal(msg.GetSigners(), []sdk.AccAddress{ts.assignerAddress})
	ts.Require().Equal(msg.GetSignBytes(), signedBytes)
	ts.Require().Equal(msg.GetMsgs(), []sdk.Msg{msg})
}

func (ts *MsgModuleOwnerTestSuite) TestMsgModuleOwnerValidateBasic() {
	testCases := []struct {
		description string
		assigner    sdk.AccAddress
		address     sdk.AccAddress
		publicKey   []byte
		expPass     bool
	}{
		{
			description: "MsgModuleOwnerTestSuite: passing case - all valid values",
			assigner:    ts.assignerAddress,
			address:     ts.newModuleOwnerAddress,
			publicKey:   []byte(ts.newModuleOwnerPublicKey),
			expPass:     true,
		},
		{
			description: "MsgModuleOwnerTestSuite: failing case - address and publicKey does not match",
			assigner:    ts.assignerAddress,
			address:     ts.newModuleOwnerAddress,
			publicKey:   ts.invalidModOwnerPublicKey,
			expPass:     false,
		},
	}

	for i, tc := range testCases {
		msg := NewMsgModuleOwner(
			tc.assigner,
			tc.address,
			tc.publicKey,
		)
		err := msg.ValidateBasic()

		if tc.expPass {
			ts.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.description)
		} else {
			ts.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.description)
		}
	}
}

type MsgModuleOwnershipTransferTestSuite struct {
	suite.Suite
	assignerAddress             sdk.AccAddress
	assignerPublicKey           []byte
	newModuleOwnerAddress       sdk.AccAddress
	newModuleOwnerPublicKey     []byte
	invalidModuleOwnerAddress   sdk.AccAddress
	invalidModuleOwnerPublicKey []byte
}

func TestMsgModuleOwnershipTransferTestSuite(t *testing.T) {
	suite.Run(t, new(MsgModuleOwnershipTransferTestSuite))
}

func (ts *MsgModuleOwnershipTransferTestSuite) SetupTest() {
	// assigner is a different account than the address + publicKey
	_, pubkey, addr := GenerateAccount()
	ts.assignerAddress = addr
	ts.assignerPublicKey = []byte(pubkey)

	_, pubkey, addr = GenerateAccount()
	ts.newModuleOwnerAddress = addr
	ts.newModuleOwnerPublicKey = []byte(pubkey)

	_, pubkey, addr = GenerateAccount()
	ts.invalidModuleOwnerAddress = addr
	ts.invalidModuleOwnerPublicKey = []byte(pubkey)
}

func (ts *MsgModuleOwnershipTransferTestSuite) TestMsgModuleOwnershipTransferConstructor() {
	msg := NewMsgModuleOwnershipTransfer(
		ts.assignerAddress,
		ts.newModuleOwnerAddress,
		ts.newModuleOwnerPublicKey,
	)

	bz := ModuleCdc.MustMarshalJSON(msg)
	signedBytes := sdk.MustSortJSON(bz)

	ts.Require().Equal(msg.Route(), RouterKey)
	ts.Require().Equal(msg.Type(), ModuleOwnershipTransfer)
	ts.Require().Equal(msg.GetSigners(), []sdk.AccAddress{ts.assignerAddress})
	ts.Require().Equal(msg.GetSignBytes(), signedBytes)
}

func (ts *MsgModuleOwnershipTransferTestSuite) TestMsgModuleOwnershipTransferValidateBasic() {
	testCases := []struct {
		description string
		assigner    sdk.AccAddress
		address     sdk.AccAddress
		publicKey   []byte
		expPass     bool
	}{
		{
			description: "MsgModuleOwnershipTransferTestSuite: passing case - all valid values",
			assigner:    ts.assignerAddress,
			address:     ts.newModuleOwnerAddress,
			publicKey:   []byte(ts.newModuleOwnerPublicKey),
			expPass:     true,
		},
		{
			description: "MsgModuleOwnershipTransferTestSuite: failing case - assigner address is empty",
			assigner:    nil,
			address:     ts.newModuleOwnerAddress,
			publicKey:   ts.invalidModuleOwnerPublicKey,
			expPass:     false,
		},
		{
			description: "MsgModuleOwnershipTransferTestSuite: failing case - address and publicKey does not match",
			assigner:    ts.assignerAddress,
			address:     ts.newModuleOwnerAddress,
			publicKey:   ts.invalidModuleOwnerPublicKey,
			expPass:     false,
		},
	}

	for i, tc := range testCases {
		msg := NewMsgModuleOwner(
			tc.assigner,
			tc.address,
			tc.publicKey,
		)
		err := msg.ValidateBasic()

		if tc.expPass {
			ts.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.description)
		} else {
			ts.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.description)
		}
	}
}

type MsgFeedTestSuite struct {
	suite.Suite
	feedOwner     sdk.AccAddress
	moduleOwner   sdk.AccAddress
	dataProviders []*DataProvider
}

func TestMsgFeedTestSuite(t *testing.T) {
	suite.Run(t, new(MsgFeedTestSuite))
}

func (ts *MsgFeedTestSuite) SetupTest() {
	_, _, feedOwnerAddr := GenerateAccount()
	ts.feedOwner = feedOwnerAddr

	_, _, moduleOwnerAddr := GenerateAccount()
	ts.moduleOwner = moduleOwnerAddr

	_, dpPubkey, dpAddr := GenerateAccount()
	dp := &DataProvider{
		Address: dpAddr,
		PubKey:  []byte(dpPubkey),
	}

	var dps []*DataProvider
	dps = append(dps, dp)
	ts.dataProviders = dps
}

func (ts *MsgFeedTestSuite) MsgFeedConstructor() {
	feedId := "feedId1"
	desc := "feedDesc1"
	submissionCount := uint32(1)
	heartbeatTrigger := uint32(2)
	deviationThresholdTrigger := uint32(3)
	feedReward := uint32(4)

	msg := NewMsgFeed(
		feedId,
		desc,
		ts.feedOwner,
		ts.moduleOwner,
		ts.dataProviders,
		submissionCount,
		heartbeatTrigger,
		deviationThresholdTrigger,
		feedReward,
	)

	bz := ModuleCdc.MustMarshalJSON(msg)
	signedBytes := sdk.MustSortJSON(bz)

	ts.Require().Equal(msg.Route(), RouterKey)
	ts.Require().Equal(msg.Type(), AddFeed)
	ts.Require().Equal(msg.GetSigners(), []sdk.AccAddress{ts.moduleOwner})
	ts.Require().Equal(msg.GetSignBytes(), signedBytes)
}

func (ts *MsgFeedTestSuite) TestMsgFeedValidateBasic() {
	testCases := []struct {
		description               string
		feedId                    string
		desc                      string
		feedOwner                 sdk.AccAddress
		moduleOwner               sdk.AccAddress
		dataProviders             []*DataProvider
		submissionCount           uint32
		heartbeatTrigger          uint32
		deviationThresholdTrigger uint32
		feedReward                uint32
		expPass                   bool
	}{
		{
			description:               "MsgFeedTestSuite: passing case - all valid values",
			feedId:                    "feed1",
			desc:                      "feedDescription1",
			feedOwner:                 ts.feedOwner,
			moduleOwner:               ts.moduleOwner,
			dataProviders:             ts.dataProviders,
			submissionCount:           uint32(1),
			heartbeatTrigger:          uint32(2),
			deviationThresholdTrigger: uint32(3),
			feedReward:                uint32(4),
			expPass:                   true,
		},
		{
			description:               "MsgFeedTestSuite: failing case - empty feed owner",
			feedId:                    "feed1",
			desc:                      "feedDescription1",
			feedOwner:                 nil,
			moduleOwner:               ts.moduleOwner,
			dataProviders:             ts.dataProviders,
			submissionCount:           uint32(1),
			heartbeatTrigger:          uint32(2),
			deviationThresholdTrigger: uint32(3),
			feedReward:                uint32(4),
			expPass:                   false,
		},
		{
			description:               "MsgFeedTestSuite: failing case - empty module owner",
			feedId:                    "feed1",
			desc:                      "feedDescription1",
			feedOwner:                 ts.feedOwner,
			moduleOwner:               nil,
			dataProviders:             ts.dataProviders,
			submissionCount:           uint32(1),
			heartbeatTrigger:          uint32(2),
			deviationThresholdTrigger: uint32(3),
			feedReward:                uint32(4),
			expPass:                   false,
		},
		{
			description:               "MsgFeedTestSuite: failing case - empty data providers",
			feedId:                    "feed1",
			desc:                      "feedDescription1",
			feedOwner:                 ts.feedOwner,
			moduleOwner:               ts.moduleOwner,
			dataProviders:             nil,
			submissionCount:           uint32(1),
			heartbeatTrigger:          uint32(2),
			deviationThresholdTrigger: uint32(3),
			feedReward:                uint32(4),
			expPass:                   false,
		},
	}

	for i, tc := range testCases {
		msg := NewMsgFeed(
			tc.feedId,
			tc.desc,
			tc.feedOwner,
			tc.moduleOwner,
			tc.dataProviders,
			tc.submissionCount,
			tc.heartbeatTrigger,
			tc.deviationThresholdTrigger,
			tc.feedReward,
		)
		err := msg.ValidateBasic()

		if tc.expPass {
			ts.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.description)
		} else {
			ts.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.description)
		}
	}

}
