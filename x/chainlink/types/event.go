package types

import (
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
)

func EmitEvent(e proto.Message, manager *types.EventManager) error {
	err := manager.EmitTypedEvent(e)
	if err != nil {
		return err
	}
	return nil
}
