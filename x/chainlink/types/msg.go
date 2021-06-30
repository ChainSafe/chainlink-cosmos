package types

import (
	githubcosmossdktypes "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	SubmitFeedData          = "SubmitFeedData"
	AddModuleOwner          = "AddModuleOwner"
	ModuleOwnershipTransfer = "ModuleOwnershipTransfer"
)

var _, _, _ sdk.Msg = &MsgFeedData{}, &MsgModuleOwnershipTransfer{}, &MsgModuleOwner{}
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

func NewModuleOwner(assigner, address sdk.Address, pubKey []byte) *MsgModuleOwner {
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
	panic("implement me")
}

func (m *MsgModuleOwnershipTransfer) Type() string {
	return ModuleOwnershipTransfer
}

func (m *MsgModuleOwnershipTransfer) ValidateBasic() error {
	panic("implement me")
}

func (m *MsgModuleOwnershipTransfer) GetSignBytes() []byte {
	panic("implement me")
}

func (m *MsgModuleOwnershipTransfer) GetSigners() []githubcosmossdktypes.AccAddress {
	panic("implement me")
}
