package chainlink

import (
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/keeper"
	"github.com/ChainSafe/chainlink-cosmos/x/chainlink/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the chainlink module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, owner := range genState.GetModuleOwners() {
		m := types.MsgModuleOwner{
			Address: owner.GetAddress(),
			PubKey:  owner.GetPubKey(),
		}
		k.SetModuleOwner(ctx, &m)
	}
}

// ExportGenesis returns the chainlink module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	moduleOwners := k.GetModuleOwnerList(ctx)
	genesis.ModuleOwners = moduleOwners.GetModuleOwner()

	return genesis
}
