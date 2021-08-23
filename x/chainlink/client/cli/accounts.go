// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package cli

import (
	"context"
	"fmt"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CmdAddChainlinkAccount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-chainlink-account <chainlink_oracle_public_key> <chainlink_oracle_signing_key> [piggy_cosmos_address]",
		Short: "Add a chainlink account to the store.",
		Long: `Add a chainlink oracle account to the network. The chainlink account will be associated with a
		Cosmos account. The piggyAddress will be set to the submitter's Cosmos account by default.
`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(args)
			argsChainlinkPublicKey := args[0]
			argsChainlinkSigningKey := args[1]
			var piggyAddress sdk.AccAddress

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if len(args) > 2 {
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

func CmdGetAccountInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-account-info <cosmos address>",
		Short: "Gets the Chainlink account information.",
		Long: `Retreives the Chainlink account information associated with the sender's Cosmos account address. Optional cosmos address can be provided as an argument to look up. Default will retrieve the account associated with the FromAddress.
		`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadPersistentCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			accountString := args[0]
			account, err := sdk.AccAddressFromBech32(accountString)
			if err != nil {
				return err
			}

			req := &types.GetAccountRequest{AccountAddress: account}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.GetAccountInfo(context.Background(), req)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
