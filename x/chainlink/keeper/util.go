// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package keeper

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
)

// feedDataFilter filters the feedData query result by feedId and roundId
func feedDataFilter(requiredFeedID string, requiredRoundID uint64, feedData types.OCRFeedDataInStore) *types.RoundData {
	if feedData.GetRoundId() == requiredRoundID {
		if requiredFeedID == feedData.GetFeedData().GetFeedId() {
			roundData := &types.RoundData{
				FeedId:   feedData.GetFeedData().GetFeedId(),
				FeedData: feedData.GetDeserializedOCRReport(),
			}
			return roundData
		}
		if requiredFeedID == "" {
			roundData := &types.RoundData{
				FeedId:   feedData.GetFeedData().GetFeedId(),
				FeedData: feedData.GetDeserializedOCRReport(),
			}
			return roundData
		}
	}

	return nil
}

func i64tob(val uint64) []byte {
	r := make([]byte, 8)
	for i := uint64(0); i < 8; i++ {
		r[i] = byte((val >> (i * 8)) & 0xff)
	}
	return r
}

func btoi64(val []byte) uint64 {
	r := uint64(0)
	for i := uint64(0); i < 8; i++ {
		r |= uint64(val[i]) << (8 * i)
	}
	return r
}
