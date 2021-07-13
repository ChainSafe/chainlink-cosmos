package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdSubmitFeedData())
	cmd.AddCommand(CmdAddModuleOwner())
	cmd.AddCommand(CmdGenesisModuleOwner())
	cmd.AddCommand(CmdTransferModuleOwnership())
	cmd.AddCommand(CmdAddFeed())
	cmd.AddCommand(CmdAddDataProvider())
	cmd.AddCommand(CmdRemoveDataProvider())
	cmd.AddCommand(CmdSetSubmissionCount())
	cmd.AddCommand(CmdSetHeartbeatTrigger())
	cmd.AddCommand(CmdSetDeviationThreshold())

	return cmd
}
