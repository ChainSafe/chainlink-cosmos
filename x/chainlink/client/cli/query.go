// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: MIT

package cli

import (
	"fmt"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdGetFeedDataByRound())
	cmd.AddCommand(CmdGetLatestFeedData())
	cmd.AddCommand(CmdGetModuleOwnerList())
	cmd.AddCommand(CmdGetFeedInfo())
	cmd.AddCommand(CmdGetAccountInfo())

	return cmd
}
