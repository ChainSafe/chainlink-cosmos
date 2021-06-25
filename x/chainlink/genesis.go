package chainlink

import (
	"fmt"

	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/keeper"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// InitGenesis initializes the chainlink module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// TODO: ideally, genState.ModuleOwners[0].GetCachedValue() is nil here, need to cast Any to GenesisAccount type and unpack it
	// TODO: but for some unknown reasons, the UnpackAccounts not working here
	//accounts, err := authtypes.UnpackAccounts(genState.GetModuleOwners())
	//if err != nil {
	//	panic(err)
	//}

	addr, err := sdk.AccAddressFromBech32(string(genState.ModuleOwners[0].Value)[2:])
	if err != nil {
		panic("invalid init module owner address")
	}

	m := types.ModuleOwner{
		Address: addr,
		PubKey:  nil,
	}
	k.SetModuleOwner(ctx, &m)
}

// ExportGenesis returns the chainlink module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	moduleOwners := k.GetModuleOwnerList(ctx)

	baseAccounts := make(authtypes.GenesisAccounts, len(moduleOwners.ModuleOwner))
	for _, owner := range moduleOwners.GetModuleOwner() {
		baseAccount := authtypes.NewBaseAccount(owner.GetAddress(), nil, 0, 0)
		baseAccounts = append(baseAccounts, baseAccount)
	}

	accs := authtypes.SanitizeGenesisAccounts(baseAccounts)
	genAccs, err := authtypes.PackAccounts(accs)
	if err != nil {
		ctx.Logger().With("chainLink_export_genesis", fmt.Errorf("failed to convert accounts into any's: %w", err))
		return genesis
	}

	genesis.ModuleOwners = genAccs

	return genesis
}
