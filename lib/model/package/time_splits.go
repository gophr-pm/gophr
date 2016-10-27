package pkg

// TimeSplit represents a span of time over which a query can be applied.
type TimeSplit int

const (
	// Daily is the last 24 hours.
	Daily = TimeSplit(iota)
	// Weekly is the last 7 days.
	Weekly
	// Monthly is the last 30 days.
	Monthly
	// AllTime is all-time.
	AllTime
)
