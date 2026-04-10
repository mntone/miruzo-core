package contract

type Capability string

const (
	SupportsLastInsertID      Capability = "supports_last_insert_id"
	SupportsReturningClause   Capability = "supports_returning_clause"
	SupportsInfinityTimestamp Capability = "supports_infinity_timestamp"
)

var allCapabilities = [...]Capability{
	SupportsLastInsertID,
	SupportsReturningClause,
	SupportsInfinityTimestamp,
}

func AllCapabilities() []Capability {
	return allCapabilities[:]
}
