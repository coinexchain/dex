package types

// query endpoints supported by the asset Querier
const (
	QueryToken           = "token-info"
	QueryTokenList       = "token-list"
	QueryWhitelist       = "token-whitelist"
	QueryForbiddenAddr   = "addr-forbidden"
	QueryReservedSymbols = "reserved-symbols"
	QueryParameters      = "parameters"
)

// QueryTokenParams defines the params for query: "custom/asset/token-info"
type QueryTokenParams struct {
	Symbol string
}

func NewQueryAssetParams(s string) QueryTokenParams {
	return QueryTokenParams{
		Symbol: s,
	}
}

// QueryWhitelistParams defines the params for query: "custom/asset/token-whitelist"
type QueryWhitelistParams struct {
	Symbol string
}

func NewQueryWhitelistParams(s string) QueryWhitelistParams {
	return QueryWhitelistParams{
		Symbol: s,
	}
}

// QueryForbiddenAddrParams defines the params for query: "custom/asset/addr-forbidden"
type QueryForbiddenAddrParams struct {
	Symbol string
}

func NewQueryForbiddenAddrParams(s string) QueryForbiddenAddrParams {
	return QueryForbiddenAddrParams{
		Symbol: s,
	}
}
