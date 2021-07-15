# module owner
# List existing keys
chainlinkd keys list --keyring-backend test

aliceAddr=$(chainlinkd keys show alice -a)
alicePK=$(chainlinkd keys show alice -p)

bobAddr=$(chainlinkd keys show bob -a)
bobPK=$(chainlinkd keys show bob -p)

cerloAddr=$(chainlinkd keys show cerlo -a)
cerloPK=$(chainlinkd keys show cerlo -p)

# List all module owner
chainlinkd query chainlink getModuleOwnerList --chain-id testchain -o json

# Add new module owner by alice
chainlinkd tx chainlink addModuleOwner "$bobAddr" "$bobPK" --from alice --keyring-backend test --chain-id testchain

# module ownership transfer by bob
chainlinkd tx chainlink moduleOwnershipTransfer "$aliceAddr" "$alicePK" --from bob --keyring-backend test --chain-id testchain

# feed
# Add new feed
chainlinkd tx chainlink addFeed feedid1 "$cerloAddr" 1 2 3 4 "$cerloAddr,$cerloPK" --from alice --keyring-backend test --chain-id testchain

# Query feed info by feedId
chainlinkd query chainlink getFeedInfo feedid1 --chain-id testchain

# Add feed data provider
chainlinkd tx chainlink addDataProvider feedid1 "$bobAddr" "$bobPK" --from alice --keyring-backend test --chain-id testchain

# Query feed info by feedId
chainlinkd query chainlink getFeedInfo feedid1 --chain-id testchain

# Remove feed data provider
chainlinkd tx chainlink removeDataProvider feedid1 "$cerloAddr" --from alice --keyring-backend test --chain-id testchain

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
