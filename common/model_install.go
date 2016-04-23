package main

import "fmt"

// Constants directly related to interacting with the install models in the
// cassandra database.
const (
	// TableNameAllTimeInstallTotals is the name of the table containing the
	// install total model.
	TableNameAllTimeInstallTotals               = "all_time_install_totals"
	ColumnNameAllTimeInstallTotalsTotal         = "total"
	ColumnNameAllTimeInstallTotalsPackageRepo   = "package_repo"
	ColumnNameAllTimeInstallTotalsPackageAuthor = "package_author"
	// TableNameRangedInstallTotals is the name of the table containing the
	// install total model.
	TableNameRangedInstallTotals               = "ranged_install_totals"
	ColumnNameRangedInstallTotalsDate          = "date"
	ColumnNameRangedInstallTotalsDailyTotal    = "daily_total"
	ColumnNameRangedInstallTotalsWeeklyTotal   = "weekly_total"
	ColumnNameRangedInstallTotalsPackageRepo   = "package_repo"
	ColumnNameRangedInstallTotalsMonthlyTotal  = "monthly_total"
	ColumnNameRangedInstallTotalsPackageAuthor = "package_author"
	// TableNameRangedVersionedInstall is the name of the table containing the
	// install total model.
	TableNameRangedVersionedInstallTotals          = "ranged_versioned_install_totals"
	ColumnNameRangedVersionedInstallDate           = "date"
	ColumnNameRangedVersionedInstallDailyTotal     = "daily_total"
	ColumnNameRangedVersionedInstallWeeklyTotal    = "weekly_total"
	ColumnNameRangedVersionedInstallPackageRepo    = "package_repo"
	ColumnNameRangedVersionedInstallMonthlyTotal   = "monthly_total"
	ColumnNameRangedVersionedInstallPackageAuthor  = "package_author"
	ColumnNameRangedVersionedInstallPackageVersion = "package_version"
)

const (
	cqlQueryReadAllTimeInstallTotal = fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s = ? AND %s = ?`,
		ColumnNameAllTimeInstallTotalsTotal,
		TableNameAllTimeInstallTotals,
		ColumnNameAllTimeInstallTotalsPackageAuthor,
		ColumnNameAllTimeInstallTotalsPackageRepo,
	)

	cqlQueryUpdateAllTimeInstallTotal = fmt.Sprintf(
		`UPDATE %s SET %s = ? WHERE %s = ? AND %s = ?`,
		TableNameAllTimeInstallTotals,
		ColumnNameAllTimeInstallTotalsTotal,
		ColumnNameAllTimeInstallTotalsPackageAuthor,
		ColumnNameAllTimeInstallTotalsPackageRepo,
	)

	cqlQueryReadRangedInstallTotal = fmt.Sprintf(
		`SELECT %s, %s, %s FROM %s WHERE %s = ? AND %s = ?`,
		ColumnNameRangedInstallTotalsDailyTotal,
		ColumnNameRangedInstallTotalsWeeklyTotal,
		ColumnNameRangedInstallTotalsMonthlyTotal,
		TableNameRangedInstallTotals,
		ColumnNameRangedInstallTotalsPackageAuthor,
		ColumnNameRangedInstallTotalsPackageRepo,
	)

	cqlQueryUpdateRangedInstallTotal = fmt.Sprintf(
		`UPDATE %s SET %s = ?, %s = ?, %s = ? WHERE %s = ? AND %s = ? AND %s = ?`,
		TableNameRangedInstallTotals,
		ColumnNameRangedInstallTotalsDailyTotal,
		ColumnNameRangedInstallTotalsWeeklyTotal,
		ColumnNameRangedInstallTotalsMonthlyTotal,
		ColumnNameRangedInstallTotalsPackageAuthor,
		ColumnNameRangedInstallTotalsPackageRepo,
		ColumnNameRangedVersionedInstallDate,
	)

	cqlQueryReadRangedVersionedInstallTotal = fmt.Sprintf(
		`SELECT %s, %s, %s FROM %s WHERE %s = ? AND %s = ? AND %s = ?`,
		ColumnNameRangedVersionedInstallDailyTotal,
		ColumnNameRangedVersionedInstallWeeklyTotal,
		ColumnNameRangedVersionedInstallMonthlyTotal,
		TableNameRangedVersionedInstallTotals,
		ColumnNameRangedVersionedInstallPackageAuthor,
		ColumnNameRangedVersionedInstallPackageRepo,
		ColumnNameRangedVersionedInstallPackageVersion,
	)

	cqlQueryUpdateRangedVersionedInstallTotal = fmt.Sprintf(
		`UPDATE %s SET %s = ?, %s = ?, %s = ? WHERE %s = ? AND %s = ? AND %s = ? AND %s = ?`,
		TableNameRangedVersionedInstallTotals,
		ColumnNameRangedVersionedInstallTotalsDailyTotal,
		ColumnNameRangedVersionedInstallTotalsWeeklyTotal,
		ColumnNameRangedVersionedInstallTotalsMonthlyTotal,
		ColumnNameRangedVersionedInstallTotalsPackageAuthor,
		ColumnNameRangedVersionedInstallTotalsPackageRepo,
		ColumnNameRangedVersionedInstallPackageVersion,
		ColumnNameRangedVersionedInstallDate,
	)
)

type InstallModel struct {
}

type VersionedInstallModel struct {
}

func ReadInstallTotal(packageAuthor string, packageRepo string) (uint64, error) {

}

func BumpInstallTotal(packageAuthor string, packageRepo string) error {

}
