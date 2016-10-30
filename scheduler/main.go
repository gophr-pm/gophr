package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/model"
	"github.com/gophr-pm/gophr/scheduler/indexer/awesome"
	"github.com/gophr-pm/gophr/scheduler/indexer/github"
	"github.com/gophr-pm/gophr/scheduler/indexer/godoc"
	"github.com/robfig/cron"
)

func main() {
	// Instantiate new cron instance.
	c := cron.New()

	// List cron jobs.
	c.AddFunc("0 0 * * * *", func() {
		if err := awesome.Index(awesome.IndexArgs{
			Init:            common.Init,
			DoHTTPGet:       awesome.DoHTTPGet,
			BatchExecutor:   awesome.ExecBatch,
			PackageFetcher:  awesome.FetchAwesomeGoList,
			PersistPackages: awesome.PersistAwesomePackages,
		}); err != nil {
			// TODO(Shikkic): Send error somewhere, possibly deadman's snitch?
			log.Println(err)
		}
	})

	c.AddFunc("0 0 * * * *", godoc.Index)

	c.AddFunc("0 0 0 * * *", func() {
		if err := githubIndexer.Index(githubIndexer.IndexArgs{
			Init:                    common.Init,
			PackageDeleter:          models.DeletePackageModel,
			PackageRetriever:        models.ScanAllPackageModels,
			PackageInserter:         models.InsertPackage,
			NewGithubRequestService: github.NewRequestService,
			RequestTimeBuffer:       50 * time.Millisecond,
		}); err != nil {
			// TODO(Shikkic): Send error somewhere, possibly deadman's snitch?
			log.Println(err)
		}
	})

	// Start the cron process.
	c.Start()

	// Create handlers for generating metrics.
	http.HandleFunc("/daily-downloads", dailyDownloadsHandler)

	// Start HTTP server.
	http.ListenAndServe(":8080", nil)
}

// dailyDownloadsHandler responsible running a metric job via a handler.
func dailyDownloadsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("placeholder")
}
