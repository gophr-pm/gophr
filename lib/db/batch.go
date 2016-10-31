package db

// Batch represents a batch of CQL statements that can be executed.
type Batch interface {
	// Query adds the query to the batch operation
	Query(stmt string, values ...interface{})
	// Execute executes the batch operation.
	Execute() error
}
