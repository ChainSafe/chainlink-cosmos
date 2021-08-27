// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package types

import (
	"encoding/binary"
)

// DomainSeparationTag consists of:
// 11-byte zero padding
// 16-byte configDigest
// 4-byte epoch
// 1-byte round
// It uniquely identifies a message to a particular group-epoch-round tuple.
// It is used in signature verification
type DomainSeparationTag [32]byte

func (r ReportContext) DomainSeparationTag() (d DomainSeparationTag) {
	serialization := r.ConfigDigest.Value[:]
	serialization = append(serialization, []byte{0, 0, 0, 0}...)
	binary.BigEndian.PutUint32(serialization[len(serialization)-4:], r.Epoch)
	serialization = append(serialization, byte(r.Round))
	copy(d[11:], serialization)
	return d
}

func (r ReportContext) Equal(r2 ReportContext) bool {
	return r.ConfigDigest == r2.ConfigDigest && r.Epoch == r2.Epoch && r.Round == r2.Round
}
