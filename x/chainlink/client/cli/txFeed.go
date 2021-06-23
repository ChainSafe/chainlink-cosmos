package cli

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
	"strings"
)

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
