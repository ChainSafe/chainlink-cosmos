#!/bin/bash

# Clean up
rm -rf ~/.chainlinkd/

# Add two accounts
chainlinkd keys add alice --keyring-backend test
chainlinkd keys add bob --keyring-backend test

# Initialize configuration files and genesis file
chainlinkd init testchain --chain-id testchain

# Add both accounts, with coins to the genesis file
chainlinkd add-genesis-account $(chainlinkd keys show alice -a --keyring-backend test) 1000token,100000000stake
chainlinkd add-genesis-account $(chainlinkd keys show bob -a --keyring-backend test) 1000token,100000000stake

# Add init chainLink module owner to the genesis file
chainlinkd tx chainlink add-genesis-module-owner $(chainlinkd keys show alice -a --keyring-backend test) $(chainlinkd keys show alice -p --keyring-backend test) --keyring-backend test --chain-id testchain

# Generate the gen tx that creates a validator with a self-delegation,
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
# chainlinkd tx chainlink submitFeedData "testfeedid1" "feed 1 test data" "dummy signatures" --from alice --keyring-backend test --chain-id testchain

# Query feed data by txHash
#chainlinkd query tx A0B849C7A5ABB51B3FA9DC723A6C1CB8C4B6C255DB98D0EC0FD3DCD04316E387 --chain-id testchain -o json

# Query feed data by roundId and feedId
#chainlinkd query chainlink getRoundFeedData 1 "testfeedid1" --chain-id testchain -o json

# Query feed data by roundId only
#chainlinkd query chainlink getRoundFeedData 2 --chain-id testchain -o json

# Query the latest round feed data with feedId
#chainlinkd query chainlink getLatestFeedData "testfeedid1" --chain-id testchain -o json

# Query the latest round of feed data
#chainlinkd query chainlink getLatestFeedData --chain-id testchain -o json

# List all module owner
#chainlinkd query chainlink getModuleOwnerList --chain-id testchain -o json

# List existing keys
#chainlinkd keys list --keyring-backend test

# Add new module owner
# chainlinkd tx chainlink addModuleOwner "cosmos1stdn5v0tcdc6s2vy79rk8yxujlahw4jydyntt9" "cosmospub1addwnpepqfh6a0rsu9m8q5tfqkp97whexdc4jtdgnel7xvf2hv26c6m860e3gt4tf9u" --from bob --keyring-backend test --chain-id testchain

