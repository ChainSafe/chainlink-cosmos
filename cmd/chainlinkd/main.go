package main

import (
	"os"

	"github.com/ChainSafe/chainlink-cosmos/app"
	"github.com/ChainSafe/chainlink-cosmos/cmd/chainlinkd/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
