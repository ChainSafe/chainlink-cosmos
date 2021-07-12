# Cosmos ChainLink Module Documentation

Cosmos ChainLink Module is a Cosmos SDK module that allows developers to add Chainlink data feed support to their applications.

## Basic Concept

ChainLink Module allows developers to add Chainlink data feed support to their applications with a permission control. 
There are three level account permission to manage the chainlink feed data:

1. module owner
2. feed owner
3. feed data provider

Module owner is a list of cosmos accounts to manage the chainlink module; however，an init module owner assigned in the genesis.json. Module owners are trusted 
cosmos accounts which has ability to create new data feed with init feed parameters and owner of the feed. Module ownership is able to be transferred.

Feed owner is a cosmos account and only one owner per a feed. The init owner of a feed is assigned by module owner when the feed is created. Feed owner is able to manage the feed parameters such as 
heartbeat, feed data submission count, data providers etc. Feed ownership is able to be transferred.

Feed data provider is also a cosmos account that is able to sign the feed data submit transaction to the module. Only the valid data provider of a feed is able to submit the fee data
to that feed, valid feed data provider list is managed by the feed owner.

Currently the module is in development, all the transactions and queries are available throught CLI, we will be working on REST, JSON-RPC 2.0 and gRPC endpoints later.

All CLI commands are using Cosmos CLI command format. For example get all module owner CLI: `chainlinkd query chainlink getModuleOwnerList --chain-id testchain -o json`

## CLI Endpoints

### Module owner

####Transaction

1. Add init module owner in genesis file  
The init module owner address and public key are required.  
Note: this transaction is only available through CLI, and can only be executed before chain launching.
```bash
add-genesis-module-owner [address] [pubKey]
```

2. Add new module owner  
Can be signed by existing module owner only.
```bash
addModuleOwner [address] [pubKey]
```

3. Module ownership transfer  
Can be signed by existing module owner only.   
   The address and pubKey should match the new module owner.
```bash
moduleOwnershipTransfer [address] [pubKey]
```

4. Add new feed  
   Can be signed by existing module owner only.  
   `initDataProviderList` is a string with data providers' address and pubKey connecting with comma.   
   For example:`address1,keyKey1,address2,pubKey2`
```bash
addFeed [feedId] [feedOwnerAddress] [submissionCount] [heartbeatTrigger] [deviationThresholdTrigger] [initDataProviderList]
```

#### Query

1. Get all module owners
```bash
getModuleOwnerList
```

### Feed Owner

####Transaction
WIP

#### Query

1. Get feed info by feedId
```bash
getFeedInfo [feedId]
```

### Feed Data Provider

#### Transaction

1. Submit feed data  
Only valid data provider(signer of this transaction) is able to submit feed data to particular feed base on feedId.
```bash
submitFeedData [feedId] [feedData] [signatures]
```

#### Query 
1. Query feed data by round  
`feedId` is optional
```bash
getRoundFeedData [roundId] [feedId]
```

2. Query the latest round of feed data  
`feedId` is optional
```bash
getLatestFeedData [feedId]
```






