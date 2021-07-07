package types

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgModuleOwners []*MsgModuleOwner

// Contains returns true if the given address exists in a slice of ModuleOwners.
func (mo MsgModuleOwners) Contains(addr sdk.Address) bool {
	for _, acc := range mo {
		if acc.GetAddress().Equals(addr) {
			return true
		}
	}

	return false
}

func (m *DataProvider) Verify() bool {
	bech32PubKey := sdk.MustGetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, string(m.GetPubKey()))

	// address and pubKey must match
	return bytes.Equal(bech32PubKey.Address().Bytes(), m.GetAddress().Bytes())
}

type DataProviders []*DataProvider

// Contains returns true if the given address exists in a slice of DataProviders.
func (dp DataProviders) Contains(addr sdk.Address) bool {
	for _, acc := range dp {
		if acc.GetAddress().Equals(addr) {
			return true
		}
	}

	return false
}
