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
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgFeedData{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgModuleOwner{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgModuleOwnershipTransfer{})
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
