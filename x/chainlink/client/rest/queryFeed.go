package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
)

func listFeedDataHandler(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		feedId := vars["feedId"]

		fmt.Println("???????", feedId)

		params := &types.QueryFeedDataRequest{
			FeedId:     feedId,
			Pagination: nil,
		}

		queryClient := types.NewQueryClient(clientCtx)
		res, err := queryClient.FeedDataByID(context.Background(), params)
		if err != nil {
			return
		}

		//res, height, err := clientCtx.QueryWithData(fmt.Sprintf("%s/list-feedData", types.QuerierRoute), nil)
		//if err != nil {
		//	rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		//	return
		//}

		//clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}
