// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Validation is the interface to perform any validation against sdk.Msg
type Validation interface {
	Validate(func(msg sdk.Msg) bool) bool
}
