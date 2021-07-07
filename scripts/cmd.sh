# module owner
# List existing keys
chainlinkd keys list --keyring-backend test

# List all module owner
chainlinkd query chainlink getModuleOwnerList --chain-id testchain -o json

# Add new module owner by alice
chainlinkd tx chainlink addModuleOwner "cosmos1sw9p3pl9237rz25dreg9qdqyum6scjgyzj4vah" "cosmospub1addwnpepqwfed02ccmkv6feyl93ym3n62eml6xvsttzr8klc3fnvgcgau2ff2p3q5tc" --from alice --keyring-backend test --chain-id testchain

# module ownership transfer by bob
chainlinkd tx chainlink moduleOwnershipTransfer "cosmos1wxzkyuqnte8z6m0vt7g4r3j9z6rj2v7mclj92k" "cosmospub1addwnpepq03r94dzyvw70rff4kc72h90kra7vu6yurq2l5cfggtevsgrzdem598g989" --from bob --keyring-backend test --chain-id testchain

# feed
# Add new feed
chainlinkd tx chainlink addFeed feedid1 cosmos199a3lwv0w2amta8h4pp96jr0mm4f7ssa78jyft 1 2 3 cosmos199a3lwv0w2amta8h4pp96jr0mm4f7ssa78jyft,cosmospub1addwnpepqwt7anmkmmvw7sw9a2uex520munmq4yxuyj5nyyscz2uwddc37us2qkcvt8 --from alice --keyring-backend test --chain-id testchain

# Query feed by feedId
chainlinkd query chainlink getFeed "feedid1" --chain-id testchain

# feed data (report)
# Submit feed data by alice
chainlinkd tx chainlink submitFeedData "feedid1" "feed 1 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain

# Query feed data by txHash
chainlinkd query tx C350CAD4673DB75005C6215262633375ECE318BAEDC794820EE43FA958FB8174 --chain-id testchain -o json

# Query feed data by roundId and feedId
chainlinkd query chainlink getRoundFeedData 1 "feedid1" --chain-id testchain -o json

# Query feed data by roundId only
chainlinkd query chainlink getRoundFeedData 1 --chain-id testchain -o json

# Query the latest round feed data with feedId
chainlinkd query chainlink getLatestFeedData "feedid1" --chain-id testchain -o json

# Query the latest round of feed data
chainlinkd query chainlink getLatestFeedData --chain-id testchain -o json

