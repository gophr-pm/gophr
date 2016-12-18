package github

import (
	"sync"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/model/package"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/scheduler/worker/common"
)

// packageUpdaterArgs is the arguments struct for packageUpdater.
type packageUpdaterArgs struct {
	q         db.Queryable
	wg        *sync.WaitGroup
	errs      chan error
	ghSvc     github.RequestService
	logger    common.JobLogger
	summaries chan pkg.Summary
}

// packageUpdater is a worker for the Update function. It reads incoming
// packages from the summaries channel and updates each package's metrics. If
// any errors are encountered in the process, then they are put into the errors
// channel.
func packageUpdater(args packageUpdaterArgs) {
	// Guarantee that the waitgroup is notified at the end.
	defer args.wg.Done()

	// For each package summary, attempt an update in the database.
	for summary := range args.summaries {
		args.logger.Infof(
			"Now updating package %s/%s\n",
			summary.Author,
			summary.Repo)

		repoData, err := args.ghSvc.FetchRepoData(summary.Author, summary.Repo)
		if err != nil {
			args.errs <- err
			continue
		}

		err = pkg.UpdateGithubMetadata(pkg.UpdateGithubMetadataArgs{
			Repo:        summary.Repo,
			Stars:       repoData.Stars,
			Author:      summary.Author,
			Queryable:   args.q,
			Description: repoData.Description,
		})
		if err != nil {
			args.errs <- err
		}

		args.logger.Infof(
			"Finished updating package %s/%s\n",
			summary.Author,
			summary.Repo)
	}
}
