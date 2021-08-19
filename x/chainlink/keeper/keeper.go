// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	Keeper struct {
		cdc                 codec.Marshaler
		bankKeeper          types.BankKeeper
		feedDataStoreKey    sdk.StoreKey
		roundStoreKey       sdk.StoreKey
		moduleOwnerStoreKey sdk.StoreKey
		feedInfoStoreKey    sdk.StoreKey
		accountStoreKey     sdk.StoreKey
		memKey              sdk.StoreKey
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	bk types.BankKeeper,
	feedDataStoreKey,
	roundStoreKey,
	moduleOwnerStoreKey,
	feedInfoStoreKey,
	accountStoreKey,
	memKey sdk.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		bankKeeper:          bk,
		feedDataStoreKey:    feedDataStoreKey,
		roundStoreKey:       roundStoreKey,
		moduleOwnerStoreKey: moduleOwnerStoreKey,
		feedInfoStoreKey:    feedInfoStoreKey,
		accountStoreKey:     accountStoreKey,
		memKey:              memKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetFeedData(ctx sdk.Context, feedData *types.MsgFeedData) (int64, []byte, error) {
	// emit FeedDataValidationFailedEvent event if feedData validation is false
	if !feedData.IsFeedDataValid {
		feed := k.GetFeed(ctx, feedData.GetFeedId())
		err := types.EmitEvent(&types.MsgFeedDataValidationFailedEvent{
			FeedId:        feedData.GetFeedId(),
			DataProviders: feed.GetFeed().GetDataProviders(),
			FeedOwner:     feed.GetFeed().GetFeedOwner(),
			Submitter:     feedData.GetSubmitter(),
			FeedData:      feedData.GetFeedData(),
			Signatures:    feedData.GetSignatures(),
		}, ctx.EventManager())
		if err != nil {
			// only trigger error message here cuz this will not affect the whole flow
			k.Logger(ctx).Error("failed to emit FeedDataValidationFailedEvent: ", err.Error())
		}
	}

	roundStore := ctx.KVStore(k.roundStoreKey)
	currentLatestRoundId := k.GetLatestRoundId(ctx, feedData.FeedId)
	roundId := currentLatestRoundId + 1

	// update the latest roundId of the current feedId
	roundStore.Set(types.GetRoundIdKey(feedData.GetFeedId()), i64tob(roundId))

	// TODO: add more complex feed validation here such as verify against other modules

	// TODO: deserialize the feedData.FeedData if it's an OCR report, assume all the feedData is OCR report for now.
	// this is simulating the OCR report deserialization lib
	/****************/
	observations := make([]*types.Observation, 0, len(feedData.GetFeedData()))
	for _, b := range feedData.GetFeedData() {
		o := &types.Observation{Data: []byte(string(b))}
		observations = append(observations, o)
	}
	deserializedOCRReport := types.OCRAbiEncoded{
		Context:      []byte(fmt.Sprintf("%d", roundId)),
		Oracles:      feedData.Submitter.Bytes(),
		Observations: observations,
	}
	/****************/
	// TODO: verify deserializedOCRReport here
	finalFeedDataInStore := types.OCRFeedDataInStore{
		FeedData:              feedData,
		DeserializedOCRReport: &deserializedOCRReport,
		RoundId:               roundId,
	}

	feedDataStore := ctx.KVStore(k.feedDataStoreKey)

	f := k.cdc.MustMarshalBinaryBare(&finalFeedDataInStore)

	feedDataStore.Set(types.GetFeedDataKey(feedData.GetFeedId(), strconv.FormatUint(roundId, 10)), f)

	// emit NewRoundData event
	err := types.EmitEvent(&types.MsgNewRoundDataEvent{
		FeedId:   feedData.FeedId,
		RoundId:  roundId,
		FeedData: feedData.FeedData,
	}, ctx.EventManager())
	if err != nil {
		return 0, nil, err
	}

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (k Keeper) GetRoundFeedDataByFilter(ctx sdk.Context, req *types.GetRoundDataRequest) (*types.GetRoundDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var feedRoundData []*types.RoundData

	feedDataStore := ctx.KVStore(k.feedDataStoreKey)

	pageRes, err := query.Paginate(feedDataStore, req.Pagination, func(key []byte, value []byte) error {
		var feedData types.OCRFeedDataInStore

		if err := k.cdc.UnmarshalBinaryBare(value, &feedData); err != nil {
			return err
		}

		data := feedDataFilter(req.GetFeedId(), req.GetRoundId(), feedData)
		if data != nil {
			feedRoundData = append(feedRoundData, data)
		}

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.GetRoundDataResponse{
		RoundData:  feedRoundData,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) GetLatestRoundFeedDataByFilter(ctx sdk.Context, req *types.GetLatestRoundDataRequest) (*types.GetLatestRoundDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var feedRoundData []*types.RoundData

	// get the roundId based on given feedId
	latestRoundId := k.GetLatestRoundId(ctx, req.GetFeedId())

	feedDataStore := ctx.KVStore(k.feedDataStoreKey)
	iterator := sdk.KVStorePrefixIterator(feedDataStore, types.GetFeedDataKey("", ""))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var feedData types.OCRFeedDataInStore
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &feedData)

		data := feedDataFilter(req.GetFeedId(), latestRoundId, feedData)
		if data != nil {
			feedRoundData = append(feedRoundData, data)
		}
	}

	return &types.GetLatestRoundDataResponse{
		RoundData: feedRoundData,
	}, nil
}

// GetLatestRoundId returns the current existing latest roundId of a feedId
// returns the global latest roundId in roundStore regardless of feedId if feedId is not given.
func (k Keeper) GetLatestRoundId(ctx sdk.Context, feedId string) uint64 {
	roundStore := ctx.KVStore(k.roundStoreKey)

	if feedId != "" {
		roundIdBytes := roundStore.Get(types.GetRoundIdKey(feedId))

		if len(roundIdBytes) == 0 {
			return 0
		}
		return btoi64(roundIdBytes)
	}

	var latestRoundId uint64
	roundIdIterator := sdk.KVStorePrefixIterator(roundStore, types.GetRoundIdKey(""))
	defer roundIdIterator.Close()

	for ; roundIdIterator.Valid(); roundIdIterator.Next() {
		roundId := btoi64(roundIdIterator.Value())
		if roundId > latestRoundId {
			latestRoundId = roundId
		}
	}

	return latestRoundId
}

func (k Keeper) SetModuleOwner(ctx sdk.Context, moduleOwner *types.MsgModuleOwner) (int64, []byte) {
	moduleStore := ctx.KVStore(k.moduleOwnerStoreKey)

	f := k.cdc.MustMarshalBinaryBare(moduleOwner)

	moduleStore.Set(types.GetModuleOwnerKey(moduleOwner.GetAddress().String()), f)

	return ctx.BlockHeight(), ctx.TxBytes()
}

func (k Keeper) RemoveModuleOwner(ctx sdk.Context, transfer *types.MsgModuleOwnershipTransfer) (int64, []byte) {
	moduleStore := ctx.KVStore(k.moduleOwnerStoreKey)

	moduleStore.Delete(types.GetModuleOwnerKey(transfer.GetAssignerAddress().String()))

	return ctx.BlockHeight(), ctx.TxBytes()
}

func (k Keeper) GetModuleOwnerList(ctx sdk.Context) *types.GetModuleOwnerResponse {
	moduleStore := ctx.KVStore(k.moduleOwnerStoreKey)
	iterator := sdk.KVStorePrefixIterator(moduleStore, types.GetModuleOwnerKey(""))

	defer iterator.Close()

	moduleOwners := make([]*types.MsgModuleOwner, 0)

	for ; iterator.Valid(); iterator.Next() {
		var moduleOwner types.MsgModuleOwner
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &moduleOwner)

		moduleOwners = append(moduleOwners, &moduleOwner)
	}

	return &types.GetModuleOwnerResponse{
		ModuleOwner: moduleOwners,
	}
}

func (k Keeper) SetFeed(ctx sdk.Context, feed *types.MsgFeed) (int64, []byte) {
	feedInfoStore := ctx.KVStore(k.feedInfoStoreKey)

	f := k.cdc.MustMarshalBinaryBare(feed)

	feedInfoStore.Set(types.GetFeedInfoKey(feed.GetFeedId()), f)

	return ctx.BlockHeight(), ctx.TxBytes()
}

func (k Keeper) GetFeed(ctx sdk.Context, feedId string) *types.GetFeedByIdResponse {
	feedInfoStore := ctx.KVStore(k.feedInfoStoreKey)
	feedIdBytes := feedInfoStore.Get(types.GetFeedInfoKey(feedId))

	if feedIdBytes == nil {
		return &types.GetFeedByIdResponse{
			Feed: nil,
		}
	}

	var feed types.MsgFeed
	k.cdc.MustUnmarshalBinaryBare(feedIdBytes, &feed)

	return &types.GetFeedByIdResponse{
		Feed: &feed,
	}
}

func (k Keeper) AddDataProvider(ctx sdk.Context, addDataProvider *types.MsgAddDataProvider) (int64, []byte, error) {
	// retrieve feed from store
	resp := k.GetFeed(ctx, addDataProvider.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", addDataProvider.GetFeedId())
	}

	// add new data provider
	feed.DataProviders = append(feed.DataProviders, addDataProvider.DataProvider)

	// put back feed in the store
	k.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (k Keeper) RemoveDataProvider(ctx sdk.Context, removeDataProvider *types.MsgRemoveDataProvider) (int64, []byte, error) {
	// retrieve feed from store
	resp := k.GetFeed(ctx, removeDataProvider.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", removeDataProvider.GetFeedId())
	}

	// remove data provider from the list
	feed.DataProviders = (types.DataProviders)(feed.DataProviders).Remove(removeDataProvider.GetAddress())

	// put back feed in the store
	k.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (k Keeper) SetSubmissionCount(ctx sdk.Context, setSubmissionCount *types.MsgSetSubmissionCount) (int64, []byte, error) {
	// retrieve feed from store
	resp := k.GetFeed(ctx, setSubmissionCount.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", setSubmissionCount.GetFeedId())
	}

	// update submission count
	feed.SubmissionCount = setSubmissionCount.GetSubmissionCount()

	// put back feed in the store
	k.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (k Keeper) SetHeartbeatTrigger(ctx sdk.Context, setHeartbeatTrigger *types.MsgSetHeartbeatTrigger) (int64, []byte, error) {
	// retrieve feed from store
	resp := k.GetFeed(ctx, setHeartbeatTrigger.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", setHeartbeatTrigger.GetFeedId())
	}

	// update heartbeat trigger
	feed.HeartbeatTrigger = setHeartbeatTrigger.GetHeartbeatTrigger()

	// put back feed in the store
	k.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (k Keeper) SetDeviationThresholdTrigger(ctx sdk.Context, setDeviationThresholdTrigger *types.MsgSetDeviationThresholdTrigger) (int64, []byte, error) {
	// retrieve feed from store
	resp := k.GetFeed(ctx, setDeviationThresholdTrigger.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", setDeviationThresholdTrigger.GetFeedId())
	}

	// update deviation threshold trigger
	feed.DeviationThresholdTrigger = setDeviationThresholdTrigger.GetDeviationThresholdTrigger()

	// put back feed in the store
	k.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (k Keeper) SetFeedReward(ctx sdk.Context, setFeedReward *types.MsgSetFeedReward) (int64, []byte, error) {
	// retrieve feed from store
	resp := k.GetFeed(ctx, setFeedReward.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", setFeedReward.GetFeedId())
	}

	// update feed reward
	feed.FeedReward = setFeedReward.GetFeedReward()

	// put back feed in the store
	k.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

// DistributeReward will mint the reward from the module
// then transfer the reward to the receiver (data provider)
func (k Keeper) DistributeReward(ctx sdk.Context, msg *types.MsgFeedData, dataProviders []*types.DataProvider, feedReward uint32) error {
	// calculate the total reward to mint (minus fee compensation)
	totalFeedReward := int64(feedReward) * int64(len(dataProviders))
	tokensToMint := types.NewLinkCoinInt64(totalFeedReward)
	tokensToSend := types.NewLinkCoinInt64(int64(feedReward))

	OraclePaidEvent := &types.MsgOraclePaidEvent{
		FeedId: msg.FeedId,
	}

	// mint new tokens if the source of the transfer is the same chain
	if err := k.bankKeeper.MintCoins(
		ctx, types.ModuleName, sdk.NewCoins(tokensToMint),
	); err != nil {
		return err
	}

	// distribute reward to all data providers except submitter
	for _, dp := range dataProviders {
		if dp.Address.String() != msg.Submitter.String() {
			if err := k.bankKeeper.SendCoinsFromModuleToAccount(
				ctx, types.ModuleName, dp.Address, sdk.NewCoins(tokensToSend),
			); err != nil {
				return err
			}
			// emit OraclePaid event for valid data providers
			OraclePaidEvent.Account = dp.Address
			OraclePaidEvent.Value = uint64(feedReward)
			err := types.EmitEvent(OraclePaidEvent, ctx.EventManager())
			if err != nil {
				return err
			}
		}
	}

	// send to submitter
	// TODO: include fees - need to mint this amount as well
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, msg.Submitter, sdk.NewCoins(tokensToSend),
	); err != nil {
		return err
	}

	// emit OraclePaid event to submitter including the fee
	OraclePaidEvent.Account = msg.Submitter
	// TODO: event.Value =  uint64(feedReward) + fee
	OraclePaidEvent.Value = uint64(feedReward)
	err := types.EmitEvent(OraclePaidEvent, ctx.EventManager())
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) FeedOwnershipTransfer(ctx sdk.Context, feedOwnershipTransfer *types.MsgFeedOwnershipTransfer) (int64, []byte, error) {
	// retrieve feed from store
	resp := k.GetFeed(ctx, feedOwnershipTransfer.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", feedOwnershipTransfer.GetFeedId())
	}

	// update the feed owner
	feed.FeedOwner = feedOwnershipTransfer.GetNewFeedOwnerAddress()

	// put back feed in the store
	k.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

// RequestNewRound will be a transaction sent by the FeedOwner to request a new report to the chainlink network
// The event emitted will expect a data provider to submit a new report.
func (k Keeper) RequestNewRound(ctx sdk.Context, requestNewRound *types.MsgRequestNewRound) (int64, []byte, error) {
	// emit NewRoundData event
	err := types.EmitEvent(&types.MsgNewRoundRequestEvent{
		FeedId: requestNewRound.FeedId,
	}, ctx.EventManager())
	if err != nil {
		return 0, nil, err
	}

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (k Keeper) AddAccount(ctx sdk.Context, acc *types.MsgAccount) (int64, []byte) {
	accStore := ctx.KVStore(k.accountStoreKey)

	a := k.cdc.MustMarshalBinaryBare(acc)

	accStore.Set(types.GetAccountKey(acc.GetSubmitter().String()), a)

	return ctx.BlockHeight(), ctx.TxBytes()
}

func (k Keeper) EditAccount(ctx sdk.Context, acc *types.MsgEditAccount) (int64, []byte, error) {
	accStore := ctx.KVStore(k.accountStoreKey)
	accSubmitter := acc.Submitter.String()
	accountBytes := accStore.Get(types.GetAccountKey(accSubmitter))

	var account types.MsgAccount
	k.cdc.MustUnmarshalBinaryBare(accountBytes, &account)

	// submitters must match
	if acc.Submitter.String() != account.Submitter.String() {
		return 0, []byte{}, fmt.Errorf("submitter does not match")
	}

	// overwrite the piggy address
	account.PiggyAddress = acc.PiggyAddress

	a := k.cdc.MustMarshalBinaryBare(&account)
	accStore.Set(types.GetAccountKey(acc.GetSubmitter().String()), a)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (k Keeper) GetAccount(ctx sdk.Context, accReq *types.GetAccountRequest) *types.GetAccountResponse {
	acc := accReq.AccountAddress.String()
	accStore := ctx.KVStore(k.accountStoreKey)
	accountBytes := accStore.Get(types.GetAccountKey(acc))

	var account types.MsgAccount
	k.cdc.MustUnmarshalBinaryBare(accountBytes, &account)

	return &types.GetAccountResponse{
		Account: &account,
	}
}
