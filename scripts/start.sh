#!/bin/bash

# Clean up
rm -rf ~/.chainlinkd/

# Build
make install

# Initialize configuration files and genesis file
chainlinkd init testchain --chain-id testchain

# Add two accounts
chainlinkd keys add alice --keyring-backend test
chainlinkd keys add bob --keyring-backend test

# Add both accounts, with coins to the genesis file
chainlinkd add-genesis-account $(chainlinkd keys show alice -a --keyring-backend test) 1000token,100000000stake
chainlinkd add-genesis-account $(chainlinkd keys show bob -a --keyring-backend test) 1000token,100000000stake

# Generate the gen tx
chainlinkd gentx alice 100000000stake --amount=100000000stake --keyring-backend test --chain-id testchain

# Input the genTx into the genesis file, so that the chain is aware of the validators
chainlinkd collect-gentxs

# Make sure your genesis file is correct.
chainlinkd validate-genesis

# Replace app.toml API config to enable rest API server
perl -0777 -i.original -pe 's/API server should be enabled.\nenable = false/API server should be enabled.\nenable = true/igs' ~/.chainlinkd/config/app.toml

# Start chain
chainlinkd start

# Submit feed data
# chainlinkd tx chainlink submit-feedData "testfeedid1" "feed 1 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain

# Query feed data by txHash
#chainlinkd query tx AD1CEB561E3225D26E682918FDFFD4B507B0FE15200072AB3EC2C40380280B8F --chain-id testchain -o json

# Query feed data by roundId and feedId
#chainlinkd query chainlink getRoundFeedData 1 "testfeedid1" --chain-id testchain -o json

# Query feed data by roundId only
#chainlinkd query chainlink getRoundFeedData 2 --chain-id testchain -o json

# invalid query with incorrect FeedId
#chainlinkd query chainlink getRoundFeedData 999 "testfeedid1" --chain-id testchain -o json

# Query the latest round feed data with feedId
#chainlinkd query chainlink getLatestFeedData "testfeedid1" --chain-id testchain -o json

# Query the latest round of feed data
#chainlinkd query chainlink getLatestFeedData --chain-id testchain -o json

# chainlinkd tx chainlink submit-feedData "testfeedid1" "feed 1 test data1" "dummy signatures" --from alice --keyring-backend test --chain-id testchain
# chainlinkd tx chainlink submit-feedData "testfeedid1" "feed 1 test data2" "dummy signatures" --from alice --keyring-backend test --chain-id testchain

# chainlinkd tx chainlink submit-feedData "testfeedid2" "feed 2 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain

# chainlinkd query chainlink getLatestFeedData "testfeedid1" --chain-id testchain -o json

# chainlinkd query chainlink getLatestFeedData "testfeedid2" --chain-id testchain -o json