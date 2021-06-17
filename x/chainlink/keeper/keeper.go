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

var roundId uint64 = 1

type (
	Keeper struct {
		cdc      codec.Marshaler
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetFeedData(ctx sdk.Context, feedData *types.MsgFeedData) {
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
		Submitter:    feedData.GetSubmitter(),
	}
	/****************/
	// TODO: verify deserializedOCRReport here

	finalFeedDataInStore := types.OCRFeedDataInStore{
		FeedData:              feedData,
		DeserializedOCRReport: &deserializedOCRReport,
		RoundId:               roundId,
	}
	// TODO: add round index here, need to figure out a way to generate round number, below line is a tmp solution
	roundId += 1

	// use store with gas meter
	store := ctx.KVStore(k.storeKey)

	f := k.cdc.MustMarshalBinaryBare(&finalFeedDataInStore)
	store.Set(types.KeyPrefix(types.FeedDataKey), f)
}

func (k Keeper) GetRoundFeedDataByFilter(ctx sdk.Context, req *types.GetRoundDataRequest) (*types.GetRoundDataResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var feedRoundData []*types.RoundData

	// use store with gas meter
	store := ctx.KVStore(k.storeKey)

	pageRes, err := query.Paginate(store, req.Pagination, func(key []byte, value []byte) error {
		var feedData types.OCRFeedDataInStore
		if err := k.cdc.UnmarshalBinaryBare(value, &feedData); err != nil {
			return err
		}

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
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefix(types.FeedDataKey))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var feedData types.OCRFeedDataInStore
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &feedData)

		latestRoundId := roundId - 1
		feedRoundData = feedDataFilter(req.GetFeedId(), latestRoundId, feedData)
	}

	return &types.GetLatestRoundDataResponse{
		RoundData: feedRoundData,
	}, nil
}
