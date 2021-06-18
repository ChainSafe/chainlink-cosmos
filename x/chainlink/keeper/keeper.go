package keeper

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: implement a map to maintain roundId in memory for now
// data structure:  map[feedId]roundId
// var roundId uint64 = 1
var roundId uint64

type (
	Keeper struct {
		cdc           codec.Marshaler
		feedStoreKey  sdk.StoreKey
		roundStoreKey sdk.StoreKey
		memKey        sdk.StoreKey
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	feedStoreKey,
	roundStoreKey,
	memKey sdk.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:           cdc,
		feedStoreKey:  feedStoreKey,
		roundStoreKey: roundStoreKey,
		memKey:        memKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetFeedData(ctx sdk.Context, feedData *types.MsgFeedData) (int64, []byte) {
	// use store with gas meter
	roundStore := ctx.KVStore(k.roundStoreKey)

	feedRoundIdKey := types.KeyPrefix(types.RoundIdKey + feedData.FeedId)
	roundIdBytes := roundStore.Get(feedRoundIdKey)
	if len(roundIdBytes) == 0 {
		roundId = 1
		roundStore.Set(types.KeyPrefix(types.RoundIdKey+feedData.FeedId), i64tob(roundId))
	} else {
		oldRoundId := btoi64(roundIdBytes)
		roundId = oldRoundId + 1
		roundStore.Set(types.KeyPrefix(types.RoundIdKey+feedData.FeedId), i64tob(roundId))
	}

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
		Context:      []byte("testcontext"),
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

	feedStore := ctx.KVStore(k.feedStoreKey)

	f := k.cdc.MustMarshalBinaryBare(&finalFeedDataInStore)
	// will require the feedDataKey + feedId + roundId to set.
	feedStore.Set(types.KeyPrefix(types.FeedDataKey+feedData.FeedId+fmt.Sprintf("%d", roundId)), f)

	return ctx.BlockHeight(), ctx.TxBytes()
}

// func (k Keeper) GetRoundFeedDataByFeedAndRoundId()

// func (k Keeper) GerRoundFeedDataByRoundId()

func (k Keeper) GetRoundFeedDataByFilter(ctx sdk.Context, req *types.GetRoundDataRequest) (*types.GetRoundDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var feedRoundData []*types.RoundData

	// use store with gas meter
	feedStore := ctx.KVStore(k.feedStoreKey)

	pageRes, err := query.Paginate(feedStore, req.Pagination, func(key []byte, value []byte) error {
		var feedData types.OCRFeedDataInStore

		if err := k.cdc.UnmarshalBinaryBare(value, &feedData); err != nil {
			return err
		}

		// TODO: verify that the roundId procures the correct feedData
		feedRoundData = feedDataFilter(req.GetFeedId(), req.GetRoundId(), feedData)

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

	// use store with gas meter
	roundStore := ctx.KVStore(k.roundStoreKey)
	feedRoundIdKey := types.KeyPrefix(types.RoundIdKey + req.FeedId)

	// use store with gas meter
	feedStore := ctx.KVStore(k.feedStoreKey)
	iterator := sdk.KVStorePrefixIterator(feedStore, types.KeyPrefix(types.FeedDataKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var feedData types.OCRFeedDataInStore
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &feedData)

		roundIdBytes := roundStore.Get(feedRoundIdKey)
		if len(roundIdBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The provided feedId in req does not have any roundId associated.")
		}
		roundId := btoi64(roundIdBytes)
		// TODO: update the feedDataFilter according to the in memory roundId
		feedRoundData = feedDataFilter(req.GetFeedId(), roundId, feedData)
	}

	return &types.GetLatestRoundDataResponse{
		RoundData: feedRoundData,
	}, nil
}

func (k Keeper) GetLatestRoundId(ctx sdk.Context, feedId string) []byte {
	return getLatestRoundId(k, ctx, feedId)
}

func getLatestRoundId(k Keeper, ctx sdk.Context, feedId string) []byte {
	// use store with gas meter
	roundStore := ctx.KVStore(k.roundStoreKey)
	feedRoundIdKey := types.KeyPrefix(types.RoundIdKey + feedId)
	roundIdBytes := roundStore.Get(feedRoundIdKey)
	if len(roundIdBytes) == 0 {
		return []byte{}
	} else {
		return roundIdBytes
	}
}
