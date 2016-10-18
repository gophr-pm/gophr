package awesome

import (
	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/lib/db/query"
)

// ExecBatch executes a batch cassandra query and returns errors via
// an error channel.
func ExecBatch(
	session query.BatchingQueryable,
	batch *gocql.Batch,
	resultChan chan error,
) {
	err := session.ExecuteBatch(batch)
	resultChan <- err
}
