package cli

import (
	"encoding/json"
	"fmt"

	chainlinktypes "github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
)

func CmdGenesisModuleOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-module-owner [address_or_key_name]",
		Short: "Add init module owner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			address := args[0]
			// TODO: add pubKey support once the UnpackAccounts issue resolved

			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.JSONMarshaler
			cdc := depCdc.(codec.Marshaler)
			serverCtx := server.GetServerContextFromCmd(cmd)
			conf := serverCtx.Config
			conf.SetRoot(clientCtx.HomeDir)

			// checking init module owner address
			addr, err := sdk.AccAddressFromBech32(address)
			if err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}

			baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)
			if err := baseAccount.Validate(); err != nil {
				return fmt.Errorf("failed to validate new genesis account: %w", err)
			}

			// TODO: add keyring support
			//inBuf := bufio.NewReader(cmd.InOrStdin())
			//keyringBackend, err := cmd.Flags().GetString(flags.FlagKeyringBackend)
			//if err != nil {
			//	return err
			//}
			//
			//// attempt to lookup address from Keybase if no address was provided
			//kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.HomeDir, inBuf)
			//if err != nil {
			//	return err
			//}
			//info, err := kb.Key(address)
			//if err != nil {
			//	return fmt.Errorf("failed to get address from Keybase: %w", err)
			//}

			genFile := conf.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			chainLinkGenState := chainlinktypes.GetGenesisStateFromAppState(cdc, appState)

			// check if the new address is already in the genesis
			// TODO: UnpackAccounts not working if chainLinkGenState.ModuleOwners not empty
			// TODO: this is the same issue in InitGenesis func when write init module owner into store.
			accs, err := authtypes.UnpackAccounts(chainLinkGenState.ModuleOwners)
			if err != nil {
				return fmt.Errorf("failed to get accounts from any: %w", err)
			}

			if accs.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}
			accs = append(accs, baseAccount)

			// Add the new account to the set of genesis accounts and sanitize the accounts afterwards.
			accs = authtypes.SanitizeGenesisAccounts(accs)
			genAccs, err := authtypes.PackAccounts(accs)
			if err != nil {
				return fmt.Errorf("failed to convert accounts into any's: %w", err)
			}

			chainLinkGenState.ModuleOwners = genAccs

			chainlinkGenStateBz, err := cdc.MarshalJSON(chainLinkGenState)
			if err != nil {
				return fmt.Errorf("failed to marshal auth genesis state: %w", err)
			}

			appState[chainlinktypes.ModuleName] = chainlinkGenStateBz

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
