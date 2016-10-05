package query

import "github.com/gocql/gocql"

// Queryable is an interface representing anything that can make and return
// queries against Cassandra.
type Queryable interface {
	Query(stmt string, values ...interface{}) *gocql.Query
}

// VoidQueryable is an interface representing anything that can make queries
// against Cassandra, but not return anything.
type VoidQueryable interface {
	Query(stmt string, values ...interface{})
}

// Batchable is an interface representing anything that can make batch Cassandra
// queries.
type Batchable interface {
	NewBatch(gocql.BatchType) *gocql.Batch
	ExecuteBatch(*gocql.Batch) error
}

// BatchingQueryable is an interface representing anything that can make batch
// Cassandra queries, but also execute individual ones.
type BatchingQueryable interface {
	Queryable
	Batchable
}
