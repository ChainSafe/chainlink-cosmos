package types

const (
	// ModuleName defines the module name
	ModuleName = "chainlink"

	// StoreKey defines the primary module store key
	FeedStoreKey = ModuleName + "feed"

	// RoundKey defines the secondary module store key
	RoundStoreKey = ModuleName + "round"

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_chainlink"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	/*
		FeedStore key pattern: types.FeedDataKey + feedId + roundId
	*/
	FeedDataKey = "feedData"

	/*
		RoundStore key pattern: types.RoundIdKey + feedId
	*/
	RoundIdKey = "roundId"
)
