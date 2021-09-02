# List existing keys
chainlinkd keys list --keyring-backend test

aliceAddr=$(chainlinkd keys show alice -a)
alicePK=$(chainlinkd keys show alice -p)

bobAddr=$(chainlinkd keys show bob -a)
bobPK=$(chainlinkd keys show bob -p)

cerloAddr=$(chainlinkd keys show cerlo -a)
cerloPK=$(chainlinkd keys show cerlo -p)

# ======
# Module
# ======

# List all module owner
chainlinkd query chainlink get-module-owner-list --chain-id testchain -o json

# Add new module owner by alice
chainlinkd tx chainlink add-module-owner "$bobAddr" "$bobPK" --from alice --keyring-backend test --chain-id testchain --fees 3link

# Module ownership transfer by bob to alice
chainlinkd tx chainlink module-ownership-transfer "$aliceAddr" "$alicePK" --from bob --keyring-backend test --chain-id testchain --fees 3link

# ====
# Feed
# ====

# Add new feed
chainlinkd tx chainlink add-feed feedid1 "this is test feed 1" "$cerloAddr" 1 2 3 100 "" "$bobAddr,$bobPK,$cerloAddr,$cerloPK" --from alice --keyring-backend test --chain-id testchain --gas 100000 --gas-prices 0.000006link

# Query feed info by feedId
chainlinkd query chainlink get-feed-info feedid1 --chain-id testchain

# Query available feed reward strategy
chainlinkd query chainlink get-feed-reward-avail-strategy --chain-id testchain

# Add feed data provider
chainlinkd tx chainlink add-data-provider feedid1 "$bobAddr" "$bobPK" --from cerlo --keyring-backend test --chain-id testchain --fees 3link

# Query feed info by feedId
chainlinkd query chainlink get-feed-info feedid1 --chain-id testchain

# Remove feed data provider
chainlinkd tx chainlink remove-data-provider feedid1 "$cerloAddr" --from cerlo --keyring-backend test --chain-id testchain --fees 3link

# Query feed info by feedId
chainlinkd query chainlink get-feed-info feedid1 --chain-id testchain

# Feed ownership transfer by cerlo to bob
chainlinkd tx chainlink feed-ownership-transfer feedid1 "$bobAddr" --from cerlo --keyring-backend test --chain-id testchain --fees 3link

# Query feed info by feedId
chainlinkd query chainlink get-feed-info feedid1 --chain-id testchain

# Update submission count parameter
chainlinkd tx chainlink set-submission-count feedid1 100 --from bob --keyring-backend test --chain-id testchain --fees 3link

# Update heartbeat trigger parameter
chainlinkd tx chainlink set-heartbeat-trigger feedid1 200 --from bob --keyring-backend test --chain-id testchain --fees 3link

# Update deviation threshold trigger parameter
chainlinkd tx chainlink set-deviation-threshold-trigger feedid1 300 --from bob --keyring-backend test --chain-id testchain --fees 3link

# Update feed reward parameter
chainlinkd tx chainlink set-feed-reward feedid1 1000 "" --from cerlo --keyring-backend test --chain-id testchain --fees 3link

# ==================
# Feed Data (Report)
# ==================

# Submit feed data by bob
chainlinkd tx chainlink submit-feed-data feedid1 "feed 1 test data" "signatures_bob,signatures_cerlo" "$bobPK,$cerloPK" --from cerlo --keyring-backend test --chain-id testchain --fees 3link

# Query feed data by txHash
chainlinkd query tx C350CAD4673DB75005C6215262633375ECE318BAEDC794820EE43FA958FB8174 --chain-id testchain -o json

# Query feed data by roundId and feedId
chainlinkd query chainlink get-round-feed-data 1 feedid1 --chain-id testchain -o json

# Query feed data by roundId only
chainlinkd query chainlink get-round-feed-data 1 --chain-id testchain -o json

# Query the latest round feed data with feedId
chainlinkd query chainlink get-lates-feed-data feedid2 --chain-id testchain -o json

# Query the latest round of feed data
chainlinkd query chainlink get-latest-feed-data --chain-id testchain -o json
