package cli

import (
	"errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
	"strings"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func CmdAddModuleOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addModuleOwner [address] [publicKey]",
		Short: "Add ChainLink Module Owner. Signer must be an existing module owner.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsAddress := args[0]
			argsPublicKey := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addr, err := sdk.AccAddressFromBech32(argsAddress)
			if err != nil {
				return err
			}

			msg := types.NewMsgModuleOwner(clientCtx.GetFromAddress(), addr, []byte(argsPublicKey))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdTransferModuleOwnership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "moduleOwnershipTransfer [newModuleOwnerAddress] [newModuleOwnerPublicKey]",
		Short: "Transfer ChainLink Module Ownership from an existing module owner account to another account. Signer must be the existing module owner.",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsAddress := args[0]
			argsPublicKey := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			addr, err := sdk.AccAddressFromBech32(argsAddress)
			if err != nil {
				return err
			}

			msg := types.NewMsgModuleOwnershipTransfer(clientCtx.GetFromAddress(), addr, []byte(argsPublicKey))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdAddFeed() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addFeed [feedId] [feedOwnerAddress] [submissionCount] [heartbeatTrigger] [deviationThresholdTrigger] [initDataProviderList]",
		Short: "Add new feed. Signer must be the existing module owner.",
		Args:  cobra.MinimumNArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFeedId := strings.TrimSpace(args[0])
			argsFeedOwnerAddr := strings.TrimSpace(args[1])
			argsSubmissionCount := strings.TrimSpace(args[2])
			argsHeartbeatTrigger := strings.TrimSpace(args[3])
			argsDeviationThresholdTrigger := strings.TrimSpace(args[4])
			argsInitDataProviderListStr := strings.TrimSpace(args[5])

			submissionCount, err := strconv.Atoi(argsSubmissionCount)
			if err != nil {
				return err
			}
			heartbeatTrigger, err := strconv.Atoi(argsHeartbeatTrigger)
			if err != nil {
				return err
			}
			deviationThresholdTrigger, err := strconv.Atoi(argsDeviationThresholdTrigger)
			if err != nil {
				return err
			}

			argsInitDataProviderList := strings.Split(strings.TrimSpace(argsInitDataProviderListStr), " ")
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

			msg := types.NewMsgFeed(argsFeedId, feedOwnerAddr, clientCtx.GetFromAddress(), initDataProviderList, uint32(submissionCount), uint32(heartbeatTrigger), uint32(deviationThresholdTrigger))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
