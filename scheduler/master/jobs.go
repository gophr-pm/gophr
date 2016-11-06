package main

type job struct {
	name string
	path string
}

var (
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
)
