package types

// NewGenesisState creates a new genesis state instance
func NewGenesisState(params Params) *GenesisState {
	return &GenesisState{
		Params: params,
	}
}

// DefaultGenesisState returns the default epochs genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams())
}

func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}
