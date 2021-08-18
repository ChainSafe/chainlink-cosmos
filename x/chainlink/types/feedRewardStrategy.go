package types

// RewardPayout describes the calculated reward that data provider gets after the selected strategy applied
type RewardPayout struct {
	DataProvider *DataProvider `json:"dataProvider"`
	Amount       uint32        `json:"amount"`
}

type FeedRewardStrategyFunc func(*MsgFeed, *MsgFeedData) ([]RewardPayout, error)

var FeedRewardStrategyConvertor = map[string]FeedRewardStrategyFunc{}

// NewFeedRewardStrategyRegister registers the reward calculation strategies when the chain launches
func NewFeedRewardStrategyRegister(feedRewardStrategyFns map[string]FeedRewardStrategyFunc) {
	if feedRewardStrategyFns == nil {
		return
	}

	for name := range feedRewardStrategyFns {
		if name == "" {
			panic("feed reward strategy name can not be empty")
		}
	}

	FeedRewardStrategyConvertor = feedRewardStrategyFns
}
