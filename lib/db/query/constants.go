package query

const (
	// DBProtoVersion is the cassandra protocol version used by gophr.
	DBProtoVersion = 4
	// DBKeyspaceName is the name of the gophr cassandra keyspace.
	DBKeyspaceName = "gophr"

	// countOperator is the operator that specifies that rows should be counted.
	countOperator = `count(*)`
	// Now is the current timestamp.
	Now = `toTimestamp(now())`
)
