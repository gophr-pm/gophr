package db

// Batchable is an interface representing anything that can make batch Cassandra
// queries.
type Batchable interface {
	NewLoggedBatch() Batch
	NewUnloggedBatch() Batch
}
