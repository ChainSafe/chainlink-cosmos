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
			feedId:                    "feedId1",
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
			feedId:                    "feedId1",
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
			feedId:                    "feedId1",
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
			feedId:                    "feedId1",
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

type MsgAddDataProviderTestSuite struct {
	suite.Suite
	signer              sdk.AccAddress
	validDataProvider   *DataProvider
	invalidDataProvider *DataProvider
}

func TestMsgAddDataProviderTestSuite(t *testing.T) {
	suite.Run(t, new(MsgAddDataProviderTestSuite))
}

func (ts *MsgAddDataProviderTestSuite) SetupTest() {
	_, dpPubkey, dpAddr := GenerateAccount()
	vdp := &DataProvider{
		Address: dpAddr,
		PubKey:  []byte(dpPubkey),
	}

	ts.validDataProvider = vdp

	_, dpPubkey, _ = GenerateAccount()
	_, _, dpAddr = GenerateAccount()

	idp := &DataProvider{
		Address: dpAddr,
		PubKey:  []byte(dpPubkey),
	}

	ts.invalidDataProvider = idp

	_, _, signerAddr := GenerateAccount()
	ts.signer = signerAddr
}

func (ts *MsgAddDataProviderTestSuite) MsgAddDataProviderConstructor() {
	msg := NewMsgAddDataProvider(
		ts.signer,
		"feedId1",
		ts.validDataProvider,
	)

	bz := ModuleCdc.MustMarshalJSON(msg)
	signedBytes := sdk.MustSortJSON(bz)

	ts.Require().Equal(msg.Route(), RouterKey)
	ts.Require().Equal(msg.Type(), AddDataProvider)
	ts.Require().Equal(msg.GetSigners(), []sdk.AccAddress{ts.signer})
	ts.Require().Equal(msg.GetSignBytes(), signedBytes)
}

func (ts *MsgAddDataProviderTestSuite) TestMsgAddDataProviderValidateBasic() {
	testCases := []struct {
		description  string
		feedId       string
		signer       sdk.AccAddress
		dataProvider *DataProvider
		expPass      bool
	}{
		{
			description:  "MsgAddDataProviderTestSuite: passing case - all valid values",
			feedId:       "feedId1",
			signer:       ts.signer,
			dataProvider: ts.validDataProvider,
			expPass:      true,
		},
		{
			description:  "MsgAddDataProviderTestSuite: failing case - invalid feedId",
			feedId:       "",
			signer:       ts.signer,
			dataProvider: ts.validDataProvider,
			expPass:      false,
		},
		{
			description:  "MsgAddDataProviderTestSuite: failing case - data provider address and pubKey does not match",
			feedId:       "feedId1",
			signer:       ts.signer,
			dataProvider: ts.invalidDataProvider,
			expPass:      false,
		},
	}

	for i, tc := range testCases {
		msg := NewMsgAddDataProvider(
			tc.signer,
			tc.feedId,
			tc.dataProvider,
		)
		err := msg.ValidateBasic()

		if tc.expPass {
			ts.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.description)
		} else {
			ts.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.description)
		}
	}
}

type MsgRemoveDataProviderTestSuite struct {
	suite.Suite
	signer  sdk.AccAddress
	address sdk.AccAddress
}

func TestMsgRemoveDataProviderTestSuite(t *testing.T) {
	suite.Run(t, new(MsgRemoveDataProviderTestSuite))
}

func (ts *MsgRemoveDataProviderTestSuite) SetupTest() {
	_, _, signerAddr := GenerateAccount()
	ts.signer = signerAddr

	_, _, validAddr := GenerateAccount()
	ts.address = validAddr
}

func (ts *MsgRemoveDataProviderTestSuite) MsgRemoveDataProviderConstructor() {
	msg := NewMsgRemoveDataProvider(
		ts.signer,
		"feedId1",
		ts.address,
	)

	bz := ModuleCdc.MustMarshalJSON(msg)
	signedBytes := sdk.MustSortJSON(bz)

	ts.Require().Equal(msg.Route(), RouterKey)
	ts.Require().Equal(msg.Type(), RemoveDataProvider)
	ts.Require().Equal(msg.GetSigners(), []sdk.AccAddress{ts.signer})
	ts.Require().Equal(msg.GetSignBytes(), signedBytes)
}

func (ts *MsgRemoveDataProviderTestSuite) TestMsgRemoveDataProviderValidateBasic() {
	testCases := []struct {
		description string
		feedId      string
		signer      sdk.AccAddress
		address     sdk.AccAddress
		expPass     bool
	}{
		{
			description: "MsgAddDataProviderTestSuite: passing case - all valid values",
			feedId:      "feedId1",
			signer:      ts.signer,
			address:     ts.address,
			expPass:     true,
		},
		{
			description: "MsgAddDataProviderTestSuite: failing case - invalid feedId",
			feedId:      "",
			signer:      ts.signer,
			address:     ts.address,
			expPass:     false,
		},
		{
			description: "MsgAddDataProviderTestSuite: failing case - data provider address is empty",
			feedId:      "feedId1",
			signer:      ts.signer,
			address:     nil,
			expPass:     false,
		},
	}

	for i, tc := range testCases {
		msg := NewMsgRemoveDataProvider(
			tc.signer,
			tc.feedId,
			tc.address,
		)
		err := msg.ValidateBasic()

		if tc.expPass {
			ts.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.description)
		} else {
			ts.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.description)
		}
	}
}

type MsgSetSubmissionCountTestSuite struct {
	suite.Suite
	signer sdk.AccAddress
}

func TestMsgSetSubmissionCountTestSuite(t *testing.T) {
	suite.Run(t, new(MsgSetSubmissionCountTestSuite))
}

func (ts *MsgSetSubmissionCountTestSuite) SetupTest() {
	_, _, signerAddr := GenerateAccount()
	ts.signer = signerAddr
}

func (ts *MsgSetSubmissionCountTestSuite) MsgSetSubmissionCountConstructor() {
	msg := NewMsgSetSubmissionCount(
		ts.signer,
		"feedId1",
		uint32(1),
	)

	bz := ModuleCdc.MustMarshalJSON(msg)
	signedBytes := sdk.MustSortJSON(bz)

	ts.Require().Equal(msg.Route(), RouterKey)
	ts.Require().Equal(msg.Type(), SetSubmissionCount)
	ts.Require().Equal(msg.GetSigners(), []sdk.AccAddress{ts.signer})
	ts.Require().Equal(msg.GetSignBytes(), signedBytes)
}

func (ts *MsgSetSubmissionCountTestSuite) MsgSetSubmissionCountValidateBasic() {
	testCases := []struct {
		description     string
		feedId          string
		submissionCount uint32
		signer          sdk.AccAddress
		expPass         bool
	}{
		{
			description:     "MsgSetSubmissionCountTestSuite: passing case - all valid values",
			feedId:          "feedId1",
			signer:          ts.signer,
			submissionCount: uint32(1),
			expPass:         true,
		},
		{
			description:     "MsgSetSubmissionCountTestSuite: failing case - invalid feedId",
			feedId:          "",
			signer:          ts.signer,
			submissionCount: uint32(1),
			expPass:         false,
		},
		{
			description:     "MsgSetSubmissionCountTestSuite: failing case - submissionCount must not be 0",
			feedId:          "feedId1",
			signer:          ts.signer,
			submissionCount: uint32(0),
			expPass:         false,
		},
	}

	for i, tc := range testCases {
		msg := NewMsgSetSubmissionCount(
			tc.signer,
			tc.feedId,
			tc.submissionCount,
		)
		err := msg.ValidateBasic()

		if tc.expPass {
			ts.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.description)
		} else {
			ts.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.description)
		}
	}
}

type MsgSetHeartbeatTriggerTestSuite struct {
	suite.Suite
	signer sdk.AccAddress
}

func TestMsgSetHeartbeatTriggerTestSuite(t *testing.T) {
	suite.Run(t, new(MsgSetHeartbeatTriggerTestSuite))
}

func (ts *MsgSetHeartbeatTriggerTestSuite) SetupTest() {
	_, _, signerAddr := GenerateAccount()
	ts.signer = signerAddr
}

func (ts *MsgSetHeartbeatTriggerTestSuite) MsgSetHeartbeatTriggerConstructor() {
	msg := NewMsgSetHeartbeatTrigger(
		ts.signer,
		"feedId1",
		uint32(1),
	)

	bz := ModuleCdc.MustMarshalJSON(msg)
	signedBytes := sdk.MustSortJSON(bz)

	ts.Require().Equal(msg.Route(), RouterKey)
	ts.Require().Equal(msg.Type(), SetHeartbeatTrigger)
	ts.Require().Equal(msg.GetSigners(), []sdk.AccAddress{ts.signer})
	ts.Require().Equal(msg.GetSignBytes(), signedBytes)
}

func (ts *MsgSetHeartbeatTriggerTestSuite) MsgSetHeartbeatTriggeralidateBasic() {
	testCases := []struct {
		description      string
		feedId           string
		heartbeatTrigger uint32
		signer           sdk.AccAddress
		expPass          bool
	}{
		{
			description:      "MsgSetHeartbeatTriggerTestSuite: passing case - all valid values",
			feedId:           "feedId1",
			signer:           ts.signer,
			heartbeatTrigger: uint32(1),
			expPass:          true,
		},
		{
			description:      "MsgSetHeartbeatTriggerTestSuite: failing case - invalid feedId",
			feedId:           "",
			signer:           ts.signer,
			heartbeatTrigger: uint32(1),
			expPass:          false,
		},
		{
			description:      "MsgSetHeartbeatTriggerTestSuite: failing case - heartbeatTrigger must not be 0",
			feedId:           "feedId1",
			signer:           ts.signer,
			heartbeatTrigger: uint32(0),
			expPass:          false,
		},
	}

	for i, tc := range testCases {
		msg := NewMsgSetHeartbeatTrigger(
			tc.signer,
			tc.feedId,
			tc.heartbeatTrigger,
		)
		err := msg.ValidateBasic()

		if tc.expPass {
			ts.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.description)
		} else {
			ts.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.description)
		}
	}
}

type MsgSetDeviationThresholdTestSuite struct {
	suite.Suite
	signer sdk.AccAddress
}

func TestNewMsgSetDeviationThresholdTestSuite(t *testing.T) {
	suite.Run(t, new(MsgSetDeviationThresholdTestSuite))
}

func (ts *MsgSetDeviationThresholdTestSuite) SetupTest() {
	_, _, signerAddr := GenerateAccount()
	ts.signer = signerAddr
}

func (ts *MsgSetDeviationThresholdTestSuite) MsgSetDeviationThresholdConstructor() {
	msg := NewMsgSetDeviationThreshold(
		ts.signer,
		"feedId1",
		uint32(1),
	)

	bz := ModuleCdc.MustMarshalJSON(msg)
	signedBytes := sdk.MustSortJSON(bz)

	ts.Require().Equal(msg.Route(), RouterKey)
	ts.Require().Equal(msg.Type(), SetDeviationThresholdTrigger)
	ts.Require().Equal(msg.GetSigners(), []sdk.AccAddress{ts.signer})
	ts.Require().Equal(msg.GetSignBytes(), signedBytes)
}

func (ts *MsgSetDeviationThresholdTestSuite) MsgSetDeviationThresholdValidateBasic() {
	testCases := []struct {
		description               string
		feedId                    string
		deviationThresholdTrigger uint32
		signer                    sdk.AccAddress
		expPass                   bool
	}{
		{
			description:               "MsgSetHeartbeatTriggerTestSuite: passing case - all valid values",
			feedId:                    "feedId1",
			signer:                    ts.signer,
			deviationThresholdTrigger: uint32(1),
			expPass:                   true,
		},
		{
			description:               "MsgSetHeartbeatTriggerTestSuite: failing case - invalid feedId",
			feedId:                    "",
			signer:                    ts.signer,
			deviationThresholdTrigger: uint32(1),
			expPass:                   false,
		},
		{
			description:               "MsgSetHeartbeatTriggerTestSuite: failing case - deviationThresholdTrigger must not be 0",
			feedId:                    "feedId1",
			signer:                    ts.signer,
			deviationThresholdTrigger: uint32(0),
			expPass:                   false,
		},
	}

	for i, tc := range testCases {
		msg := NewMsgSetDeviationThreshold(
			tc.signer,
			tc.feedId,
			tc.deviationThresholdTrigger,
		)
		err := msg.ValidateBasic()

		if tc.expPass {
			ts.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.description)
		} else {
			ts.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.description)
		}
	}
}

type MsgSetFeedRewardTestSuite struct {
	suite.Suite
	signer sdk.AccAddress
}

func TestMsgSetFeedRewardTestSuite(t *testing.T) {
	suite.Run(t, new(MsgSetFeedRewardTestSuite))
}

func (ts *MsgSetFeedRewardTestSuite) SetupTest() {
	_, _, signerAddr := GenerateAccount()
	ts.signer = signerAddr
}

func (ts *MsgSetFeedRewardTestSuite) MsgSetFeedRewardConstructor() {
	msg := NewMsgSetFeedReward(
		ts.signer,
		"feedId1",
		uint32(1),
	)

	bz := ModuleCdc.MustMarshalJSON(msg)
	signedBytes := sdk.MustSortJSON(bz)

	ts.Require().Equal(msg.Route(), RouterKey)
	ts.Require().Equal(msg.Type(), SetFeedReward)
	ts.Require().Equal(msg.GetSigners(), []sdk.AccAddress{ts.signer})
	ts.Require().Equal(msg.GetSignBytes(), signedBytes)
}

func (ts *MsgSetFeedRewardTestSuite) MsgSetFeedRewardValidateBasic() {
	testCases := []struct {
		description string
		feedId      string
		feedReward  uint32
		signer      sdk.AccAddress
		expPass     bool
	}{
		{
			description: "MsgSetFeedRewardTestSuite: passing case - all valid values",
			feedId:      "feedId1",
			signer:      ts.signer,
			feedReward:  uint32(1),
			expPass:     true,
		},
		{
			description: "MsgSetFeedRewardTestSuite: failing case - invalid feedId",
			feedId:      "",
			signer:      ts.signer,
			feedReward:  uint32(1),
			expPass:     false,
		},
		{
			description: "MsgSetFeedRewardTestSuite: failing case - feedReward must not be 0",
			feedId:      "feedId1",
			signer:      ts.signer,
			feedReward:  uint32(0),
			expPass:     false,
		},
	}

	for i, tc := range testCases {
		msg := NewMsgSetFeedReward(
			tc.signer,
			tc.feedId,
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

type MsgFeedOwnershipTransferTestSuite struct {
	suite.Suite
	signer       sdk.AccAddress
	newFeedOwner sdk.AccAddress
}

func TestMsgFeedOwnershipTransferTestSuite(t *testing.T) {
	suite.Run(t, new(MsgFeedOwnershipTransferTestSuite))
}

func (ts *MsgFeedOwnershipTransferTestSuite) SetupTest() {
	_, _, signerAddr := GenerateAccount()
	ts.signer = signerAddr

	_, _, newFeedOwnerAddr := GenerateAccount()
	ts.newFeedOwner = newFeedOwnerAddr
}

func (ts *MsgFeedOwnershipTransferTestSuite) MsgFeedOwnershipTransferConstructor() {
	msg := NewMsgFeedOwnershipTransfer(
		ts.signer,
		"feedId1",
		ts.newFeedOwner,
	)

	bz := ModuleCdc.MustMarshalJSON(msg)
	signedBytes := sdk.MustSortJSON(bz)

	ts.Require().Equal(msg.Route(), RouterKey)
	ts.Require().Equal(msg.Type(), FeedOwnershipTransfer)
	ts.Require().Equal(msg.GetSigners(), []sdk.AccAddress{ts.signer})
	ts.Require().Equal(msg.GetSignBytes(), signedBytes)
}

func (ts *MsgFeedOwnershipTransferTestSuite) MsgFeedOwnershipTransferValidateBasic() {
	testCases := []struct {
		description  string
		feedId       string
		signer       sdk.AccAddress
		newFeedOwner sdk.AccAddress
		expPass      bool
	}{
		{
			description:  "MsgFeedOwnershipTransferTestSuite: passing case - all valid values",
			feedId:       "feedId1",
			signer:       ts.signer,
			newFeedOwner: ts.newFeedOwner,
			expPass:      true,
		},
		{
			description:  "MsgFeedOwnershipTransferTestSuite: failing case - signer can not be empty",
			feedId:       "feedId1",
			signer:       nil,
			newFeedOwner: ts.newFeedOwner,
			expPass:      false,
		},
		{
			description:  "MsgFeedOwnershipTransferTestSuite: failing case - feedId can not be empty",
			feedId:       "",
			signer:       ts.signer,
			newFeedOwner: ts.newFeedOwner,
			expPass:      false,
		},
	}

	for i, tc := range testCases {
		msg := NewMsgFeedOwnershipTransfer(
			tc.signer,
			tc.feedId,
			tc.newFeedOwner,
		)
		err := msg.ValidateBasic()

		if tc.expPass {
			ts.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.description)
		} else {
			ts.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.description)
		}
	}
}

type MsgRequestNewRoundTestSuite struct {
	suite.Suite
	signer       sdk.AccAddress
	newFeedOwner sdk.AccAddress
}

func TestMsgRequestNewRoundTestSuite(t *testing.T) {
	suite.Run(t, new(MsgRequestNewRoundTestSuite))
}

func (ts *MsgRequestNewRoundTestSuite) SetupTest() {
	_, _, signerAddr := GenerateAccount()
	ts.signer = signerAddr

	_, _, newFeedOwnerAddr := GenerateAccount()
	ts.newFeedOwner = newFeedOwnerAddr
}

func (ts *MsgRequestNewRoundTestSuite) MsgRequestNewRoundConstructor() {
	msg := NewMsgRequestNewRound(
		ts.signer,
		"feedId1",
	)

	bz := ModuleCdc.MustMarshalJSON(msg)
	signedBytes := sdk.MustSortJSON(bz)

	ts.Require().Equal(msg.Route(), RouterKey)
	ts.Require().Equal(msg.Type(), FeedOwnershipTransfer)
	ts.Require().Equal(msg.GetSigners(), []sdk.AccAddress{ts.signer})
	ts.Require().Equal(msg.GetSignBytes(), signedBytes)
}

func (ts *MsgRequestNewRoundTestSuite) MsgRequestNewRoundValidateBasic() {
	testCases := []struct {
		description  string
		feedId       string
		signer       sdk.AccAddress
		newFeedOwner sdk.AccAddress
		expPass      bool
	}{
		{
			description: "MsgRequestNewRoundTestSuite: passing case - all valid values",
			feedId:      "feedId1",
			signer:      ts.signer,
			expPass:     true,
		},
		{
			description: "MsgRequestNewRoundTestSuite: failing case - signer can not be empty",
			feedId:      "feedId1",
			signer:      nil,
			expPass:     false,
		},
		{
			description: "MsgRequestNewRoundTestSuite: failing case - feedId can not be empty",
			feedId:      "",
			signer:      ts.signer,
			expPass:     false,
		},
	}

	for i, tc := range testCases {
		msg := NewMsgRequestNewRound(
			tc.signer,
			tc.feedId,
		)
		err := msg.ValidateBasic()

		if tc.expPass {
			ts.Require().NoError(err, "valid test %d failed: %s, %v", i, tc.description)
		} else {
			ts.Require().Error(err, "invalid test %d passed: %s, %v", i, tc.description)
		}
	}
}
