package db

// Queryable is an interface representing anything that can make and return
// queries against Cassandra.
type Queryable interface {
	// Query generates a new query object for interacting with the database.
	// Further details of the query may be tweaked using the resulting query value
	// before the query is executed. Query is automatically prepared if it has not
	// previously been executed.
	Query(stmt string, values ...interface{}) Query
}
