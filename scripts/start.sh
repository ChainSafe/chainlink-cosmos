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

# Add init chainlink module owner to the genesis file
#chainlinkd tx chainlink add-genesis-module-owner "cosmos17sxrwr30r8nlf7ahykscpkpu0muqy8k0hegadc" "cosmospub1addwnpepq0z2excjzlqn3l6g6eh6kj56h6rt38f5yfu7aeqxdy6eg7dn0a8wy5alhhx" --keyring-backend test --chain-id testchain
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

# invalid query with incorrect FeedId
#chainlinkd query chainlink getRoundFeedData 999 "testfeedid1" --chain-id testchain -o json

# Query the latest round feed data with feedId
#chainlinkd query chainlink getLatestFeedData "testfeedid1" --chain-id testchain -o json

# Query the latest round of feed data
#chainlinkd query chainlink getLatestFeedData --chain-id testchain -o json

# List all module owner
#chainlinkd query chainlink getModuleOwnerList --chain-id testchain -o json

# List existing keys
#chainlinkd keys list --keyring-backend test

# Add new module owner
# chainlinkd tx chainlink addModuleOwner "cosmos1tq9q6sfcfzkj5vq90la2nwfhz974whm8x9jl9s" "cosmospub1addwnpepqgr9xvm3s5ks9naq00ghvatww338q6p7jr4apg64vds0auva3vk4c3ddpl7" --from alice --keyring-backend test --chain-id testchain

