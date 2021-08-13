package types

type RewardPayout struct {
	DataProvider *DataProvider `json:"dataProvider"`
	Amount       uint32        `json:"amount"`
}

type feedRewardStrategyFunc func(*MsgFeed, *MsgFeedData) []RewardPayout

var feedRewardStrategyConvertor = map[string]feedRewardStrategyFunc{}

func NewFeedRewardStrategyRegister(feedRewardStrategyFuncs map[string]feedRewardStrategyFunc) {
	if feedRewardStrategyFuncs == nil {
		return
	}
	feedRewardStrategyConvertor = feedRewardStrategyFuncs
}
