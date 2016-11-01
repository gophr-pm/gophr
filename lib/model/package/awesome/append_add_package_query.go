package awesome

import (
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/query"
)

// AppendAddPackageQuery records that an archive of a package version exists.
func AppendAddPackageQuery(
	b db.Batch,
	author string,
	repo string,
) {
	query.InsertInto(tableName).
		Value(columnNameAuthor, author).
		Value(columnNameRepo, repo).
		AppendTo(b)
}
