package query

const (
	// DBKeyspaceName is the name of the gophr cassandra keyspace.
	DBKeyspaceName = "gophr"

	// countOperator is the operator that specifies that rows should be counted.
	countOperator = `count(*)`
	// Now is the current timestamp.
	Now = `toTimestamp(now())`
)
