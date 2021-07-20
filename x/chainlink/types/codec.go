// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
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
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgFeedData{},
		&MsgModuleOwner{},
		&MsgModuleOwnershipTransfer{},
		&MsgFeed{},
		&MsgAddDataProvider{},
		&MsgRemoveDataProvider{},
		&MsgSetSubmissionCount{},
		&MsgSetHeartbeatTrigger{},
		&MsgSetDeviationThresholdTrigger{},
		&MsgFeedOwnershipTransfer{},
	)

	/*registry.RegisterInterface(
		"chainlink.v1beta.RoundDataI",
		(*exported.RoundDataI)(nil),
		&RoundData{},
	)

	registry.RegisterInterface(
		"chainlink.v1beta.ObservationI",
		(*exported.ObservationI)(nil),
		&Observation{},
	)

	registry.RegisterInterface(
		"chainlink.OCRAbiEncodedI.",
		(*exported.OCRAbiEncodedI)(nil),
		&OCRAbiEncoded{},
	)*/

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
