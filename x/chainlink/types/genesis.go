// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
// This is where the init genesis can be defined
func DefaultGenesis() *GenesisState {
	return &GenesisState{ModuleOwners: nil}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if len(gs.GetModuleOwners()) == 0 {
		return errors.New("module owner size cannot be the zero")
	}

	for _, moduleOwner := range gs.GetModuleOwners() {
		err := moduleOwner.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m MsgModuleOwner) Validate() error {
	if m.GetAddress().Empty() {
		return errors.New("module owner address cannot be the empty")
	}
	if len(m.GetPubKey()) == 0 {
		return errors.New("module owner public key cannot be the empty")
	}

	bech32PubKey := sdk.MustGetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, string(m.PubKey))
	if !bytes.Equal(bech32PubKey.Address().Bytes(), m.Address.Bytes()) {
		return errors.New("module owner address and pubKey not match")
	}

	return nil
}

// GetGenesisStateFromAppState returns chainlink module GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONMarshaler, appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return &genesisState
}
