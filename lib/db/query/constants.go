package query

const (
	// DBKeyspaceName is the name of the gophr cassandra keyspace.
	DBKeyspaceName = "gophr"
	// countOperator is the operator that specifies that rows should be counted.
	countOperator = `count(*)`
	// sumOperatorTemplate is the template for the operator that returns an
	// aggregate of number types.
	sumOperatorTemplate = `sum(%s)`
)
