# Submit feed data
chainlinkd tx chainlink submitFeedData "testfeedid1" "feed 1 test data" "dummy signatures" --from "cosmos15pfql3pfx4z0v7vgynj5fnmfh237cq5jn6vlz8" --keyring-backend test --chain-id testchain

# Query feed data by txHash
chainlinkd query tx A0B849C7A5ABB51B3FA9DC723A6C1CB8C4B6C255DB98D0EC0FD3DCD04316E387 --chain-id testchain -o json

# Query feed data by roundId and feedId
chainlinkd query chainlink getRoundFeedData 1 "testfeedid1" --chain-id testchain -o json

# Query feed data by roundId only
chainlinkd query chainlink getRoundFeedData 2 --chain-id testchain -o json

# Query the latest round feed data with feedId
chainlinkd query chainlink getLatestFeedData "testfeedid1" --chain-id testchain -o json

# Query the latest round of feed data
chainlinkd query chainlink getLatestFeedData --chain-id testchain -o json

# List all module owner
chainlinkd query chainlink getModuleOwnerList --chain-id testchain -o json

# List existing keys
chainlinkd keys list --keyring-backend test

# Add new module owner
chainlinkd tx chainlink addModuleOwner "cosmos12uxqzq7aae5ew2rksy0cl0cua7e75cu5t3rjxf" "cosmospub1addwnpepqw2xec9wutvxdhvgke029v3hx97jdpzqpr6d33jx62drd92qjxxrsm9t7q5" --from cosmos1mxea03k4y8n7c4fyeseu3c889ejw7kncwhp0uv --keyring-backend test --chain-id testchain

