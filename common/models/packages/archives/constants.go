package archives

const (
	tableName = "package_archive_records"
	// TODO(skeswa): migrate column to "sha" instead of "ref".
	columnNameSHA    = "ref"
	columnNameRepo   = "repo"
	columnNameAuthor = "author"
)
