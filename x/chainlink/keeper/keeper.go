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

func (s Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (s Keeper) SetFeedData(ctx sdk.Context, feedData *types.MsgFeedData) (int64, []byte) {
	roundStore := ctx.KVStore(s.roundStoreKey)
	currentLatestRoundId := s.GetLatestRoundId(ctx, feedData.FeedId)
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

	feedDateStore := ctx.KVStore(s.feedDataStoreKey)

	f := s.cdc.MustMarshalBinaryBare(&finalFeedDataInStore)

	feedDateStore.Set(types.GetFeedDataKey(feedData.GetFeedId(), strconv.FormatUint(roundId, 10)), f)

	return ctx.BlockHeight(), ctx.TxBytes()
}

func (s Keeper) GetRoundFeedDataByFilter(ctx sdk.Context, req *types.GetRoundDataRequest) (*types.GetRoundDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var feedRoundData []*types.RoundData

	feedDataStore := ctx.KVStore(s.feedDataStoreKey)

	pageRes, err := query.Paginate(feedDataStore, req.Pagination, func(key []byte, value []byte) error {
		var feedData types.OCRFeedDataInStore

		if err := s.cdc.UnmarshalBinaryBare(value, &feedData); err != nil {
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

func (s Keeper) GetLatestRoundFeedDataByFilter(ctx sdk.Context, req *types.GetLatestRoundDataRequest) (*types.GetLatestRoundDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var feedRoundData []*types.RoundData

	// get the roundId based on given feedId
	latestRoundId := s.GetLatestRoundId(ctx, req.GetFeedId())

	feedDataStore := ctx.KVStore(s.feedDataStoreKey)
	iterator := sdk.KVStorePrefixIterator(feedDataStore, types.GetFeedDataKey("", ""))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var feedData types.OCRFeedDataInStore
		s.cdc.MustUnmarshalBinaryBare(iterator.Value(), &feedData)

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
func (s Keeper) GetLatestRoundId(ctx sdk.Context, feedId string) uint64 {
	roundStore := ctx.KVStore(s.roundStoreKey)

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

func (s Keeper) SetModuleOwner(ctx sdk.Context, moduleOwner *types.MsgModuleOwner) (int64, []byte) {
	moduleStore := ctx.KVStore(s.moduleOwnerStoreKey)

	f := s.cdc.MustMarshalBinaryBare(moduleOwner)

	moduleStore.Set(types.GetModuleOwnerKey(moduleOwner.GetAddress().String()), f)

	return ctx.BlockHeight(), ctx.TxBytes()
}

func (s Keeper) RemoveModuleOwner(ctx sdk.Context, transfer *types.MsgModuleOwnershipTransfer) (int64, []byte) {
	moduleStore := ctx.KVStore(s.moduleOwnerStoreKey)

	moduleStore.Delete(types.GetModuleOwnerKey(transfer.GetAssignerAddress().String()))

	return ctx.BlockHeight(), ctx.TxBytes()
}

func (s Keeper) GetModuleOwnerList(ctx sdk.Context) *types.GetModuleOwnerResponse {
	moduleStore := ctx.KVStore(s.moduleOwnerStoreKey)
	iterator := sdk.KVStorePrefixIterator(moduleStore, types.GetModuleOwnerKey(""))

	defer iterator.Close()

	moduleOwners := make([]*types.MsgModuleOwner, 0)

	for ; iterator.Valid(); iterator.Next() {
		var moduleOwner types.MsgModuleOwner
		s.cdc.MustUnmarshalBinaryBare(iterator.Value(), &moduleOwner)

		moduleOwners = append(moduleOwners, &moduleOwner)
	}

	return &types.GetModuleOwnerResponse{
		ModuleOwner: moduleOwners,
	}
}

func (s Keeper) SetFeed(ctx sdk.Context, feed *types.MsgFeed) (int64, []byte) {
	feedInfoStore := ctx.KVStore(s.feedInfoStoreKey)

	f := s.cdc.MustMarshalBinaryBare(feed)

	feedInfoStore.Set(types.GetFeedInfoKey(feed.GetFeedId()), f)

	return ctx.BlockHeight(), ctx.TxBytes()
}

func (s Keeper) GetFeed(ctx sdk.Context, feedId string) *types.GetFeedByIdResponse {
	feedInfoStore := ctx.KVStore(s.feedInfoStoreKey)
	feedIdBytes := feedInfoStore.Get(types.GetFeedInfoKey(feedId))

	if feedIdBytes == nil {
		return &types.GetFeedByIdResponse{
			Feed: nil,
		}
	}

	var feed types.MsgFeed
	s.cdc.MustUnmarshalBinaryBare(feedIdBytes, &feed)

	return &types.GetFeedByIdResponse{
		Feed: &feed,
	}
}

func (s Keeper) AddDataProvider(ctx sdk.Context, addDataProvider *types.MsgAddDataProvider) (int64, []byte, error) {
	// retrieve feed from store
	resp := s.GetFeed(ctx, addDataProvider.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", addDataProvider.GetFeedId())
	}

	// add new data provider
	feed.DataProviders = append(feed.DataProviders, addDataProvider.DataProvider)

	// put back feed in the store
	s.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (s Keeper) RemoveDataProvider(ctx sdk.Context, removeDataProvider *types.MsgRemoveDataProvider) (int64, []byte, error) {
	// retrieve feed from store
	resp := s.GetFeed(ctx, removeDataProvider.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", removeDataProvider.GetFeedId())
	}

	// remove data provider from the list
	feed.DataProviders = (types.DataProviders)(feed.DataProviders).Remove(removeDataProvider.GetAddress())

	// put back feed in the store
	s.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (s Keeper) SetSubmissionCount(ctx sdk.Context, setSubmissionCount *types.MsgSetSubmissionCount) (int64, []byte, error) {
	// retrieve feed from store
	resp := s.GetFeed(ctx, setSubmissionCount.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", setSubmissionCount.GetFeedId())
	}

	// update submission count
	feed.SubmissionCount = setSubmissionCount.GetSubmissionCount()

	// put back feed in the store
	s.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (s Keeper) SetHeartbeatTrigger(ctx sdk.Context, setHeartbeatTrigger *types.MsgSetHeartbeatTrigger) (int64, []byte, error) {
	// retrieve feed from store
	resp := s.GetFeed(ctx, setHeartbeatTrigger.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", setHeartbeatTrigger.GetFeedId())
	}

	// update heartbeat trigger
	feed.HeartbeatTrigger = setHeartbeatTrigger.GetHeartbeatTrigger()

	// put back feed in the store
	s.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (s Keeper) SetDeviationThresholdTrigger(ctx sdk.Context, setDeviationThresholdTrigger *types.MsgSetDeviationThresholdTrigger) (int64, []byte, error) {
	// retrieve feed from store
	resp := s.GetFeed(ctx, setDeviationThresholdTrigger.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", setDeviationThresholdTrigger.GetFeedId())
	}

	// update deviation threshold trigger
	feed.DeviationThresholdTrigger = setDeviationThresholdTrigger.GetDeviationThresholdTrigger()

	// put back feed in the store
	s.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}

func (s Keeper) FeedOwnershipTransfer(ctx sdk.Context, feedOwnershipTransfer *types.MsgFeedOwnershipTransfer) (int64, []byte, error) {
	// retrieve feed from store
	resp := s.GetFeed(ctx, feedOwnershipTransfer.GetFeedId())
	feed := resp.GetFeed()
	if feed == nil {
		return 0, nil, fmt.Errorf("feed '%s' not found", feedOwnershipTransfer.GetFeedId())
	}

	// update the feed owner
	feed.FeedOwner = feedOwnershipTransfer.GetNewFeedOwnerAddress()

	// put back feed in the store
	s.SetFeed(ctx, feed)

	return ctx.BlockHeight(), ctx.TxBytes(), nil
}
