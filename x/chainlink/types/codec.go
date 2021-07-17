// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgFeedData{}, "chainlink/SubmitFeedData", nil)
	cdc.RegisterConcrete(MsgModuleOwner{}, "chainlink/AddModuleOwner", nil)
	cdc.RegisterConcrete(MsgModuleOwnershipTransfer{}, "chainlink/ModuleOwnershipTransfer", nil)
	cdc.RegisterConcrete(MsgFeed{}, "chainlink/AddFeed", nil)
	cdc.RegisterConcrete(MsgAddDataProvider{}, "chainlink/AddDataProvider", nil)
	cdc.RegisterConcrete(MsgRemoveDataProvider{}, "chainlink/RemoveDataProvider", nil)
	cdc.RegisterConcrete(MsgSetSubmissionCount{}, "chainlink/SetSubmissionCount", nil)
	cdc.RegisterConcrete(MsgSetHeartbeatTrigger{}, "chainlink/SetHeartbeatTrigger", nil)
	cdc.RegisterConcrete(MsgSetDeviationThresholdTrigger{}, "chainlink/SetDeviationThresholdTrigger", nil)
	cdc.RegisterConcrete(MsgFeedOwnershipTransfer{}, "chainlink/FeedOwnershipTransfer", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgFeedData{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgModuleOwner{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgModuleOwnershipTransfer{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgFeed{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgAddDataProvider{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRemoveDataProvider{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSetSubmissionCount{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSetHeartbeatTrigger{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSetDeviationThresholdTrigger{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgFeedOwnershipTransfer{})
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
