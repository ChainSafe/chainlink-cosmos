package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgFeed{}

func NewMsgFeed(creator sdk.Address, feedId string, feedData string) *MsgFeed {
	return &MsgFeed{
		FeedId:    feedId,
		Submitter: creator.Bytes(),
		FeedData:  feedData,
	}
}

func (m *MsgFeed) Route() string {
	return RouterKey
}

func (m *MsgFeed) Type() string {
	return "CreateFeed"
}

func (m *MsgFeed) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(m.Submitter)}
}

func (m *MsgFeed) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

func (m *MsgFeed) ValidateBasic() error {
	if m.FeedId == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "feedId can not be empty")
	}
	if m.FeedData == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "feed data can not be empty")
	}
	if m.Submitter.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "submitter can not be empty")
	}
	return nil
}
