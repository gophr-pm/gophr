package main

type job struct {
	name string
	path string
}

var (
	updateGithubMetadata = job{
		name: "updateGithubMetadata",
		path: "update/github-metadata",
	}
	updatePackageMetrics = job{
		name: "updatePackageMetrics",
		path: "update/metrics",
	}
	indexAwesomePackages = job{
		name: "indexAwesomePackages",
		path: "index/awesome",
	}
	indexGoSearchPackages = job{
		name: "indexGoSearchPackages",
		path: "index/go-search",
	}
	deleteOldDownloadsPackages = job{
		name: "deleteOldDownloadsPackages",
		path: "delete/old-downloads",
	}
)
