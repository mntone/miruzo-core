package contract

type Capability string

const (
	SupportsInfinityTimestamp     Capability = "supports_infinity_timestamp"
	SupportsLastInsertID          Capability = "supports_last_insert_id"
	SupportsNumberedPlaceholder   Capability = "supports_numbered_placeholder"
	SupportsReturningClause       Capability = "supports_returning_clause"
	SupportsUnnumberedPlaceholder Capability = "supports_unnumbered_placeholder"
)

var allCapabilities = [...]Capability{
	SupportsInfinityTimestamp,
	SupportsLastInsertID,
	SupportsNumberedPlaceholder,
	SupportsReturningClause,
	SupportsUnnumberedPlaceholder,
}

func AllCapabilities() []Capability {
	return allCapabilities[:]
}
