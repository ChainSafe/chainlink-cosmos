package cli

import (
	"errors"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdAddModule() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addModuleOwner [address] [publicKey]",
		Short: "Add ChainLink Module Owner",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsAddress := args[0]
			argsPublicKey := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if clientCtx.GetFromAddress().String() != argsAddress {
				return errors.New("address not match the signer")
			}

			msg := types.NewModuleOwner(clientCtx.GetFromAddress(), []byte(argsPublicKey))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
