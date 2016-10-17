package godoc

import "github.com/gophr-pm/gophr/lib/db/query"

// AppendAddPackageQuery records that an archive of a package version exists.
func AppendAddPackageQuery(
	session query.VoidQueryable,
	author,
	repo,
	description string,
) {
	query.InsertInto(tableName).
		Value(columnNameAuthor, author).
		Value(columnNameRepo, repo).
		Value(columnNameDescription, description).
		CreateVoid(session)
}
