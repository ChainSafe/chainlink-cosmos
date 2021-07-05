package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m *DataProvider) Verify() bool {
	bech32PubKey := sdk.MustGetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, string(m.GetPubKey()))

	// address and pubKey must match
	return bytes.Equal(bech32PubKey.Address().Bytes(), m.GetAddress().Bytes())
}
