package db

import "github.com/gocql/gocql"

type clientImpl struct {
	session *gocql.Session
}

// Close closes all connections. The client is unusable after this operation.
func (c clientImpl) Close() {
	c.session.Close()
}

// Query generates a new query object for interacting with the database.
// Further details of the query may be tweaked using the resulting query value
// before the query is executed. Query is automatically prepared if it has not
// previously been executed.
func (c clientImpl) Query(stmt string, values ...interface{}) Query {
	return queryImpl{query: c.session.Query(stmt, values...)}
}

// NewLoggedBatch creates a new logged batch.
func (c clientImpl) NewLoggedBatch() Batch {
	return batchImpl{
		batch:   c.session.NewBatch(gocql.LoggedBatch),
		session: c.session,
	}
}

// NewUnloggedBatch creates a new unlogged batch.
func (c clientImpl) NewUnloggedBatch() Batch {
	return batchImpl{
		batch:   c.session.NewBatch(gocql.UnloggedBatch),
		session: c.session,
	}
}
