// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package cli

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

func CmdAddFeed() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addFeed [feedId] [feedOwnerAddress] [submissionCount] [heartbeatTrigger] [deviationThresholdTrigger] [feedReward] [initDataProviderList]",
		Short: "Add new feed. Signer must be the existing module owner.",
		Long:  "The following fields are required:\n\tThe feedId will be a string that uniquely identifies the feed. The feedOwnerAddress must be a valid cosmos address.\n\tThe submissionCount in the required number of signatures.\n\tThe deviationThresholdTrigger is the fraction of deviation in the feed data required to trigger a new round.\n\tThe initDataProviderList is a string contains each data provider's address with pubkey and split by comma.\n\tThe feedReward is a uint32 value that represents the data provider reward for submitting data to a feed.",
		Args:  cobra.MinimumNArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFeedId := args[0]
			argsFeedOwnerAddr := args[1]
			argsSubmissionCount := args[2]
			argsHeartbeatTrigger := args[3]
			argsDeviationThresholdTrigger := args[4]
			argsFeedReward := args[5]
			argsInitDataProviderListStr := strings.TrimSpace(args[6])

			submissionCount, err := strconv.ParseUint(argsSubmissionCount, 10, 32)
			if err != nil {
				return err
			}
			heartbeatTrigger, err := strconv.ParseUint(argsHeartbeatTrigger, 10, 32)
			if err != nil {
				return err
			}
			deviationThresholdTrigger, err := strconv.ParseUint(argsDeviationThresholdTrigger, 10, 32)
			if err != nil {
				return err
			}
			feedReward, err := strconv.ParseUint(argsFeedReward, 10, 32)
			if err != nil {
				return err
			}

			argsInitDataProviderList := strings.Split(argsInitDataProviderListStr, ",")
			if len(argsInitDataProviderList)%2 != 0 {
				return errors.New("invalid init data provider pairs")
			}

			initDataProviderList := make([]*types.DataProvider, 0, len(argsInitDataProviderList)/2)
			i := 0
			for i < len(argsInitDataProviderList) {
				addr, err := sdk.AccAddressFromBech32(strings.TrimSpace(argsInitDataProviderList[i]))
				if err != nil {
					return sdkerrors.Wrapf(err, "invalid init data provider address: %s", argsInitDataProviderList[i])
				}

				initDataProviderList = append(initDataProviderList, &types.DataProvider{
					Address: addr,
					PubKey:  []byte(strings.TrimSpace(argsInitDataProviderList[i+1])),
				})
				i = i + 2
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			feedOwnerAddr, err := sdk.AccAddressFromBech32(argsFeedOwnerAddr)
			if err != nil {
				return err
			}

			msg := types.NewMsgFeed(argsFeedId, feedOwnerAddr, clientCtx.GetFromAddress(), initDataProviderList, uint32(submissionCount), uint32(heartbeatTrigger), uint32(deviationThresholdTrigger), uint32(feedReward))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdAddDataProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addDataProvider [feedId] [address] [publicKey]",
		Short: "Add new data provider to the feed. Signer must be the existing module owner.",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFeedId := args[0]
			argsAddress := args[1]
			argsPublicKey := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addr, err := sdk.AccAddressFromBech32(argsAddress)
			if err != nil {
				return err
			}

			msg := types.NewMsgAddDataProvider(clientCtx.GetFromAddress(), argsFeedId, &types.DataProvider{
				Address: addr,
				PubKey:  []byte(argsPublicKey),
			})
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdRemoveDataProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "removeDataProvider [feedId] [address]",
		Short: "Remove data provider from the feed. Signer must be the existing module owner.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFeedId := args[0]
			argsAddress := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addr, err := sdk.AccAddressFromBech32(argsAddress)
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveDataProvider(clientCtx.GetFromAddress(), argsFeedId, addr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSetSubmissionCount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setSubmissionCount [feedId] [count]",
		Short: "Sets a new submission count for a given feed",
		Long:  "Set the required number of signatures. Signer must be the existing module owner.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFeedId := args[0]
			argsSubmissionCount := args[1]

			submissionCount, err := strconv.ParseUint(argsSubmissionCount, 10, 32)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetSubmissionCount(clientCtx.GetFromAddress(), argsFeedId, uint32(submissionCount))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSetHeartbeatTrigger() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setHeartbeatTrigger [feedId] [heartbeatTrigger]",
		Short: "Sets a new heartbeat trigger for the given feed",
		Long:  "Set the interval between which a new round should automatically be triggered. Signer must be the existing module owner.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFeedId := args[0]
			argsHeartbeatTrigger := args[1]

			heartbeatTrigger, err := strconv.ParseUint(argsHeartbeatTrigger, 10, 32)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetHeartbeatTrigger(clientCtx.GetFromAddress(), argsFeedId, uint32(heartbeatTrigger))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSetDeviationThreshold() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setDeviationThresholdTrigger [feedId] [deviationThresholdTrigger]",
		Short: "Sets a new deviation threshold trigger for the given feed",
		Long:  "Set the fraction of deviation in the feed data required to trigger a new round. Signer must be the existing module owner.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFeedId := args[0]
			argsDeviationThresholdTrigger := args[1]

			deviationThresholdTrigger, err := strconv.ParseUint(argsDeviationThresholdTrigger, 10, 32)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetDeviationThreshold(clientCtx.GetFromAddress(), argsFeedId, uint32(deviationThresholdTrigger))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSetFeedReward() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setFeedReward [feedId] [feedReward]",
		Short: "Sets a new feed reward for the given feed",
		Long:  "Set the feed reward for a given feed, the reward will be distributed in tokens denominated as 'link'. Signer must be the existing module owner.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFeedId := args[0]
			argsFeedReward := args[1]

			feedReward, err := strconv.ParseUint(argsFeedReward, 10, 32)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetFeedReward(clientCtx.GetFromAddress(), argsFeedId, uint32(feedReward))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdTransferFeedOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feedOwnershipTransfer [feedId] [newFeedOwnerAddress]",
		Short: "Transfer chainLink feed ownership from an existing feed owner account to another account. Signer must be an existing feed owner.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFeedId := args[0]
			argsAddress := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addr, err := sdk.AccAddressFromBech32(argsAddress)
			if err != nil {
				return err
			}

			msg := types.NewMsgFeedOwnershipTransfer(clientCtx.GetFromAddress(), argsFeedId, addr)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSubmitFeedData() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submitFeedData [feedId] [feedData] [signatures]",
		Short: "Submit feed data",
		Long:  "Submit feed data, called by an OCR round leader to submit an off-chain report of data signed by a number of oracles.",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFeedId := args[0]
			argsFeedData := args[1]
			argsSignatures := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// TODO: this is dummy data to simulate the data providers signature set
			signatures := strings.Split(argsSignatures, " ")
			s := make([][]byte, 0)
			for _, sign := range signatures {
				s = append(s, []byte(sign))
			}

			msg := types.NewMsgFeedData(clientCtx.GetFromAddress(), argsFeedId, []byte(argsFeedData), s)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
