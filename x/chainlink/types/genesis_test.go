// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypes_GenesisState_Validate(t *testing.T) {
	genstate := DefaultGenesis()
	emptyGenesis := &GenesisState{ModuleOwners: nil}
	require.Equal(t, genstate, emptyGenesis)
	require.Error(t, genstate.Validate())

	genstate.ModuleOwners = make([]*MsgModuleOwner, 0)
	require.Error(t, genstate.Validate())

	_, pubKey, addr := GenerateAccount()

	genstate.ModuleOwners = append(genstate.ModuleOwners, &MsgModuleOwner{Address: addr, PubKey: []byte(pubKey), AssignerAddress: nil})

	require.NoError(t, genstate.Validate())
}

func TestTypes_MsgModuleOwner_Validate(t *testing.T) {
	_, validPubKey, validAddr := GenerateAccount()

	mo := &MsgModuleOwner{Address: validAddr, PubKey: []byte(validPubKey), AssignerAddress: nil}

	require.NoError(t, mo.Validate())

	_, invalidPubKey, _ := GenerateAccount()
	_, _, invalidAddr := GenerateAccount()

	imo := &MsgModuleOwner{Address: invalidAddr, PubKey: []byte(invalidPubKey), AssignerAddress: nil}

	require.Error(t, imo.Validate())
}

func TestTypes_GetGenesisStateFromAppState(t *testing.T) {
	// TODO: figure out how to get appstate to pass into GetGenesisStateFromAppState
	// issue with cyclic import if i import chainlink-cosmos/app

	// userHomeDir, err := os.UserHomeDir()
	// require.NoError(t, err)
	// home := filepath.Join(userHomeDir, ".chainlinkd")
	// encodingConfig := &params.EncodingConfig{}
	// clientCtx := client.Context{}.
	// 	WithHomeDir(home).
	// 	WithViper("").
	// 	WithAccountRetriever(authtypes.AccountRetriever{}).
	// 	WithInterfaceRegistry(encodingConfig.InterfaceRegistry)

	// depCdc := clientCtx.JSONMarshaler
	// cdc := depCdc.(codec.Marshaler)

	// appstate := map[string]json.RawMessage{}

	// genstate := GetGenesisStateFromAppState(cdc, appstate)

	// require.NoError(t, genstate.Validate())
}
