package db

import "github.com/gocql/gocql"

type queryImpl struct {
	query *gocql.Query
}

// Exec executes the query without returning any rows.
func (q queryImpl) Exec() error {
	return q.query.Exec()
}

// Iter executes the query and returns an iterator capable of iterating over
// all results.
func (q queryImpl) Iter() ResultsIterator {
	return resultsIteratorImpl{iter: q.query.Iter()}
}

// Scan executes the query, copies the columns of the first selected row into
// the values pointed at by dest and discards the rest. If no rows were
// selected, ErrNotFound is returned.
func (q queryImpl) Scan(dest ...interface{}) error {
	return q.query.Scan(dest...)
}
