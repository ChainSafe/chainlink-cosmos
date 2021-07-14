package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Chainlink token denom
	LinkDenom string = "link"
)

// NewLinkCoin will create a "link" coin with the provided amount.
func NewLinkCoin(amount sdk.Int) sdk.Coin {
	return sdk.NewCoin(LinkDenom, amount)
}

// NewLinkDecCoin will create a decimal "link" coin with the provided amount.
func NewLinkDecCoin(amount sdk.Int) sdk.DecCoin {
	return sdk.NewDecCoin(LinkDenom, amount)
}

// NewLinkCoinInt64 will create a "link" coin with the given int64 amount.
func NewLinkCoinInt64(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(LinkDenom, amount)
}
