package db

import "github.com/gocql/gocql"

type batchImpl struct {
	batch   *gocql.Batch
	session *gocql.Session
}

// Query adds the query to the batch operation
func (b batchImpl) Query(stmt string, values ...interface{}) {
	b.batch.Query(stmt, values...)
}

// Execute executes the batch operation.
func (b batchImpl) Execute() error {
	return b.session.ExecuteBatch(b.batch)
}
