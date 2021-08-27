// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import "fmt"

func (c ConfigDigest) Hex() string {
	return fmt.Sprintf("%x", c.Value)
}
