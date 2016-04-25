package common

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

// Constants directly related to interacting with the install models in the
// cassandra database.
const (
	// TableNameAllTimeInstallTotals is the name of the table containing the
	// install total model.
	TableNameAllTimeInstallTotals                    = "all_time_install_totals"
	IndexNameAllTimeInstallTotals                    = "all_time_install_totals_index"
	ColumnNameAllTimeInstallTotalsTotal              = "total"
	ColumnNameAllTimeInstallTotalsPackageRepo        = "package_repo"
	ColumnNameAllTimeInstallTotalsPackageAuthor      = "package_author"
	ColumnNameAllTimeInstallTotalsPackageDescription = "package_description"
	// TableNameRangedInstallTotals is the name of the table containing the
	// install total model.
	TableNameRangedInstallTotals                    = "ranged_install_totals"
	ColumnNameRangedInstallTotalsDate               = "date"
	ColumnNameRangedInstallTotalsDailyTotal         = "daily_total"
	ColumnNameRangedInstallTotalsWeeklyTotal        = "weekly_total"
	ColumnNameRangedInstallTotalsPackageRepo        = "package_repo"
	ColumnNameRangedInstallTotalsMonthlyTotal       = "monthly_total"
	ColumnNameRangedInstallTotalsPackageAuthor      = "package_author"
	ColumnNameRangedInstallTotalsPackageDescription = "package_description"
	// TableNameRangedVersionedInstall is the name of the table containing the
	// install total model.
	TableNameRangedVersionedInstallTotals                = "ranged_versioned_install_totals"
	ColumnNameRangedVersionedInstallTotalsDate           = "date"
	ColumnNameRangedVersionedInstallTotalsDailyTotal     = "daily_total"
	ColumnNameRangedVersionedInstallTotalsWeeklyTotal    = "weekly_total"
	ColumnNameRangedVersionedInstallTotalsPackageRepo    = "package_repo"
	ColumnNameRangedVersionedInstallTotalsMonthlyTotal   = "monthly_total"
	ColumnNameRangedVersionedInstallTotalsPackageAuthor  = "package_author"
	ColumnNameRangedVersionedInstallTotalsPackageVersion = "package_version"
)

const (
	timeRangeDayLength   = time.Hour * 24
	timeRangeWeekLength  = 7
	timeRangeMonthLength = 30
)

var (
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

	cqlQueryGetAllTimeInstallTotalTopTen = fmt.Sprintf(
		`SELECT %s, %s, %s, %s FROM %s WHERE expr(%s, '{
      sort: {
        fields: [
          { field: "%s", reverse: true }
        ]
      }
    }') LIMIT 10`,
		ColumnNameAllTimeInstallTotalsPackageAuthor,
		ColumnNameAllTimeInstallTotalsPackageRepo,
		ColumnNameAllTimeInstallTotalsPackageDescription,
		ColumnNameAllTimeInstallTotalsTotal,
		TableNameAllTimeInstallTotals,
		IndexNameAllTimeInstallTotals,
		ColumnNameAllTimeInstallTotalsTotal,
	)

	cqlQueryReadRangedInstallTotal = fmt.Sprintf(
		`SELECT %s, %s, %s FROM %s WHERE %s = ? AND %s = ? AND %s = ?`,
		ColumnNameRangedInstallTotalsDailyTotal,
		ColumnNameRangedInstallTotalsWeeklyTotal,
		ColumnNameRangedInstallTotalsMonthlyTotal,
		TableNameRangedInstallTotals,
		ColumnNameRangedInstallTotalsDate,
		ColumnNameRangedInstallTotalsPackageAuthor,
		ColumnNameRangedInstallTotalsPackageRepo,
	)

	cqlQueryUpdateRangedInstallTotal = fmt.Sprintf(
		`UPDATE %s SET %s = ?, %s = ?, %s = ? WHERE %s = ? AND %s = ? AND %s = ?`,
		TableNameRangedInstallTotals,
		ColumnNameRangedInstallTotalsDailyTotal,
		ColumnNameRangedInstallTotalsWeeklyTotal,
		ColumnNameRangedInstallTotalsMonthlyTotal,
		ColumnNameRangedInstallTotalsDate,
		ColumnNameRangedInstallTotalsPackageAuthor,
		ColumnNameRangedInstallTotalsPackageRepo,
	)

	cqlQueryReadRangedVersionedInstallTotal = fmt.Sprintf(
		`SELECT %s, %s, %s FROM %s WHERE %s = ? AND %s = ? AND %s = ?`,
		ColumnNameRangedVersionedInstallTotalsDailyTotal,
		ColumnNameRangedVersionedInstallTotalsWeeklyTotal,
		ColumnNameRangedVersionedInstallTotalsMonthlyTotal,
		TableNameRangedVersionedInstallTotals,
		ColumnNameRangedVersionedInstallTotalsPackageAuthor,
		ColumnNameRangedVersionedInstallTotalsPackageRepo,
		ColumnNameRangedVersionedInstallTotalsPackageVersion,
	)

	cqlQueryUpdateRangedVersionedInstallTotal = fmt.Sprintf(
		`UPDATE %s SET %s = ?, %s = ?, %s = ? WHERE %s = ? AND %s = ? AND %s = ? AND %s = ?`,
		TableNameRangedVersionedInstallTotals,
		ColumnNameRangedVersionedInstallTotalsDailyTotal,
		ColumnNameRangedVersionedInstallTotalsWeeklyTotal,
		ColumnNameRangedVersionedInstallTotalsMonthlyTotal,
		ColumnNameRangedVersionedInstallTotalsPackageAuthor,
		ColumnNameRangedVersionedInstallTotalsPackageRepo,
		ColumnNameRangedVersionedInstallTotalsPackageVersion,
		ColumnNameRangedVersionedInstallTotalsDate,
	)
)

type AllTimeInstallTotalModel struct {
	PackageAuthor      string
	PackageRepo        string
	PackageDescription string
	TotalInstalls      int64
}

type RangedInstallTotalModel struct {
	Date                 time.Time
	PackageAuthor        string
	PackageRepo          string
	PackageDescription   string
	TotalDailyInstalls   int64
	TotalWeeklyInstalls  int64
	TotalMonthlyInstalls int64
}

func ReadAllTimeInstallTotal(
	session *gocql.Session,
	packageAuthor string,
	packageRepo string,
) (int64, error) {
	var (
		err   error
		total int64

		iter = session.Query(
			cqlQueryReadAllTimeInstallTotal,
			packageAuthor,
			packageRepo,
		).Iter()
	)

	if iter.Scan(&total) {
		return total, nil
	}

	if err = iter.Close(); err != nil {
		return 0, NewQueryScanError(nil, err)
	}

	return 0, nil
}

func UpdateAllTimeInstallTotal(
	session *gocql.Session,
	packageAuthor string,
	packageRepo string,
	newTotal int64,
) error {
	return session.Query(
		cqlQueryUpdateAllTimeInstallTotal,
		newTotal,
		packageAuthor,
		packageRepo,
	).Exec()
}

func GetAllTimeInstallTotalTopTen(session *gocql.Session) ([]AllTimeInstallTotalModel, error) {
	var (
		err                error
		packageRepo        string
		packageAuthor      string
		totalInstalls      int64
		packageDescription string

		installModels = make([]AllTimeInstallTotalModel, 0, 10)
		iter          = session.Query(cqlQueryGetAllTimeInstallTotalTopTen).Iter()
	)

	for iter.Scan(&packageAuthor, &packageRepo, &packageDescription, &totalInstalls) {
		installModels = append(installModels, AllTimeInstallTotalModel{
			PackageRepo:        packageRepo,
			PackageAuthor:      packageAuthor,
			TotalInstalls:      totalInstalls,
			PackageDescription: packageDescription,
		})
	}

	if err = iter.Close(); err != nil {
		return nil, NewQueryScanError(nil, err)
	}

	return installModels, nil
}

func ReadRangedInstallTotal(
	session *gocql.Session,
	date time.Time,
	packageAuthor string,
	packageRepo string,
) (RangedInstallTotalModel, error) {
	var (
		err                    error
		installsToday          int64
		installsInTheLastWeek  int64
		installsInTheLastMonth int64

		iter = session.Query(
			cqlQueryReadRangedInstallTotal,
			date,
			packageAuthor,
			packageRepo,
		).Iter()
	)

	if iter.Scan(&installsToday, &installsInTheLastWeek, &installsInTheLastMonth) {
		return RangedInstallTotalModel{
			Date:                 date,
			PackageRepo:          packageRepo,
			PackageAuthor:        packageAuthor,
			TotalDailyInstalls:   installsToday,
			TotalWeeklyInstalls:  installsInTheLastWeek,
			TotalMonthlyInstalls: installsInTheLastMonth,
		}, nil
	}

	if err = iter.Close(); err != nil {
		return RangedInstallTotalModel{}, NewQueryScanError(nil, err)
	}

	return RangedInstallTotalModel{
		Date:                 date,
		PackageRepo:          packageRepo,
		PackageAuthor:        packageAuthor,
		TotalDailyInstalls:   0,
		TotalWeeklyInstalls:  0,
		TotalMonthlyInstalls: 0,
	}, nil
}

func UpdateRangedInstallTotals(
	session *gocql.Session,
	installModels []*RangedInstallTotalModel,
) error {
	batch := gocql.NewBatch(gocql.LoggedBatch)

	for _, installModel := range installModels {
		batch.Query(
			cqlQueryUpdateRangedInstallTotal,
			installModel.TotalDailyInstalls,
			installModel.TotalWeeklyInstalls,
			installModel.TotalMonthlyInstalls,
			installModel.Date,
			installModel.PackageAuthor,
			installModel.PackageRepo,
		)
	}

	return session.ExecuteBatch(batch)
}

func BumpRangedInstallTotals(
	session *gocql.Session,
	date time.Time,
	packageAuthor string,
	packageRepo string,
	bumpAmount int64,
) error {
	// Generate a list of dates for the next 30 days (incl. this one).
	dates := make([]time.Time, timeRangeMonthLength)
	startDate := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0,
		0,
		0,
		0,
		date.Location(),
	)
	currDate := startDate
	for i := 0; i < timeRangeMonthLength; i++ {
		// Clone the curr date so array elements aren't affected.
		dates[i] = currDate
		// Increment curr date by one full day.
		currDate = currDate.Add(timeRangeDayLength)
	}

	// For each of those dates, spin off a worker to go get the current totals.
	selectOutputChan := make(chan *RangedInstallTotalModel)
	for _, date := range dates {
		go (func(
			session *gocql.Session,
			date time.Time,
			packageAuthor string,
			packageRepo string,
			selectOutputChan chan *RangedInstallTotalModel) {
			// Get get data from the database for this date.
			model, err := ReadRangedInstallTotal(session, date, packageAuthor, packageRepo)
			if err != nil {
				//TODO (Sandile): log these errors.
				// If we bumped into an error, send zeroes back through the chan.
				selectOutputChan <- &RangedInstallTotalModel{
					Date:                 date,
					PackageRepo:          packageRepo,
					PackageAuthor:        packageAuthor,
					TotalDailyInstalls:   0,
					TotalWeeklyInstalls:  0,
					TotalMonthlyInstalls: 0,
				}
			} else {
				// Ship the result back up the chan.
				selectOutputChan <- &model
			}
		})(session, date, packageAuthor, packageRepo, selectOutputChan)
	}

	// Collect the output from the output chan into a slice.
	installModels := make([]*RangedInstallTotalModel, timeRangeMonthLength)
	for i := 0; i < timeRangeMonthLength; i++ {
		installModels[i] = <-selectOutputChan
	}

	// Kill the chan. once its exhausted.
	close(selectOutputChan)

	// Transform the totals on the install models we collected.
	dayDateBoundary := startDate.Add(timeRangeDayLength)
	weekDateBoundary := startDate.Add(timeRangeDayLength * timeRangeWeekLength)
	monthDateBoundary := startDate.Add(timeRangeDayLength * timeRangeMonthLength)
	for _, installModel := range installModels {
		switch {
		case installModel.Date.Before(dayDateBoundary):
			installModel.TotalDailyInstalls = installModel.TotalDailyInstalls + bumpAmount
			installModel.TotalWeeklyInstalls = installModel.TotalWeeklyInstalls + bumpAmount
			installModel.TotalMonthlyInstalls = installModel.TotalMonthlyInstalls + bumpAmount
		case installModel.Date.Before(weekDateBoundary):
			installModel.TotalWeeklyInstalls = installModel.TotalWeeklyInstalls + bumpAmount
			installModel.TotalMonthlyInstalls = installModel.TotalMonthlyInstalls + bumpAmount
		case installModel.Date.Before(monthDateBoundary):
			installModel.TotalMonthlyInstalls = installModel.TotalMonthlyInstalls + bumpAmount
		}
	}

	// Run the insert on the transformed install models.
	return UpdateRangedInstallTotals(session, installModels)
}
