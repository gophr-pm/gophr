package apikey

import (
	"fmt"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/query"
)

// InsertionTuple is a tuple of a Github API Key string and some usage metadata.
type InsertionTuple struct {
	Key              string
	ForScheduledJobs bool
}

// InsertAll inserts a collection of Github API Keys into the database.
func InsertAll(q db.BatchingQueryable, tuples []InsertionTuple) error {
	b := q.NewLoggedBatch()

	for _, tuple := range tuples {
		query.InsertInto(tableName).
			Value(columnNameKey, tuple.Key).
			Value(columnNameForScheduledJobs, tuple.ForScheduledJobs).
			AppendTo(b)
	}

	if err := b.Execute(); err != nil {
		return fmt.Errorf("Failed to insert Github API keys: %v.", err)
	}

	return nil
}
