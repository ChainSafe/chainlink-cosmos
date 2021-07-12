package keeper

import (
	"fmt"
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
		feedDataStoreKey    sdk.StoreKey
		roundStoreKey       sdk.StoreKey
		moduleOwnerStoreKey sdk.StoreKey
		feedInfoStoreKey    sdk.StoreKey
		memKey              sdk.StoreKey
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	feedDataStoreKey,
	roundStoreKey,
	moduleOwnerStoreKey,
	feedInfoStoreKey,
	memKey sdk.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:                 cdc,
		feedDataStoreKey:    feedDataStoreKey,
		roundStoreKey:       roundStoreKey,
		moduleOwnerStoreKey: moduleOwnerStoreKey,
		feedInfoStoreKey:    feedInfoStoreKey,
		memKey:              memKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetFeedData(ctx sdk.Context, feedData *types.MsgFeedData) (int64, []byte) {
	roundStore := ctx.KVStore(k.roundStoreKey)
	currentLatestRoundId := k.GetLatestRoundId(roundStore, feedData.FeedId)
	roundId := currentLatestRoundId + 1

	// update the latest roundId of the current feedId
	roundStore.Set(types.KeyPrefix(types.RoundIdKey+feedData.FeedId), i64tob(roundId))

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

	feedDateStore := ctx.KVStore(k.feedDataStoreKey)

	f := k.cdc.MustMarshalBinaryBare(&finalFeedDataInStore)

	feedDateStore.Set(types.KeyPrefix(types.FeedDataKey+feedData.FeedId+fmt.Sprintf("%d", roundId)), f)

	return ctx.BlockHeight(), ctx.TxBytes()
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
	roundStore := ctx.KVStore(k.roundStoreKey)
	latestRoundId := k.GetLatestRoundId(roundStore, req.GetFeedId())

	feedDataStore := ctx.KVStore(k.feedDataStoreKey)
	iterator := sdk.KVStorePrefixIterator(feedDataStore, types.KeyPrefix(types.FeedDataKey))

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
func (k Keeper) GetLatestRoundId(store sdk.KVStore, feedId string) uint64 {
	if feedId != "" {
		feedRoundIdKey := types.KeyPrefix(types.RoundIdKey + feedId)
		roundIdBytes := store.Get(feedRoundIdKey)

		if len(roundIdBytes) == 0 {
			return 0
		}
		return btoi64(roundIdBytes)
	}

	var latestRoundId uint64
	roundIdIterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.RoundIdKey))
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

	moduleStore.Set(types.KeyPrefix(types.ModuleOwnerKey+moduleOwner.GetAddress().String()), f)

	return ctx.BlockHeight(), ctx.TxBytes()
}

func (k Keeper) RemoveModuleOwner(ctx sdk.Context, transfer *types.MsgModuleOwnershipTransfer) (int64, []byte) {
	moduleStore := ctx.KVStore(k.moduleOwnerStoreKey)

	moduleStore.Delete(types.KeyPrefix(types.ModuleOwnerKey + transfer.GetAssignerAddress().String()))

	return ctx.BlockHeight(), ctx.TxBytes()
}

func (k Keeper) GetModuleOwnerList(ctx sdk.Context) *types.GetModuleOwnerResponse {
	moduleStore := ctx.KVStore(k.moduleOwnerStoreKey)
	iterator := sdk.KVStorePrefixIterator(moduleStore, types.KeyPrefix(types.ModuleOwnerKey))

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

	feedInfoStore.Set(types.KeyPrefix(types.FeedInfoKey+feed.GetFeedId()), f)

	return ctx.BlockHeight(), ctx.TxBytes()
}

func (k Keeper) GetFeed(ctx sdk.Context, feedId string) *types.GetFeedByIdResponse {
	feedInfoStore := ctx.KVStore(k.feedInfoStoreKey)

	feedKey := types.KeyPrefix(types.FeedInfoKey + feedId)

	feedIdBytes := feedInfoStore.Get(feedKey)

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

func (k Keeper) AddFeedProvider(ctx sdk.Context, addFeedProvider *types.MsgAddFeedProvider) (int64, []byte, error) {
	// retrieve feed from store
	resp := k.GetFeed(ctx, addFeedProvider.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", addFeedProvider.GetFeedId())
	}

	// add new data provider
	feed.DataProviders = append(feed.DataProviders, addFeedProvider.DataProvider)

	// put back feed in the store
	f := k.cdc.MustMarshalBinaryBare(feed)
	feedInfoStore := ctx.KVStore(k.feedInfoStoreKey)
	feedInfoStore.Set(types.KeyPrefix(types.FeedInfoKey+feed.GetFeedId()), f)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (k Keeper) RemoveFeedProvider(ctx sdk.Context, removeFeedProvider *types.MsgRemoveFeedProvider) (int64, []byte, error) {
	// retrieve feed from store
	resp := k.GetFeed(ctx, removeFeedProvider.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", removeFeedProvider.GetFeedId())
	}

	// remove data provider from the list
	feed.DataProviders = (types.DataProviders)(feed.DataProviders).Remove(removeFeedProvider.GetAddress())

	// put back feed in the store
	f := k.cdc.MustMarshalBinaryBare(feed)
	feedInfoStore := ctx.KVStore(k.feedInfoStoreKey)
	feedInfoStore.Set(types.KeyPrefix(types.FeedInfoKey+feed.GetFeedId()), f)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}
