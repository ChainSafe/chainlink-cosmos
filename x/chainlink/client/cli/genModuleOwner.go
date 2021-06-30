package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	chainlinktypes "github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
)

func CmdGenesisModuleOwner() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-module-owner [address_or_key_name] [pubKey]",
		Short: "Add init module owner to genesis.json",
		Long:  "Add an init ChainLink module owner account to genesis.json. If a key name is given, the provided account must be in the local Keybase.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			address := args[0]
			pubKey := args[1]

			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.JSONMarshaler
			cdc := depCdc.(codec.Marshaler)
			serverCtx := server.GetServerContextFromCmd(cmd)
			conf := serverCtx.Config
			conf.SetRoot(clientCtx.HomeDir)

			// checking init module owner address
			addr, err := sdk.AccAddressFromBech32(address)
			if err != nil {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				keyringBackend, err := cmd.Flags().GetString(flags.FlagKeyringBackend)
				if err != nil {
					return err
				}

				// attempt to lookup address from Keybase if no address was provided
				kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, clientCtx.HomeDir, inBuf)
				if err != nil {
					return err
				}

				info, err := kb.Key(address)
				if err != nil {
					return fmt.Errorf("failed to get address from Keybase: %w", err)
				}

				addr = info.GetAddress()
			}
			// get bech32 pubkey
			bech32PubKey := sdk.MustGetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, pubKey)

			// address and pubKey must match
			if !bytes.Equal(bech32PubKey.Address().Bytes(), addr.Bytes()) {
				return fmt.Errorf("address and pubKey not match")
			}

			initModuleOwner := chainlinktypes.NewModuleOwner(nil, addr, []byte(pubKey))

			genFile := conf.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			chainLinkGenState := chainlinktypes.GetGenesisStateFromAppState(cdc, appState)

			// check if the new address is already in the genesis
			accs := (chainlinktypes.MsgModuleOwners)(chainLinkGenState.GetModuleOwners())
			if accs.Contains(addr) {
				return fmt.Errorf("cannot add account at existing address %s", addr)
			}
			accs = append(accs, initModuleOwner)

			chainLinkGenState.ModuleOwners = accs

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
