// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package cli

import (
	"context"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CmdGetAccountInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-chainlink-account",
		Short: "Gets the chainlink account information associated to the submitter",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			params := &types.GetModuleOwnerRequest{}

			res, err := queryClient.GetAllModuleOwner(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func CmdAddChainlinkAccount() *cobra.Command {
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

			msg := types.NewMsgAddAccount(
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

func CmdEditPiggyAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit-piggy-address <piggy_cosmos_address>",
		Short: "Edit the Piggy Address of a Chainlink account.",
		Long: `Update the Piggy Address associated to a Chainlink account. The Piggy Address is the cosmos address that will be used to issue reward distributions in the native token.
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			argsPiggyAddress := args[0]
			piggyAddress, err := sdk.AccAddressFromBech32(argsPiggyAddress)
			if err != nil {
				return err
			}

			msg := types.NewMsgEditAccount(
				clientCtx.GetFromAddress(),
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
