// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package cli

import (
	"context"
	"errors"
	"strconv"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetFeedDataByRound() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-round-feed-data [roundId] [feedId]",
		Short: "List feed data by round. roundId is required, feedId is optional.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("roundId is required")
			}
			roundId := args[0]
			roundIdInt, err := strconv.ParseInt(roundId, 10, 64)
			if err != nil {
				return errors.New("roundId is invalid")
			}

			var feedId string
			if len(args) >= 2 {
				feedId = args[1]
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err = client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.GetRoundDataRequest{
				RoundId:    uint64(roundIdInt),
				FeedId:     feedId,
				Pagination: pageReq,
			}

			res, err := queryClient.GetRoundData(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdGetLatestFeedData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-latest-feed-data [feedId]",
		Short: "List the latest round feed data. feedId is optional.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var feedId string
			if len(args) != 0 {
				feedId = args[0]
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			params := &types.GetLatestRoundDataRequest{
				FeedId: feedId,
			}

			res, err := queryClient.LatestRoundData(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdGetFeedInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-feed-info [feedId]",
		Short: "Get feed info by feedId",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var feedId string
			if len(args) != 0 {
				feedId = args[0]
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			params := &types.GetFeedByIdRequest{FeedId: feedId}

			res, err := queryClient.GetFeedByFeedId(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdGetFeedRewardAvailStrategy() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getFeedRewardAvailStrategy",
		Short: "Get feed reward available strategies info",
		RunE: func(cmd *cobra.Command, args []string) error {

			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.GetFeedRewardAvailStrategiesRequest{}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetFeedRewardAvailStrategy(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
