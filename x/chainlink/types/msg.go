package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	SubmitFeedData = "SubmitFeedData"
	AddModuleOwner = "AddModuleOwner"
)

var _ sdk.Msg = &MsgFeedData{}

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

var _ sdk.Msg = &ModuleOwner{}

func NewModuleOwner(address sdk.Address, pubKey []byte) *ModuleOwner {
	return &ModuleOwner{
		Address: address.Bytes(),
		PubKey:  pubKey,
	}
}

func (m *ModuleOwner) Route() string {
	return RouterKey
}

func (m *ModuleOwner) Type() string {
	return AddModuleOwner
}

func (m *ModuleOwner) ValidateBasic() error {
	// TODO: add proper cosmos address and pubkey validation
	return nil
}

func (m *ModuleOwner) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *ModuleOwner) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.Address)}
}
