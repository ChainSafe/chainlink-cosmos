// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import (
	"bytes"
	"math/big"
)

func (o Observation) GoEthereumValue() *big.Int { return new(big.Int).SetBytes(o.Value) }

func (o Observation) Equal(o2 *Observation) bool {
	return bytes.Equal(o.Value, o2.Value)
}
