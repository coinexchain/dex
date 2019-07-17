package types

const (
	// RouterKey = "market"
	// StoreKey  = RouterKey
	// Topic     = RouterKey
	// Query
	// ModuleName is the name of the module
	ModuleName = "market"

	// StoreKey is string representation of the store key for asset
	StoreKey = ModuleName

	// RouterKey is the message route for asset
	RouterKey = ModuleName

	// QuerierRoute is the querier route for asset
	QuerierRoute = ModuleName

	DefaultParamspace = ModuleName

	// Kafka topic name
	Topic = ModuleName
)
