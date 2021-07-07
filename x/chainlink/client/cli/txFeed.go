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
		Use:   "addFeed [feedId] [feedOwnerAddress] [submissionCount] [heartbeatTrigger] [deviationThresholdTrigger] [initDataProviderList]",
		Short: "Add new feed. Signer must be the existing module owner. initDataProviderList is a string contains each data provider's address with pubkey and split by comma",
		Args:  cobra.MinimumNArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFeedId := args[0]
			argsFeedOwnerAddr := args[1]
			argsSubmissionCount := args[2]
			argsHeartbeatTrigger := args[3]
			argsDeviationThresholdTrigger := args[4]
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
