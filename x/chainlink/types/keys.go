package types

const (
	// ModuleName defines the module name
	ModuleName = "chainlink"

	// FeedDataStoreKey defines the store key for feed data
	FeedDataStoreKey = ModuleName + "feedData"

	// RoundStoreKey defines the store key for feed roundId
	RoundStoreKey = ModuleName + "round"

	// ModuleOwnerStoreKey defines the store key for module owner
	ModuleOwnerStoreKey = ModuleName + "moduleOwner"

	// FeedInfoStoreKey defines the store key for feed
	FeedInfoStoreKey = ModuleName + "feedInfo"

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
	// FeedDataKey FeedDataStore key pattern: types.FeedDataKey/feedId/roundId
	FeedDataKey = "feedData"

	// RoundIdKey RoundStore key pattern: types.RoundIdKey/feedId
	RoundIdKey = "roundId"

	// ModuleOwnerKey ModuleOwnerStore key pattern: types.ModuleOwnerKey/moduleOwnerAddress
	ModuleOwnerKey = "moduleOwner"

	// FeedInfoKey FeedInfoStore key pattern: types.FeedInfoKey/feedId
	FeedInfoKey = "feed"
)

func GetFeedDataKey(feedId, roundId string) []byte {
	key := FeedDataKey + "/"
	if len(feedId) > 0 {
		key += feedId + "/"
		if len(roundId) > 0 {
			key += roundId
		}
	}
	return KeyPrefix(key)
}

func GetRoundIdKey(feedId string) []byte {
	key := RoundIdKey + "/"
	if len(feedId) > 0 {
		key += feedId
	}
	return KeyPrefix(key)
}

func GetModuleOwnerKey(moduleOwnerAddress string) []byte {
	key := ModuleOwnerKey + "/"
	if len(moduleOwnerAddress) > 0 {
		key += moduleOwnerAddress
	}
	return KeyPrefix(key)
}

func GetFeedInfoKey(feedId string) []byte {
	key := FeedInfoKey + "/"
	if len(feedId) > 0 {
		key += feedId
	}
	return KeyPrefix(key)
}
