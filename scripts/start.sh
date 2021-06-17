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
#chainlinkd tx chainlink submit-feedData "testfeedid1" "feed 1 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain

# Query feed data by txHash
#chainlinkd query tx FDD655DE6F4E1B0F7A04163F856A88E4BACAC9755402B90F77D9EF9F45570168 --chain-id testchain -o json

# Query feed data by getRoundFeedData with pagination
#chainlinkd query chainlink getRoundFeedData 1 "testfeedid1" --chain-id testchain -o json

# Query feed data by getLatestFeedData
#chainlinkd query chainlink getLatestFeedData "testfeedid1" --chain-id testchain -o json

