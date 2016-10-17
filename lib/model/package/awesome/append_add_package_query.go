package awesome

import "github.com/gophr-pm/gophr/lib/db/query"

// AppendAddPackageQuery records that an archive of a package version exists.
func AppendAddPackageQuery(
	session query.VoidQueryable,
	author string,
	repo string,
) {
	query.InsertInto(tableName).
		Value(columnNameAuthor, author).
		Value(columnNameRepo, repo).
		CreateVoid(session)
}
