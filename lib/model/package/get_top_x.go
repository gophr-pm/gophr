package pkg

import (
	"errors"
	"fmt"

	"github.com/gophr-pm/gophr/lib/db/query"
)

// GetTopX (as in "get top ten") gets the top packages sorted descendingly
// within the specified time split.
func GetTopX(q query.Queryable, x int, split TimeSplit) (Summaries, error) {
	if x < 1 {
		return nil, errors.New("X must be greater than zero")
	}

	// Turn the split into a field to sort.
	var sortField string
	switch split {
	case Daily:
		sortField = packagesColumnNameDailyDownloads
	case Weekly:
		sortField = packagesColumnNameWeeklyDownloads
	case Monthly:
		sortField = packagesColumnNameMonthlyDownloads
	case AllTime:
		sortField = packagesColumnNameAllTimeDownloads
	default:
		return nil, errors.New("Invalid time split provided")
	}

	// Create and execute the query, then create in iterator for the results.
	iter := query.
		Select(
			packagesColumnNameRepo,
			packagesColumnNameStars,
			packagesColumnNameAuthor,
			packagesColumnNameAwesome,
			packagesColumnNameDescription,
			packagesColumnNameDailyDownloads,
			packagesColumnNameWeeklyDownloads,
			packagesColumnNameMonthlyDownloads,
			packagesColumnNameAllTimeDownloads).
		From(packagesTableName).
		Where(query.Index(packagesIndexName).Matches(fmt.Sprintf(
			descSortExprTemplate,
			sortField))).
		Limit(x).
		Create(q).
		Iter()

	var (
		summaries   []Summary
		nextSummary Summary
	)

	// Scan into a summary struct. Add it to the list if successful.
	for iter.Scan(
		nextSummary.Repo,
		nextSummary.Stars,
		nextSummary.Author,
		nextSummary.Awesome,
		nextSummary.Description,
		nextSummary.DailyDownloads,
		nextSummary.WeeklyDownloads,
		nextSummary.MonthlyDownloads,
		nextSummary.AllTimeDownloads) {
		summaries = append(summaries, nextSummary)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf(
			`Failed to get top %d packages from the db: %v`,
			x,
			err)
	}

	return summaries, nil
}
