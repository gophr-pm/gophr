package db

// Query represents a CQL statement that can be executed.
type Query interface {
	// Exec executes the query without returning any rows.
	Exec() error
	// Iter executes the query and returns an iterator capable of iterating over
	// all results.
	Iter() ResultsIterator
	// Scan executes the query, copies the columns of the first selected row into
	// the values pointed at by dest and discards the rest. If no rows were
	// selected, ErrNotFound is returned.
	Scan(dest ...interface{}) error
}
