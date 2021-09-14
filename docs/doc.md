# Cosmos ChainLink Module Documentation

Cosmos ChainLink Module is a Cosmos SDK module that allows developers to add Chainlink data feed support to their
applications.

## Basic Concept

ChainLink Module allows developers to add Chainlink data feed support to their applications with a permission control.
There are three level account permission to manage the chainlink feed data:

1. Module owner
2. Feed owner
3. Feed data provider

Module owner is a list of cosmos accounts to manage the chainlink module; however, an init module owner need to be
assigned in the genesis.json. Module owners are trusted cosmos accounts which have ability to create new data feed with
init feed parameters and owner of the feed. Module ownership is able to be transferred.

Feed owner is a cosmos account and only one owner per a feed. The init owner of a feed is assigned by module owner when
the feed is created. Feed owner is able to manage the feed parameters such as heartbeat, feed data submission count,
valid data provider set etc. Feed ownership is also able to be transferred.

Feed data provider is also a cosmos account that is able to sign and broadcast the feed data submit transaction to the
module. Only the valid data provider of a feed is able to submit the feed data to that feed, valid feed data provider
list is managed by the feed owner.

Currently, the module is in development, all the transactions and queries are available throught CLI, we will be working
on REST, JSON-RPC 2.0 and gRPC endpoints later.

All CLI commands are using Cosmos CLI command format. For example get all module owner CLI:

`chainlinkd query chainlink getModuleOwnerList --chain-id testchain -o json`

## CLI Endpoints

### Module owner

#### Transaction

1. Add init module owner in genesis file  
   The init module owner address and public key are required.  
   Note: this transaction is only available through CLI, and can only be executed before chain launching.

```bash
add-genesis-module-owner [address] [pubKey]
```

2. Add new module owner  
   Can be signed by existing module owner only.

```bash
add-module-owner [address] [pubKey]
```

3. Module ownership transfer  
   Can be signed by existing module owner only.   
   The address and pubKey should match the new module owner.

```bash
module-ownership-transfer [address] [pubKey]
```

4. Add new feed  
   Can be signed by existing module owner only.  
   `initDataProviderList` is a string with data providers' address and pubKey connecting with comma.   
   For example:`address1,keyKey1,address2,pubKey2`

```bash
add-feed [feedId] [feedOwnerAddress] [submissionCount] [heartbeatTrigger] [deviationThresholdTrigger] [baseFeedRewardAmount] [feedRewardStrategy] [initDataProviderList]
```

#### Query

1. Get all current module owners

```bash
get-module-owner-list
```

### Feed Owner

#### Transaction

1. Add new data provider to a feed    
   Can be signed by feed owner only.  
   `address` is the address of new data provider  
   `publicKey` is the public key of new data provider  
   `address` and `publicKey` must match

```bash
add-data-provider [feedId] [address] [publicKey]
```

2. Remove data provider from a feed    
   Can be signed by feed owner only.  
   `address` is the address of new data provider  
   `publicKey` is the public key of new data provider  
   `address` and `publicKey` must match

```bash
remove-data-provider [feedId] [address] [publicKey]
```

3. Set a new submission count of a feed   
   Can be signed by feed owner only.    
   `count` is the number of valid signatures required

```bash
set-submission-count [feedId] [count]
```

4. Set a new heart beat trigger of a feed   
   Can be signed by feed owner only.  
   `heartbeatTrigger` is a number of milliseconds

```bash
set-heartbeat-trigger [feedId] [heartbeatTrigger]
```

5. Set a new deviation threshold of a feed   
   Can be signed by feed owner only.  
   `deviationThresholdTrigger` is the deviation threshold expressed as thousandths of a percent.  
   For example if the price of `ATOM/USD` changes by 1% then a new round should occur even if the heartbeat interval has
   not elapsed.

```bash
set-deviation-threshold-trigger [feedId] [deviationThresholdTrigger]
```

6. Set a new data provider reward schema of a feed  
   Can be signed by feed owner only.  
   Currently, the feedReward is a number, the complex reward schema will be enabled later.  
   `baseFeedRewardAmount` is the base amount of app native token given to the valid data provider for each round as reward.
   `feedRewardStrategy` is the strategy name in effect.

```bash
set-feed-reward [feedId] [baseFeedRewardAmount] [feedRewardStrategy]
```

7. Feed ownership transfer  
   Can be signed by feed owner only.

```bash
feed-ownership-transfer [feedId] [newFeedOwnerAddress]
```

#### Query

1. Get feed info by feedId

```bash
get-feed-info [feedId]
```

2. Get available feed reward payout strategy list

```bash
get-feed-reward-avail-strategy
```

### Feed Data Provider

#### Transaction

1. Submit feed data  
   Only valid data provider(signer of this transaction) is able to submit feed data to particular feed base on feedId.

```bash
submit-feed-data [feedId] [feedData] [signatures] [cosmosPubKeys]
```

#### Query

1. Query feed data by round  
   `feedId` is optional

```bash
get-round-feed-data [roundId] [feedId]
```

2. Query the latest round of feed data  
   `feedId` is optional

```bash
get-latest-feed-data [feedId]
```

## Configurable Transaction Data Validation Interface

Configurable transaction data validation Interface gives the possibility to the app level devs to implement the
customizable logic of the tx data validation of any ChainLink Cosmos Module transactions against any external resources.
Currently, this interface is applicable for `Submit Feed Data` transaction only. Other tx support is WIP.

App devs can implement a func that takes a `sdk.Msg` as input and return a single boolean as the output. This func could
be injected in the `app/app.go`. One example as below:

Implement your own validation logic in a separate func, lets call it `externalTxDataValidationFuncExample`

```go
func externalTxDataValidationFuncExample(msg sdk.Msg) bool {
// make sure you do the type assertion for the tx that you want to validate
// in our case, it's MsgFeedData 

s := msg.(*types.MsgFeedData)
// some validation logic, e.g. feed data accuracy against CoinMarketCap.

return true
}
```

In the `app/app.go` file, you should see the code below where we set the `AnteHandler`:

```go
    app.SetAnteHandler(
chainlindkante.NewAnteHandler(app.AccountKeeper, app.BankKeeper, app.ChainLinkKeeper, ante.DefaultSigVerificationGasConsumer,
encodingConfig.TxConfig.SignModeHandler(), externalTxDataValidationFuncExample),
)
```

`externalTxDataValidationFuncExample` is where the tx data validation func should be injected.

Once the `Feed data submit` tx got broadcasted into the network, module will do the validation using this interface and
if the validation failed, module will trigger a transaction level event(MsgFeedDataValidationFailedEvent)
including all the feed data and feed info for further actions.

Please keep in mind that the injection func could also be `nil`, in this case there would be no data validation and the
validation result would be `true` and there is no `MsgFeedDataValidationFailedEvent` event gets triggerred.

## Register feed reward payout strategy function

Cosmos Chainlink module provides the ability that allows app developers to implement and register the feed reward payout
strategy functions in `app.go`. This will give the feed owner options to set the `feedRewardSchema` when creating a new
feed, or changing the `feedReward` parameter during the run time.

Each strategy function must have a name and the implementation associated when registering. Feed owner gets to pick
which strategy should be used to calculate the reward for the valid data providers.

example of registering strategies in `app.go`:

```go
    rewardStrategies := make(map[string]chainlinktypes.FeedRewardStrategyFunc)
rewardStrategies["accuracy"] = calculateByAccuracy
rewardStrategies["frequency"] = calculateByFrequency

chainlinktypes.NewFeedRewardStrategyRegister(rewardStrategies)
```

`NewFeedRewardStrategyRegister` takes a `map[string]chainlinktypes.FeedRewardStrategyFunc` as the argument,
`FeedRewardStrategyFunc` signature as below, it is defined in `x/types/feedRewardStrategy.go`

```go
    func(*MsgFeed, *MsgFeedData) ([]RewardPayout, error)
```

If `nil` is given when registering, no strategy will be available after chain launching, feed owner will not be able to
set any strategy in effect by issuing tx later. In which case, all the valid data providers will be rewarded by the base
amount.

Feed owner is also able to set the `feedRewardStrategy` of `feedReward` to `nil` by issuing a tx even though there are
available strategies registered, in which case, all the valid data providers will be rewarded by the base amount as
well.

Only registered strategies are available for a feed, CLI to query the list of available strategies:

```bash
get-feed-reward-avail-strategy
```


