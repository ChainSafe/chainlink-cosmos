package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	SubmitFeedData = "SubmitFeedData"
)

var _ sdk.Msg = &MsgFeedData{}

func NewMsgFeed(submitter sdk.Address, feedId string, feedData string) *MsgFeedData {
	return &MsgFeedData{
		FeedId:    feedId,
		Submitter: submitter.Bytes(),
		FeedData:  feedData,
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
	if m.FeedId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "feedId can not be empty")
	}
	if m.FeedData == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "feedData can not be empty")
	}
	if m.Submitter.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "submitter can not be empty")
	}
	return nil
}
