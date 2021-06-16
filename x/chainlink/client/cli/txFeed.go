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
		Use:   "submit-feedData [feedId] [feedData] [signatures]",
		Short: "Submit feed data",
		Args:  cobra.MinimumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsFeedId := args[0]
			argsFeedData := args[1]
			argsSignatures := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

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
