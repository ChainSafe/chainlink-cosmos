# Submit feed data by alice
chainlinkd tx chainlink submitFeedData "testfeedid1" "feed 1 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain

# Query feed data by txHash
chainlinkd query tx C350CAD4673DB75005C6215262633375ECE318BAEDC794820EE43FA958FB8174 --chain-id testchain -o json

# Query feed data by roundId and feedId
chainlinkd query chainlink getRoundFeedData 1 "testfeedid1" --chain-id testchain -o json

# Query feed data by roundId only
chainlinkd query chainlink getRoundFeedData 2 --chain-id testchain -o json

# Query the latest round feed data with feedId
chainlinkd query chainlink getLatestFeedData "testfeedid1" --chain-id testchain -o json

# Query the latest round of feed data
chainlinkd query chainlink getLatestFeedData --chain-id testchain -o json

# List existing keys
chainlinkd keys list --keyring-backend test

# List all module owner
chainlinkd query chainlink getModuleOwnerList --chain-id testchain -o json

# Add new module owner by alice
chainlinkd tx chainlink addModuleOwner "cosmos1wjthz4kmkcusava94f55pg06cqrlxm889udgjn" "cosmospub1addwnpepq00vr93dx4k88rpfupm5wxv50nmq69chxlfal279cexxjy8yl29dc2kqqqn" --from alice --keyring-backend test --chain-id testchain

# module ownership transfer by bob
chainlinkd tx chainlink moduleOwnershipTransfer "cosmos1z6r57d75mw4yzenykk4l9zjma0mjusseaz5yk3" "cosmospub1addwnpepq2uxljzkshf02yuk9k7ehmmcru0de5p5d9gw54g49jgq7djeufq32we4zr0" --from bob --keyring-backend test --chain-id testchain

chainlinkd tx chainlink addFeed "feedid2" "cosmos1avfmvhzf2kng33v0365jgv6qc3a8823ygddwr0" "1" "2" "3" "cosmos1avfmvhzf2kng33v0365jgv6qc3a8823ygddwr0 cosmospub1addwnpepqv6cn6jvja4epx42wyxvdd6lrn43p6ldskpawtqafksvksgl7p0puk8nwtu" --from alice --keyring-backend test --chain-id testchain

chainlinkd query chainlink getFeed "feedid1" --chain-id testchain
