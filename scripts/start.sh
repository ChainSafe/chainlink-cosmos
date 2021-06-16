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
#chainlinkd tx chainlink submit-feedData "testfeedID1" "this is test feed data" --from alice --keyring-backend test --chain-id testchain

# Query feed data by txHash
#chainlinkd query tx D357F65DD6AFF2292D0BB66F8B77D94CB1ECAC0D9E0A5150C05277456256F10C --chain-id testchain -o json

# Query feed data by feedID with pagination
#chainlinkd query chainlink list-feedData "testfeedID1" --chain-id testchain -o json

