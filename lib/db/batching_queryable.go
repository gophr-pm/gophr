package db

// BatchingQueryable is an interface representing anything that can make batch
// Cassandra queries, but also execute individual ones.
type BatchingQueryable interface {
	Batchable
	Queryable
}
