package types

import (
	"errors"
	githubcosmossdktypes "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	SubmitFeedData          = "SubmitFeedData"
	AddModuleOwner          = "AddModuleOwner"
	ModuleOwnershipTransfer = "ModuleOwnershipTransfer"
	AddFeed                 = "AddFeed"
)

var _, _, _, _ sdk.Msg = &MsgFeedData{}, &MsgModuleOwnershipTransfer{}, &MsgModuleOwner{}, &MsgFeed{}
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

	if m.Submitter.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "submitter can not be empty")
	}
	if m.FeedId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "feedId can not be empty")
	}
	if len(m.FeedData) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "feedData can not be empty")
	}

	// TODO: verify the number of required signatures here
	if len(m.GetSignatures()) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "number of oracle signatures does not meet the required number")
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
	// TODO: add proper cosmos address and pubkey validation
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

type MsgModuleOwners []*MsgModuleOwner

// Contains returns true if the given address exists in a slice of ModuleOwners objects.
func (mo MsgModuleOwners) Contains(addr sdk.Address) bool {
	for _, acc := range mo {
		if acc.GetAddress().Equals(addr) {
			return true
		}
	}

	return false
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
	// TODO: add proper cosmos address and pubkey validation
	return nil
}

func (m *MsgModuleOwnershipTransfer) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgModuleOwnershipTransfer) GetSigners() []githubcosmossdktypes.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.AssignerAddress)}
}

func NewMsgFeed(feedId string, feedOwner, moduleOwner sdk.Address, initDataProviders []*DataProvider, submissionCount, heartbeatTrigger, deviationThresholdTrigger uint32) *MsgFeed {
	return &MsgFeed{
		FeedId:                    feedId,
		FeedOwner:                 feedOwner.Bytes(),
		DataProviders:             initDataProviders,
		SubmissionCount:           submissionCount,
		HeartbeatTrigger:          heartbeatTrigger,
		DeviationThresholdTrigger: deviationThresholdTrigger,
		ModuleOwnerAddress:        moduleOwner.Bytes(),
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
		return errors.New("invalid module owner")
	}
	if len(m.GetFeedId()) == 0 {
		return errors.New("invalid feedId")
	}
	if m.GetFeedOwner().Empty() {
		return errors.New("invalid feed owner")
	}
	if m.GetSubmissionCount() == 0 {
		return errors.New("SubmissionCount must not be 0")
	}
	if m.GetHeartbeatTrigger() == 0 {
		return errors.New("HeartbeatTrigger must not be 0")
	}
	if m.GetDeviationThresholdTrigger() == 0 {
		return errors.New("DeviationThresholdTrigger must not be 0")
	}

	if len(m.GetDataProviders()) == 0 {
		return errors.New("init data provider must not empty")
	}
	for _, provider := range m.GetDataProviders() {
		if !provider.Verify() {
			return errors.New("init data provider address and pubKey does not match")
		}
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
