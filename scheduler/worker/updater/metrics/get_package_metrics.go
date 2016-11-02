package metrics

import (
	"fmt"
	"sync"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/lib/db/model/package/download"
)

// getPackageMetrics calculates and organizes the metrics for a specific package
// from the database. The result is an args struct for pkg.UpdateMetrics.
func getPackageMetrics(
	q db.Queryable,
	summary pkg.Summary,
) (pkg.UpdateMetricsArgs, error) {
	var (
		wg                        sync.WaitGroup
		trendScore                float32
		searchScore               float32
		getSplitsResult           getSplitsWrapperResult
		getVersionDownloadsResult getVersionDownloadsWrapperResult
	)

	wg.Add(2)
	go getSplitsWrapper(getSplitsWrapperArgs{
		q:         q,
		wg:        &wg,
		repo:      summary.Repo,
		author:    summary.Author,
		result:    &getSplitsResult,
		getSplits: download.GetSplits,
	})
	go getVersionDownloadsWrapper(getVersionDownloadsWrapperArgs{
		q:                   q,
		wg:                  &wg,
		repo:                summary.Repo,
		author:              summary.Author,
		result:              &getVersionDownloadsResult,
		fetchRefs:           lib.FetchRefs,
		getVersionDownloads: download.GetForVersions,
	})

	wg.Wait()
	if getSplitsResult.err != nil {
		return pkg.UpdateMetricsArgs{}, fmt.Errorf(
			`Failed to get downloads for package "%s/%s": Failed to get splits: %v`,
			summary.Repo,
			summary.Author,
			getSplitsResult.err)
	}
	if getVersionDownloadsResult.err != nil {
		return pkg.UpdateMetricsArgs{}, fmt.Errorf(
			`Failed to get downloads for package "%s/%s": `+
				`Failed to get version downloads: %v`,
			summary.Repo,
			summary.Author,
			getVersionDownloadsResult.err)
	}

	trendScore = pkg.CalcTrendScore(
		getSplitsResult.splits.Daily,
		getSplitsResult.splits.Weekly,
		getSplitsResult.splits.Monthly)
	searchScore = pkg.CalcSearchScore(
		summary.Stars,
		getSplitsResult.splits.AllTime,
		summary.Awesome,
		trendScore)

	return pkg.UpdateMetricsArgs{
		Repo:                    summary.Repo,
		Author:                  summary.Author,
		Queryable:               q,
		TrendScore:              trendScore,
		SearchScore:             searchScore,
		DailyDownloads:          getSplitsResult.splits.Daily,
		WeeklyDownloads:         getSplitsResult.splits.Weekly,
		MonthlyDownloads:        getSplitsResult.splits.Monthly,
		AllTimeDownloads:        getSplitsResult.splits.AllTime,
		AllTimeVersionDownloads: getVersionDownloadsResult.versionDownloads,
	}, nil
}
