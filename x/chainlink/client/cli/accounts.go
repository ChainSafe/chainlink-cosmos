// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package cli

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddGenesisAccountCmd returns add-genesis-account cobra Command.
func AddChainlinkAccountCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-chainlink-account <chainlink_oracle_public_key> <chainlink_oracle_signing_key> [piggy_cosmos_address]",
		Short: "Add a chainlink account to the store.",
		Long: `Add a chainlink oracle account to the network. The chainlink account will be associated with a
		Cosmos account. The piggyAddress will be set to the submitter's Cosmos account by default.
`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsChainlinkPublicKey := args[0]
			argsChainlinkSigningKey := args[1]
			var piggyAddress sdk.AccAddress

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if len(args[2]) > 0 {
				argsPiggyAddress := args[2]
				piggyAddress, err = sdk.AccAddressFromBech32(argsPiggyAddress)
				if err != nil {
					return err
				}
			} else {
				piggyAddress = clientCtx.GetFromAddress()
			}

			msg := types.NewMsgAccount(
				clientCtx.GetFromAddress(),
				[]byte(argsChainlinkPublicKey),
				[]byte(argsChainlinkSigningKey),
				piggyAddress,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)

		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
