package pkg

import (
	"time"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/query"
)

// UpdateMetricsArgs is the arguments struct for UpdateMetrics.
type UpdateMetricsArgs struct {
	Repo                    string
	Stars                   int
	Author                  string
	Awesome                 bool
	Queryable               db.Queryable
	TrendScore              float32
	SearchScore             float32
	Description             string
	DailyDownloads          int
	WeeklyDownloads         int
	MonthlyDownloads        int
	AllTimeDownloads        int
	AllTimeVersionDownloads map[string]int
}

// UpdateMetrics updates all of the metrics for a package.
func UpdateMetrics(args UpdateMetricsArgs) error {
	return query.
		Update(packagesTableName).
		Set(packagesColumnNameStars, args.Stars).
		Set(packagesColumnNameAwesome, args.Awesome).
		Set(packagesColumnNameTrendScore, args.TrendScore).
		Set(packagesColumnNameSearchScore, args.SearchScore).
		Set(packagesColumnNameDescription, args.Description).
		Set(packagesColumnNameDailyDownloads, args.DailyDownloads).
		Set(packagesColumnNameDateLastIndexed, time.Now()).
		Set(packagesColumnNameWeeklyDownloads, args.WeeklyDownloads).
		Set(packagesColumnNameMonthlyDownloads, args.MonthlyDownloads).
		Set(packagesColumnNameAllTimeDownloads, args.AllTimeDownloads).
		Set(
			packagesColumnNameAllTimeVersionDownloads,
			args.AllTimeVersionDownloads).
		Where(query.Column(packagesColumnNameRepo).Equals(args.Repo)).
		And(query.Column(packagesColumnNameAuthor).Equals(args.Author)).
		IfExists().
		Create(args.Queryable).
		Exec()
}
