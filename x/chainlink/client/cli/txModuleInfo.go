// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package cli

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func CmdAddModuleOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-module-owner [address] [publicKey]",
		Short: "Add a chainLink module owner. Signer must be an existing module owner.",
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
		Use:   "module-ownership-transfer [newModuleOwnerAddress] [newModuleOwnerPublicKey]",
		Short: "Transfer chainLink module ownership from an existing module owner account to another account. Signer must be an existing module owner.",
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
