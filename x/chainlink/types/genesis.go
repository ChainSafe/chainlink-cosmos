package types

import (
	"encoding/json"
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
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
	// TODO: add proper module owner validation
	if len(gs.GetModuleOwners()) == 0 {
		return errors.New("module owner size cannot be the zero")
	}

	for _, owner := range gs.GetModuleOwners() {
		err := owner.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m ModuleOwner) Validate() error {
	// TODO: add proper cosmos address and pubkey validation
	if len(m.GetAddress()) == 0 {
		return errors.New("module owner address cannot be the empty")
	}
	if len(m.GetPubKey()) == 0 {
		return errors.New("module owner public key cannot be the empty")
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
