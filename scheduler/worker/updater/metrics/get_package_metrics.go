package metrics

import (
	"fmt"
	"sync"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/lib/db/model/package/download"
	"github.com/gophr-pm/gophr/lib/github"
)

// getPackageMetrics calculates and organizes the metrics for a specific package
// from the database. The result is an args struct for pkg.UpdateMetrics.
func getPackageMetrics(
	q db.Queryable,
	ghSvc github.RequestService,
	summary pkg.Summary,
) (pkg.UpdateMetricsArgs, error) {
	var (
		wg                        sync.WaitGroup
		trendScore                float32
		searchScore               float32
		getSplitsResult           getSplitsWrapperResult
		fetchRepoDataResult       fetchRepoDataWrapperResult
		getVersionDownloadsResult getVersionDownloadsWrapperResult
	)

	wg.Add(3)
	go getSplitsWrapper(getSplitsWrapperArgs{
		q:         q,
		wg:        &wg,
		repo:      summary.Repo,
		author:    summary.Author,
		result:    &getSplitsResult,
		getSplits: download.GetSplits,
	})
	go fetchRepoDataWrapper(fetchRepoDataWrapperArgs{
		wg:            &wg,
		repo:          summary.Repo,
		author:        summary.Author,
		result:        &fetchRepoDataResult,
		fetchRepoData: ghSvc.FetchGitHubDataForPackageModel,
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
	if fetchRepoDataResult.err != nil {
		return pkg.UpdateMetricsArgs{}, fmt.Errorf(
			`Failed to get downloads for package "%s/%s": Failed to get splits: %v`,
			summary.Repo,
			summary.Author,
			fetchRepoDataResult.err)
	}
	if getVersionDownloadsResult.err != nil {
		return pkg.UpdateMetricsArgs{}, fmt.Errorf(
			`Failed to get downloads for package "%s/%s": `+
				`Failed to get version downloads: %v`,
			summary.Repo,
			summary.Author,
			getVersionDownloadsResult.err)
	}

	// Calculate the derived metrics.
	trendScore = pkg.CalcTrendScore(
		getSplitsResult.splits.Daily,
		getSplitsResult.splits.Weekly,
		getSplitsResult.splits.Monthly)
	searchScore = pkg.CalcSearchScore(
		fetchRepoDataResult.repoData.Stars,
		getSplitsResult.splits.AllTime,
		summary.Awesome,
		trendScore)

	return pkg.UpdateMetricsArgs{
		Repo:                    summary.Repo,
		Stars:                   fetchRepoDataResult.repoData.Stars,
		Author:                  summary.Author,
		Queryable:               q,
		TrendScore:              trendScore,
		SearchScore:             searchScore,
		Description:             fetchRepoDataResult.repoData.Description,
		DailyDownloads:          getSplitsResult.splits.Daily,
		WeeklyDownloads:         getSplitsResult.splits.Weekly,
		MonthlyDownloads:        getSplitsResult.splits.Monthly,
		AllTimeDownloads:        getSplitsResult.splits.AllTime,
		AllTimeVersionDownloads: getVersionDownloadsResult.versionDownloads,
	}, nil
}
