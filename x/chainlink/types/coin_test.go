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
	require.Equal(t, NewLinkCoin(sdk.NewInt(200)).Sub(hundredLink), hundredLink)
}

func TestTypes_NewLinkDecCoin(t *testing.T) {
	hundredAmount := sdk.NewInt(100)
	hundredDecLink := NewLinkDecCoin(hundredAmount)
	require.NoError(t, hundredDecLink.Validate())
	require.True(t, hundredDecLink.IsValid())

	zeroAmount := sdk.NewInt(0)
	zeroDecLink := NewLinkDecCoin(zeroAmount)
	require.True(t, zeroDecLink.IsZero())
	require.True(t, zeroDecLink.IsValid())
	require.NoError(t, zeroDecLink.Validate())

	require.Equal(t, hundredDecLink.GetDenom(), LinkDenom)

	require.Equal(t, hundredDecLink.Amount, sdk.NewDec(hundredAmount.Int64()))
	require.Equal(t, zeroDecLink.Amount, sdk.NewDec(zeroAmount.Int64()))

	require.Equal(t, hundredDecLink.Add(zeroDecLink), hundredDecLink)
	require.Equal(t, hundredDecLink.Add(hundredDecLink), NewLinkDecCoin(sdk.NewInt(200)))

	require.Equal(t, hundredDecLink.Sub(zeroDecLink), hundredDecLink)
	require.Equal(t, NewLinkDecCoin(sdk.NewInt(200)).Sub(hundredDecLink), hundredDecLink)

	require.True(t, hundredDecLink.IsGTE(zeroDecLink))
	require.True(t, zeroDecLink.IsLT(hundredDecLink))

	require.True(t, hundredDecLink.IsEqual(NewLinkDecCoin(sdk.NewInt(100))))
	require.True(t, zeroDecLink.IsEqual(NewLinkDecCoin(sdk.NewInt(0))))
}

func TestTypes_NewLinkCoin64(t *testing.T) {
	hundredAmount := int64(100)
	hundredLink := NewLinkCoinInt64(hundredAmount)
	require.NoError(t, hundredLink.Validate())
	require.True(t, hundredLink.IsValid())
	require.Equal(t, hundredLink, NewLinkCoin(sdk.NewInt(100)))

	zeroAmount := int64(0)
	zeroLink := NewLinkCoinInt64(zeroAmount)
	require.True(t, zeroLink.IsZero())
	require.True(t, zeroLink.IsValid())
	require.NoError(t, zeroLink.Validate())
	require.Equal(t, zeroLink, NewLinkCoin(sdk.NewInt(0)))
}
