// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestTypes_NewLinkCoin(t *testing.T) {
	hundredLink := NewLinkCoin(sdk.NewInt(100))
	require.NoError(t, hundredLink.Validate())
	require.True(t, hundredLink.IsValid())

	zeroLink := NewLinkCoin(sdk.NewInt(0))
	require.True(t, zeroLink.IsZero())
	require.True(t, zeroLink.IsValid())
	require.NoError(t, zeroLink.Validate())

	require.True(t, hundredLink.IsGTE(zeroLink))
	require.True(t, zeroLink.IsLT(hundredLink))

	require.True(t, hundredLink.IsEqual(NewLinkCoin(sdk.NewInt(100))))
	require.True(t, zeroLink.IsEqual(NewLinkCoin(sdk.NewInt(0))))

	require.Equal(t, hundredLink.Add(zeroLink), hundredLink)
	require.Equal(t, hundredLink.Add(hundredLink), NewLinkCoin(sdk.NewInt(200)))

	require.Equal(t, hundredLink.Sub(zeroLink), hundredLink)

}
