// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import (
	"bytes"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/stretchr/testify/require"

	"github.com/gogo/protobuf/proto"
)

func TestTypes_RegisterCodec(t *testing.T) {
	cdc := codec.NewLegacyAmino()
	require.NotEmpty(t, cdc)

	buf := new(bytes.Buffer)
	err := cdc.PrintTypes(buf)
	require.NoError(t, err)
	require.NotEmpty(t, buf.Bytes())

	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/SubmitFeedData")))
	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/AddModuleOwner")))
	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/ModuleOwnershipTransfer")))
	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/AddFeed")))
	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/AddDataProvider")))
	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/RemoveDataProvider")))
	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/SetSubmissionCount")))
	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/SetHeartbeatTrigger")))
	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/SetDeviationThresholdTrigger")))
	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/SetFeedReward")))
	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/FeedOwnershipTransfer")))
	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/AddAccount")))
	require.False(t, bytes.Contains(buf.Bytes(), []byte("chainlink/EditAccount")))

	RegisterCodec(cdc)

	err = cdc.PrintTypes(buf)
	require.NoError(t, err)
	require.NotEmpty(t, buf.Bytes())

	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/SubmitFeedData")))
	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/AddModuleOwner")))
	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/ModuleOwnershipTransfer")))
	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/AddFeed")))
	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/AddDataProvider")))
	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/RemoveDataProvider")))
	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/SetSubmissionCount")))
	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/SetHeartbeatTrigger")))
	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/SetDeviationThresholdTrigger")))
	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/SetFeedReward")))
	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/FeedOwnershipTransfer")))
	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/AddAccount")))
	require.True(t, bytes.Contains(buf.Bytes(), []byte("chainlink/EditAccount")))
}

func TestTypes_RegisterInterfaces(t *testing.T) {
	nir := types.NewInterfaceRegistry()

	_, e := nir.Resolve("/" + proto.MessageName(&MsgFeedData{}))
	require.Error(t, e)

	RegisterInterfaces(nir)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgFeedData{}))
	require.NoError(t, e)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgModuleOwner{}))
	require.NoError(t, e)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgModuleOwnershipTransfer{}))
	require.NoError(t, e)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgFeed{}))
	require.NoError(t, e)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgAddDataProvider{}))
	require.NoError(t, e)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgRemoveDataProvider{}))
	require.NoError(t, e)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgSetSubmissionCount{}))
	require.NoError(t, e)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgSetHeartbeatTrigger{}))
	require.NoError(t, e)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgSetDeviationThresholdTrigger{}))
	require.NoError(t, e)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgSetFeedReward{}))
	require.NoError(t, e)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgFeedOwnershipTransfer{}))
	require.NoError(t, e)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgAccount{}))
	require.NoError(t, e)

	_, e = nir.Resolve("/" + proto.MessageName(&MsgEditAccount{}))
	require.NoError(t, e)
}
