package awesome

import "github.com/gophr-pm/gophr/lib/db"

// ExecBatch executes a batch cassandra query and returns errors via
// an error channel.
func execBatch(
	batch db.Batch,
	resultChan chan error,
) {
	err := batch.Execute()
	resultChan <- err
}
