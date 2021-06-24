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
	// TODO: genState.ModuleOwners[0].GetCachedValue() is nil here, need to cast Any to GenesisAccount type and unpack it
	var a interface{}
	genState.ModuleOwners[0].GetCachedValue()
	fmt.Println("/////     ",genState.ModuleOwners[0].Value)

	a = genState.ModuleOwners[0].Value

	acc, ok := a.(authtypes.GenesisAccount)
	if !ok {
		panic("expected genesis account")
	}
	fmt.Println("/////     ",acc)

	accounts, err := authtypes.UnpackAccounts(genState.ModuleOwners)
	if err != nil {
		panic(err)
	}

	for _, a := range accounts {
		m := types.ModuleOwner{
			Address: a.GetAddress(),
			PubKey:  a.GetPubKey().Bytes(),
		}
		k.SetModuleOwner(ctx, &m)
	}
}

// ExportGenesis returns the chainlink module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	//moduleOwners := k.GetModuleOwnerList(ctx)
	//genesis.ModuleOwners = moduleOwners.ModuleOwner

	return genesis
}
