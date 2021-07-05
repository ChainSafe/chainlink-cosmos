package types

const (
	// ModuleName defines the module name
	ModuleName = "chainlink"

	// FeedDataStoreKey defines the primary module store key for feed data
	FeedDataStoreKey = ModuleName + "feedData"

	// RoundStoreKey defines the secondary module store key for feed roundId
	RoundStoreKey = ModuleName + "round"

	// ModuleOwnerStoreKey defines the module owner scope store key
	ModuleOwnerStoreKey = ModuleName + "moduleOwner"

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
	// FeedDataStore key pattern: types.FeedDataKey/feedId/roundId
	FeedDataKey = "feedData"

	// RoundStore key pattern: types.RoundIdKey/feedId
	RoundIdKey = "roundId"

	// ModuleOwnerStore key pattern: types.ModuleOwnerKey/moduleOwnerAddress
	ModuleOwnerKey = "moduleOwner"
)
