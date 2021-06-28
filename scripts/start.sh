#!/bin/bash

# Clean up
rm -rf ~/.chainlinkd/

# Add two accounts
chainlinkd keys add alice --keyring-backend test
chainlinkd keys add bob --keyring-backend test
chainlinkd keys add cerlo --keyring-backend test

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