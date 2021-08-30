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
	SetFeedReward                = "SetFeedReward"
	FeedOwnershipTransfer        = "FeedOwnershipTransfer"
	RequestNewRound              = "RequestNewRound"
	SetAccountPiggyAddress       = "SetAccountPiggyAddress"
)

var _, _, _, _, _, _, _, _, _, _, _, _ sdk.Msg = &MsgFeedData{}, &MsgModuleOwnershipTransfer{}, &MsgModuleOwner{},
	&MsgFeed{}, &MsgAddDataProvider{}, &MsgRemoveDataProvider{}, &MsgSetSubmissionCount{}, &MsgSetHeartbeatTrigger{},
	&MsgSetDeviationThresholdTrigger{}, &MsgFeedOwnershipTransfer{}, &MsgRequestNewRound{}, &MsgAccount{}

var _ sdk.Tx = &MsgModuleOwner{}

var _ Validation = &MsgFeedData{}

func NewMsgFeedData(submitter sdk.Address, feedId string, feedData []byte, signatures [][]byte, cosmosPubKeys [][]byte) *MsgFeedData {
	return &MsgFeedData{
		FeedId:     feedId,
		Submitter:  submitter.Bytes(),
		FeedData:   feedData,
		Signatures: signatures,
		// IsFeedDataValid will be true by default
		IsFeedDataValid: true,
		CosmosPubKeys:   cosmosPubKeys,
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
	if len(m.GetSignatures()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "number of oracle signatures does not meet the required number")
	}
	if len(m.Signatures) != len(m.CosmosPubKeys) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "number of oracle signatures does not match the data provider cosmos pubkey number")
	}

	return nil
}

func (m *MsgFeedData) Validate(fn func(msg sdk.Msg) bool) bool {
	if fn == nil {
		return true
	}
	return fn(m)
}

// RewardCalculator calculates the reward for each data provider in the current submit feed data tx
// base on the registered reward strategy
// return the slice of reward payout and the total reward amount
func (m *MsgFeedData) RewardCalculator(feed *MsgFeed, feedData *MsgFeedData) ([]RewardPayout, uint64, error) {
	// every one gets the base amount if no strategy configured when chain launching
	// or the owner of the current feed does not set a strategy
	if len(FeedRewardStrategyConvertor) == 0 || feed.GetFeedReward().GetStrategy() == "" {
		rewardPayout := make([]RewardPayout, 0, len(feedData.GetCosmosPubKeys()))

		for i := 0; i < len(feedData.GetCosmosPubKeys()); i++ {
			// err is not possible here since pubkey has been checked in anteHandler
			cosmosAddr, _ := DeriveCosmosAddrFromPubKey(string(feedData.GetCosmosPubKeys()[i]))
			dataProviderAddr, _ := sdk.AccAddressFromBech32(cosmosAddr.String())

			rp := RewardPayout{
				DataProvider: &DataProvider{
					Address: dataProviderAddr,
				},
				Amount: feed.GetFeedReward().GetAmount(),
			}
			rewardPayout = append(rewardPayout, rp)
		}
		return rewardPayout, feed.GetFeedReward().GetAmount() * uint64(len(feedData.GetSignatures())), nil
	}

	// strategy of a feed here has already been checked in anteHandler when set, ok must be true
	strategyFn, _ := FeedRewardStrategyConvertor[feed.GetFeedReward().GetStrategy()] // nolint

	RewardPayoutList, err := strategyFn(feed, feedData)
	if err != nil {
		return nil, 0, err
	}

	totalRewardAmount := uint64(0)
	for _, payout := range RewardPayoutList {
		totalRewardAmount += payout.Amount
	}

	return RewardPayoutList, totalRewardAmount, nil
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

func NewMsgFeed(feedId, feedDesc string, feedOwner, moduleOwner sdk.Address, initDataProviders []*DataProvider,
	submissionCount, heartbeatTrigger, deviationThresholdTrigger uint32, baseFeedRewardAmount uint64, feedRewardStrategy string) *MsgFeed {
	return &MsgFeed{
		FeedId:                    feedId,
		Desc:                      feedDesc,
		FeedOwner:                 feedOwner.Bytes(),
		DataProviders:             initDataProviders,
		SubmissionCount:           submissionCount,
		HeartbeatTrigger:          heartbeatTrigger,
		DeviationThresholdTrigger: deviationThresholdTrigger,
		ModuleOwnerAddress:        moduleOwner.Bytes(),
		FeedReward: &FeedRewardSchema{
			Amount:   baseFeedRewardAmount,
			Strategy: feedRewardStrategy,
		},
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
	if m.GetFeedReward().GetAmount() == 0 {
		return errors.New("baseFeedRewardAmount must not be 0")
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

func NewMsgSetFeedReward(signer githubcosmossdktypes.AccAddress, feedId string, baseFeedRewardAmount uint64, feedRewardStrategy string) *MsgSetFeedReward {
	return &MsgSetFeedReward{
		FeedId: feedId,
		FeedReward: &FeedRewardSchema{
			Amount:   baseFeedRewardAmount,
			Strategy: feedRewardStrategy,
		},
		Signer: signer,
	}
}

func (m *MsgSetFeedReward) Route() string {
	return RouterKey
}

func (m *MsgSetFeedReward) Type() string {
	return SetFeedReward
}

func (m *MsgSetFeedReward) ValidateBasic() error {
	if m.GetSigner().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "signer can not be empty")
	}
	if len(m.GetFeedId()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "feedId can not be empty")
	}
	if m.GetFeedReward().GetAmount() == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "baseFeedRewardAmount must not be 0")
	}
	return nil
}

func (m *MsgSetFeedReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgSetFeedReward) GetSigners() []githubcosmossdktypes.AccAddress {
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

func NewMsgRequestNewRound(signer githubcosmossdktypes.AccAddress, feedId string) *MsgRequestNewRound {
	return &MsgRequestNewRound{
		FeedId: feedId,
		Signer: signer,
	}
}

func (m *MsgRequestNewRound) Route() string {
	return RouterKey
}

func (m *MsgRequestNewRound) Type() string {
	return RequestNewRound
}

func (m *MsgRequestNewRound) ValidateBasic() error {
	if m.GetSigner().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "signer can not be empty")
	}
	if len(m.GetFeedId()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "feedId can not be empty")
	}
	return nil
}

func (m *MsgRequestNewRound) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgRequestNewRound) GetSigners() []githubcosmossdktypes.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

func NewMsgAddAccount(submitter githubcosmossdktypes.AccAddress, chainlinkPublicKey, chainlinkSigningKey []byte, piggyAddress githubcosmossdktypes.AccAddress) *MsgAccount {
	return &MsgAccount{
		Submitter:           submitter,
		ChainlinkPublicKey:  chainlinkPublicKey,
		ChainlinkSigningKey: chainlinkSigningKey,
		PiggyAddress:        piggyAddress,
	}
}

func (m *MsgAccount) Route() string {
	return RouterKey
}

func (m *MsgAccount) Type() string {
	return RequestNewRound
}

func (m *MsgAccount) ValidateBasic() error {
	if m.GetSubmitter().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "submitter can not be empty")
	}
	if len(m.GetChainlinkPublicKey()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "chainlink public key can not be empty")
	}
	if len(m.GetChainlinkSigningKey()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "chainlink signing key can not be empty")
	}
	if m.GetPiggyAddress().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "piggy address can not be empty")
	}

	return nil
}

func (m *MsgAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgAccount) GetSigners() []githubcosmossdktypes.AccAddress {
	return []sdk.AccAddress{m.Submitter}
}

func NewMsgEditAccount(submitter githubcosmossdktypes.AccAddress, piggyAddress githubcosmossdktypes.AccAddress) *MsgEditAccount {
	return &MsgEditAccount{
		Submitter:    submitter,
		PiggyAddress: piggyAddress,
	}
}

func (m *MsgEditAccount) Route() string {
	return RouterKey
}

func (m *MsgEditAccount) Type() string {
	return RequestNewRound
}

func (m *MsgEditAccount) ValidateBasic() error {
	if m.GetSubmitter().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "submitter can not be empty")
	}
	if m.GetPiggyAddress().Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "piggy address can not be empty")
	}

	return nil
}

func (m *MsgEditAccount) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgEditAccount) GetSigners() []githubcosmossdktypes.AccAddress {
	return []sdk.AccAddress{m.Submitter}
}
