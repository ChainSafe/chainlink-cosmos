package types

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// GenesisState defines the chainlink module genesis state
type GenesisState struct {
	Accounts []GenesisAccount `json:"accounts"`
}

type GenesisAccount struct {
	Address string `json:"address"`
	Code    string `json:"code,omitempty"`
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # ibc/genesistype/default
		// this line is used by starport scaffolding # genesis/types/default
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # ibc/genesistype/validate

	// this line is used by starport scaffolding # genesis/types/validate

	return nil
}
