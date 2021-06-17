package keeper

import "github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"

// feedDataFilter filters the feedData query result by feedId and roundId
func feedDataFilter(requiredFeedID string, requiredRoundID uint64, feedData types.OCRFeedDataInStore) []*types.RoundData {
	feedRoundData := make([]*types.RoundData, 0)

	if feedData.GetRoundId() == requiredRoundID && requiredFeedID == feedData.GetFeedData().GetFeedId() {
		roundData := types.RoundData{
			FeedId:   feedData.GetFeedData().GetFeedId(),
			FeedData: feedData.GetDeserializedOCRReport(),
		}
		feedRoundData = append(feedRoundData, &roundData)
	}
	if feedData.GetRoundId() == requiredRoundID && requiredFeedID == "" {
		roundData := types.RoundData{
			FeedId:   feedData.GetFeedData().GetFeedId(),
			FeedData: feedData.GetDeserializedOCRReport(),
		}
		feedRoundData = append(feedRoundData, &roundData)
	}

	return feedRoundData
}
