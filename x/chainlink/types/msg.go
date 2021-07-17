// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import (
	"bytes"
	"errors"
	"strings"

	githubcosmossdktypes "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	SubmitFeedData               = "SubmitFeedData"
	AddModuleOwner               = "AddModuleOwner"
	ModuleOwnershipTransfer      = "ModuleOwnershipTransfer"
	AddFeed                      = "AddFeed"
	AddDataProvider              = "AddDataProvider"
	RemoveDataProvider           = "RemoveDataProvider"
	SetSubmissionCount           = "SetSubmissionCount"
	SetHeartbeatTrigger          = "SetHeartbeatTrigger"
	SetDeviationThresholdTrigger = "SetDeviationThresholdTrigger"
	FeedOwnershipTransfer        = "FeedOwnershipTransfer"
)

var _, _, _, _, _, _, _, _, _, _ sdk.Msg = &MsgFeedData{}, &MsgModuleOwnershipTransfer{}, &MsgModuleOwner{}, &MsgFeed{}, &MsgAddDataProvider{}, &MsgRemoveDataProvider{}, &MsgSetSubmissionCount{}, &MsgSetHeartbeatTrigger{}, &MsgSetDeviationThresholdTrigger{}, &MsgFeedOwnershipTransfer{}
var _ sdk.Tx = &MsgModuleOwner{}

func NewMsgFeedData(submitter sdk.Address, feedId string, feedData []byte, signatures [][]byte) *MsgFeedData {
	return &MsgFeedData{
		FeedId:     feedId,
		Submitter:  submitter.Bytes(),
		FeedData:   feedData,
		Signatures: signatures,
	}
}

func (m *MsgFeedData) Route() string {
	return RouterKey
}

func (m *MsgFeedData) Type() string {
	return SubmitFeedData
}

func (m *MsgFeedData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.Submitter)}
}

func (m *MsgFeedData) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgFeedData) ValidateBasic() error {
	// TODO: add any basic input checking here

	if m.GetSubmitter().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "submitter can not be empty")
	}
	if len(m.GetFeedId()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "feedId can not be empty")
	}
	if strings.Contains(m.GetFeedId(), "/") {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "feedId can not contain character '/'")
	}
	if len(m.GetFeedData()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "feedData can not be empty")
	}

	// TODO: verify the number of required signatures here
	if len(m.GetSignatures()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "number of oracle signatures does not meet the required number")
	}
	return nil
}

func NewMsgModuleOwner(assigner, address sdk.Address, pubKey []byte) *MsgModuleOwner {
	mo := &MsgModuleOwner{
		Address: address.Bytes(),
		PubKey:  pubKey,
	}
	if assigner != nil {
		mo.AssignerAddress = assigner.Bytes()
	}

	return mo
}

func (m *MsgModuleOwner) Route() string {
	return RouterKey
}

func (m *MsgModuleOwner) Type() string {
	return AddModuleOwner
}

func (m *MsgModuleOwner) ValidateBasic() error {
	bech32PubKey := sdk.MustGetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, string(m.PubKey))
	if !bytes.Equal(bech32PubKey.Address().Bytes(), m.Address.Bytes()) {
		return errors.New("address and pubKey not match")
	}
	return nil
}

func (m *MsgModuleOwner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgModuleOwner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.AssignerAddress)}
}

func (m *MsgModuleOwner) GetMsgs() []githubcosmossdktypes.Msg {
	return []sdk.Msg{m}
}

func NewMsgModuleOwnershipTransfer(assigner, address sdk.Address, pubKey []byte) *MsgModuleOwnershipTransfer {
	return &MsgModuleOwnershipTransfer{
		AssignerAddress:       assigner.Bytes(),
		NewModuleOwnerAddress: address.Bytes(),
		NewModuleOwnerPubKey:  pubKey,
	}
}

func (m *MsgModuleOwnershipTransfer) Route() string {
	return RouterKey
}

func (m *MsgModuleOwnershipTransfer) Type() string {
	return ModuleOwnershipTransfer
}

func (m *MsgModuleOwnershipTransfer) ValidateBasic() error {
	if m.GetAssignerAddress().Empty() {
		return errors.New("assigner address can not be empty")
	}
	bech32PubKey := sdk.MustGetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, string(m.NewModuleOwnerPubKey))
	if !bytes.Equal(bech32PubKey.Address().Bytes(), m.NewModuleOwnerAddress.Bytes()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "new module owner address and pubKey does not match")
	}
	return nil
}

func (m *MsgModuleOwnershipTransfer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgModuleOwnershipTransfer) GetSigners() []githubcosmossdktypes.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.AssignerAddress)}
}

func NewMsgFeed(feedId string, feedOwner, moduleOwner sdk.Address, initDataProviders []*DataProvider, submissionCount, heartbeatTrigger, deviationThresholdTrigger, feedReward uint32) *MsgFeed {
	return &MsgFeed{
		FeedId:                    feedId,
		FeedOwner:                 feedOwner.Bytes(),
		DataProviders:             initDataProviders,
		SubmissionCount:           submissionCount,
		HeartbeatTrigger:          heartbeatTrigger,
		DeviationThresholdTrigger: deviationThresholdTrigger,
		ModuleOwnerAddress:        moduleOwner.Bytes(),
		FeedReward:                feedReward,
	}
}

func (m *MsgFeed) Route() string {
	return RouterKey
}

func (m *MsgFeed) Type() string {
	return AddFeed
}

func (m *MsgFeed) ValidateBasic() error {
	if m.GetModuleOwnerAddress().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "moduleOwner can not be empty")
	}
	if len(m.GetFeedId()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "feedId can not be empty")
	}
	if m.GetFeedOwner().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "feedOwner can not be empty")
	}
	if m.GetSubmissionCount() == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "submissionCount must not be 0")
	}
	if m.GetHeartbeatTrigger() == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "heartbeatTrigger must not be 0")
	}
	if m.GetDeviationThresholdTrigger() == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "deviationThresholdTrigger must not be 0")
	}
	if m.GetFeedReward() == 0 {
		return errors.New("FeedReward must not be 0")
	}

	if len(m.GetDataProviders()) == 0 {
		return errors.New("init data provider must not empty")
	}
	tmp := make(map[string][]byte)
	for _, provider := range m.GetDataProviders() {
		if !provider.Verify() {
			return errors.New("init data provider address and pubKey does not match")
		}
		tmp[provider.GetAddress().String()] = provider.GetPubKey()
	}
	if len(tmp) != len(m.GetDataProviders()) {
		return errors.New("init data provider list contains duplication")
	}
	return nil
}

func (m *MsgFeed) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgFeed) GetSigners() []githubcosmossdktypes.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.ModuleOwnerAddress)}
}

func (m *MsgFeed) Empty() bool {
	return m == nil
}

func NewMsgAddDataProvider(signer githubcosmossdktypes.AccAddress, feedId string, provider *DataProvider) *MsgAddDataProvider {
	return &MsgAddDataProvider{
		FeedId:       feedId,
		DataProvider: provider,
		Signer:       signer,
	}
}

func (m *MsgAddDataProvider) Route() string {
	return RouterKey
}

func (m *MsgAddDataProvider) Type() string {
	return AddDataProvider
}

func (m *MsgAddDataProvider) ValidateBasic() error {
	if len(m.GetFeedId()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid feedId")
	}
	provider := m.GetDataProvider()
	if !provider.Verify() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "data provider address and pubKey does not match")
	}
	return nil
}

func (m *MsgAddDataProvider) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgAddDataProvider) GetSigners() []githubcosmossdktypes.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

func NewMsgRemoveDataProvider(signer githubcosmossdktypes.AccAddress, feedId string, address githubcosmossdktypes.AccAddress) *MsgRemoveDataProvider {
	return &MsgRemoveDataProvider{
		FeedId:  feedId,
		Address: address,
		Signer:  signer,
	}
}

func (m *MsgRemoveDataProvider) Route() string {
	return RouterKey
}

func (m *MsgRemoveDataProvider) Type() string {
	return RemoveDataProvider
}

func (m *MsgRemoveDataProvider) ValidateBasic() error {
	if len(m.GetFeedId()) == 0 {
		return errors.New("invalid feedId")
	}
	if m.GetAddress().Empty() {
		return errors.New("data provider address is empty")
	}
	return nil
}

func (m *MsgRemoveDataProvider) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRemoveDataProvider) GetSigners() []githubcosmossdktypes.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

func NewMsgSetSubmissionCount(signer githubcosmossdktypes.AccAddress, feedId string, submissionCount uint32) *MsgSetSubmissionCount {
	return &MsgSetSubmissionCount{
		FeedId:          feedId,
		SubmissionCount: submissionCount,
		Signer:          signer,
	}
}

func (m *MsgSetSubmissionCount) Route() string {
	return RouterKey
}

func (m *MsgSetSubmissionCount) Type() string {
	return SetSubmissionCount
}

func (m *MsgSetSubmissionCount) ValidateBasic() error {
	if m.GetSigner().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "signer can not be empty")
	}
	if len(m.GetFeedId()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "feedId can not be empty")
	}
	if m.GetSubmissionCount() == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "submissionCount must not be 0")
	}
	return nil
}

func (m *MsgSetSubmissionCount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgSetSubmissionCount) GetSigners() []githubcosmossdktypes.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

func NewMsgSetHeartbeatTrigger(signer githubcosmossdktypes.AccAddress, feedId string, heartbeatTrigger uint32) *MsgSetHeartbeatTrigger {
	return &MsgSetHeartbeatTrigger{
		FeedId:           feedId,
		HeartbeatTrigger: heartbeatTrigger,
		Signer:           signer,
	}
}

func (m *MsgSetHeartbeatTrigger) Route() string {
	return RouterKey
}

func (m *MsgSetHeartbeatTrigger) Type() string {
	return SetHeartbeatTrigger
}

func (m *MsgSetHeartbeatTrigger) ValidateBasic() error {
	if m.GetSigner().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "signer can not be empty")
	}
	if len(m.GetFeedId()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "feedId can not be empty")
	}
	if m.GetHeartbeatTrigger() == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "heartbeatTrigger must not be 0")
	}
	return nil
}

func (m *MsgSetHeartbeatTrigger) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgSetHeartbeatTrigger) GetSigners() []githubcosmossdktypes.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

func NewMsgSetDeviationThreshold(signer githubcosmossdktypes.AccAddress, feedId string, deviationThresholdTrigger uint32) *MsgSetDeviationThresholdTrigger {
	return &MsgSetDeviationThresholdTrigger{
		FeedId:                    feedId,
		DeviationThresholdTrigger: deviationThresholdTrigger,
		Signer:                    signer,
	}
}

func (m *MsgSetDeviationThresholdTrigger) Route() string {
	return RouterKey
}

func (m *MsgSetDeviationThresholdTrigger) Type() string {
	return SetDeviationThresholdTrigger
}

func (m *MsgSetDeviationThresholdTrigger) ValidateBasic() error {
	if m.GetSigner().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "signer can not be empty")
	}
	if len(m.GetFeedId()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "feedId can not be empty")
	}
	if m.GetDeviationThresholdTrigger() == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "deviationThresholdTrigger must not be 0")
	}
	return nil
}

func (m *MsgSetDeviationThresholdTrigger) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgSetDeviationThresholdTrigger) GetSigners() []githubcosmossdktypes.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

func NewMsgFeedOwnershipTransfer(signer githubcosmossdktypes.AccAddress, feedId string, newFeedOwnerAddress sdk.AccAddress) *MsgFeedOwnershipTransfer {
	return &MsgFeedOwnershipTransfer{
		FeedId:              feedId,
		NewFeedOwnerAddress: newFeedOwnerAddress,
		Signer:              signer,
	}
}

func (m *MsgFeedOwnershipTransfer) Route() string {
	return RouterKey
}

func (m *MsgFeedOwnershipTransfer) Type() string {
	return FeedOwnershipTransfer
}

func (m *MsgFeedOwnershipTransfer) ValidateBasic() error {
	if m.GetSigner().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "signer can not be empty")
	}
	if len(m.GetFeedId()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "feedId can not be empty")
	}
	return nil
}

func (m *MsgFeedOwnershipTransfer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgFeedOwnershipTransfer) GetSigners() []githubcosmossdktypes.AccAddress {
	return []sdk.AccAddress{m.Signer}
}
