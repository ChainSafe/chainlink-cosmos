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
		{description: "MsgFeedDataTestSuite: passing case - all valid values", submitter: ts.submitter, feedId: ts.feedId, feedData: ts.feedData, signatures: ts.signatures, expPass: true},
		{description: "MsgFeedDataTestSuite: failing case - empty submitter", submitter: nil, feedId: ts.feedId, feedData: ts.feedData, signatures: ts.signatures, expPass: false},
		{description: "MsgFeedDataTestSuite: failing case - empty feedId", submitter: ts.submitter, feedId: "", feedData: ts.feedData, signatures: ts.signatures, expPass: false},
		{description: "MsgFeedDataTestSuite: failing case - invalid feedId format", submitter: ts.submitter, feedId: "BAD/FEED/ID", feedData: ts.feedData, signatures: ts.signatures, expPass: false},
		{description: "MsgFeedDataTestSuite: failing case - empty feedData", submitter: ts.submitter, feedId: ts.feedId, feedData: nil, signatures: ts.signatures, expPass: false},
		{description: "MsgFeedDataTestSuite: failing case - empty signatures", submitter: ts.submitter, feedId: ts.feedId, feedData: ts.feedData, signatures: [][]byte{}, expPass: false},
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
		{description: "MsgModuleOwnerTestSuite: passing case - all valid values", assigner: ts.assignerAddress, address: ts.newModuleOwnerAddress, publicKey: []byte(ts.newModuleOwnerPublicKey), expPass: true},
		{description: "MsgModuleOwnerTestSuite: failing case - address and publicKey does not match", assigner: ts.assignerAddress, address: ts.newModuleOwnerAddress, publicKey: ts.invalidModOwnerPublicKey, expPass: false},
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
	assignerAddress          sdk.AccAddress
	assignerPublicKey        []byte
	newModuleOwnerAddress    sdk.AccAddress
	newModuleOwnerPublicKey  []byte
	invalidModOwnerAddress   sdk.AccAddress
	invalidModOwnerPublicKey []byte
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
	ts.invalidModOwnerAddress = addr
	ts.invalidModOwnerPublicKey = []byte(pubkey)
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
		{description: "MsgModuleOwnershipTransferTestSuite: passing case - all valid values", assigner: ts.assignerAddress, address: ts.newModuleOwnerAddress, publicKey: []byte(ts.newModuleOwnerPublicKey), expPass: true},
		{description: "MsgModuleOwnershipTransferTestSuite: failing case - assigner address is empty", assigner: nil, address: ts.newModuleOwnerAddress, publicKey: ts.invalidModOwnerPublicKey, expPass: false},
		{description: "MsgModuleOwnershipTransferTestSuite: failing case - address and publicKey does not match", assigner: ts.assignerAddress, address: ts.newModuleOwnerAddress, publicKey: ts.invalidModOwnerPublicKey, expPass: false},
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
