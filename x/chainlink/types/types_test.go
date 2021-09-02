// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTypes_MsgModuleOwners_Contains(t *testing.T) {
	_, _, addr := GenerateAccount()
	assignerAddress := addr
	// assignerPublicKey = []byte(pubkey)

	_, pubkey1, addr1 := GenerateAccount()
	newModuleOwnerAddress1 := addr1
	newModuleOwnerPublicKey1 := []byte(pubkey1)

	_, pubkey2, addr2 := GenerateAccount()
	newModuleOwnerAddress2 := addr2
	newModuleOwnerPublicKey2 := []byte(pubkey2)

	var mos MsgModuleOwners

	mo := NewMsgModuleOwner(
		assignerAddress,
		newModuleOwnerAddress1,
		[]byte(newModuleOwnerPublicKey1),
	)
	err := mo.ValidateBasic()
	require.NoError(t, err)

	mos = append(mos, mo)
	require.Equal(t, 1, len(mos))

	require.True(t, mos.Contains(newModuleOwnerAddress1))
	require.False(t, mos.Contains(newModuleOwnerAddress2))

	mo = NewMsgModuleOwner(
		assignerAddress,
		newModuleOwnerAddress2,
		[]byte(newModuleOwnerPublicKey2),
	)
	err = mo.ValidateBasic()
	require.NoError(t, err)

	mos = append(mos, mo)
	require.Equal(t, 2, len(mos))

	require.True(t, mos.Contains(newModuleOwnerAddress1))
	require.True(t, mos.Contains(newModuleOwnerAddress2))
}

func TestTypes_DataProvider_Verify(t *testing.T) {
	_, pubkey, addr := GenerateAccount()
	dataProviderAddress := addr
	dataProviderPublicKey := []byte(pubkey)

	dp := &DataProvider{
		Address: dataProviderAddress,
		PubKey:  []byte(dataProviderPublicKey),
	}

	require.Equal(t, dataProviderAddress, dp.GetAddress())
	require.Equal(t, dataProviderPublicKey, dp.GetPubKey())

	require.True(t, dp.Verify())

	_, _, iaddr := GenerateAccount()
	_, ipubkey, _ := GenerateAccount()
	invalidModOwnerAddress := iaddr
	invalidModOwnerPublicKey := []byte(ipubkey)

	idp := &DataProvider{
		Address: invalidModOwnerAddress,
		PubKey:  []byte(invalidModOwnerPublicKey),
	}

	require.Equal(t, invalidModOwnerAddress, idp.GetAddress())
	require.Equal(t, invalidModOwnerPublicKey, idp.GetPubKey())

	require.False(t, idp.Verify())
}

func TestTypes_DataProviders_Contains_Remove(t *testing.T) {
	_, pubkey1, addr1 := GenerateAccount()
	dataProviderAddress1 := addr1
	dataProviderPublicKey1 := []byte(pubkey1)

	_, pubkey2, addr2 := GenerateAccount()
	dataProviderAddress2 := addr2
	dataProviderPublicKey2 := []byte(pubkey2)

	var dps DataProviders

	dp := &DataProvider{
		Address: dataProviderAddress1,
		PubKey:  []byte(dataProviderPublicKey1),
	}
	require.True(t, dp.Verify())

	dps = append(dps, dp)
	require.Equal(t, 1, len(dps))

	require.True(t, dps.Contains(dataProviderAddress1))
	require.False(t, dps.Contains(dataProviderAddress2))

	dp = &DataProvider{
		Address: dataProviderAddress2,
		PubKey:  []byte(dataProviderPublicKey2),
	}
	require.True(t, dp.Verify())

	dps = append(dps, dp)
	require.Equal(t, 2, len(dps))

	require.True(t, dps.Contains(dataProviderAddress1))
	require.True(t, dps.Contains(dataProviderAddress2))

	dps = dps.Remove(dataProviderAddress1)
	require.Equal(t, 1, len(dps))
	require.False(t, dps.Contains(dataProviderAddress1))
	require.True(t, dps.Contains(dataProviderAddress2))

	dps = dps.Remove(dataProviderAddress2)
	require.False(t, dps.Contains(dataProviderAddress1))
	require.False(t, dps.Contains(dataProviderAddress2))
}

func TestTypes_DeriveCosmosAddrFromPubKey(t *testing.T) {
	_, pubkey1, addr1 := GenerateAccount()
	_, pubkey2, addr2 := GenerateAccount()

	require.NotEqual(t, pubkey1, pubkey2)
	require.NotEqual(t, addr1, addr2)

	expAddr1, err := DeriveCosmosAddrFromPubKey(pubkey1)
	require.NoError(t, err)
	require.Equal(t, expAddr1, addr1)

	expAddr2, err := DeriveCosmosAddrFromPubKey(pubkey2)
	require.NoError(t, err)
	require.Equal(t, expAddr2, addr2)
}
