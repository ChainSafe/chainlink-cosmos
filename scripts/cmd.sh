# module owner
# List existing keys
chainlinkd keys list --keyring-backend test

# List all module owner
chainlinkd query chainlink getModuleOwnerList --chain-id testchain -o json

# Add new module owner by alice
chainlinkd tx chainlink addModuleOwner cosmos1p68ydzcyq6khyyz8up8l4pl56lzqentfguytnu cosmospub1addwnpepqwd9c7dtj8er34j4wfjc8hf50nzakgcx04tmdd0qf42ryxl85p665rvmpmy --from alice --keyring-backend test --chain-id testchain

# module ownership transfer by bob
chainlinkd tx chainlink moduleOwnershipTransfer cosmos1wxzkyuqnte8z6m0vt7g4r3j9z6rj2v7mclj92k cosmospub1addwnpepq03r94dzyvw70rff4kc72h90kra7vu6yurq2l5cfggtevsgrzdem598g989 --from bob --keyring-backend test --chain-id testchain

# feed
# Add new feed
chainlinkd tx chainlink addFeed feedid1 cosmos1j8t7v6tt98wjhzhcuwjqmnqzgaz4v8uffhd4gq 1 2 3 cosmos1j8t7v6tt98wjhzhcuwjqmnqzgaz4v8uffhd4gq,cosmospub1addwnpepq2cxc37a5kwle7rhtj0qvhlx8nhrujvf2r6h66vhx2leakl2wpn2qnj0j8m --from alice --keyring-backend test --chain-id testchain

# Query feed info by feedId
chainlinkd query chainlink getFeedInfo feedid1 --chain-id testchain

# feed data (report)
# Submit feed data by alice
chainlinkd tx chainlink submitFeedData feedid1 "feed 1 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain

# Query feed data by txHash
chainlinkd query tx C350CAD4673DB75005C6215262633375ECE318BAEDC794820EE43FA958FB8174 --chain-id testchain -o json

# Query feed data by roundId and feedId
chainlinkd query chainlink getRoundFeedData 1 feedid1 --chain-id testchain -o json

# Query feed data by roundId only
chainlinkd query chainlink getRoundFeedData 1 --chain-id testchain -o json

# Query the latest round feed data with feedId
chainlinkd query chainlink getLatestFeedData feedid2 --chain-id testchain -o json

# Query the latest round of feed data
chainlinkd query chainlink getLatestFeedData --chain-id testchain -o json

