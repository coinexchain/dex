package types

// GenesisState - all asset state that must be provided at genesis
type GenesisState struct {
	Params             Params   `json:"params"`
	Tokens             []Token  `json:"tokens"`
	Whitelist          []string `json:"whitelist"`
	ForbiddenAddresses []string `json:"forbidden_addresses"`
}

// NewGenesisState - Create a new genesis state
func NewGenesisState(params Params, tokens []Token, whitelist []string, forbiddenAddresses []string) GenesisState {
	return GenesisState{
		Params:             params,
		Tokens:             tokens,
		Whitelist:          whitelist,
		ForbiddenAddresses: forbiddenAddresses,
	}
}

// DefaultGenesisState - Return a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultParams(), []Token{}, []string{}, []string{})
}
