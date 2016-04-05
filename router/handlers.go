package main

import "net/http"

// var (
// 	// The regular expression used to deduce user, repository and version information from an incoming request.
// 	// This expressions expects input matching the following formats:
// 	// - Matches any version:
// 	//   - /user/repo
// 	// - Matches specific version:
// 	//   - /user/repo@1.1.1
// 	//   - /user/repo@1.1.1-alpha.1
// 	// - Matches version range:
// 	//   - /user/repo@1.x
// 	//   - /user/repo@1.1.x
// 	//   - /user/repo@1.1.1-alpha.x
// 	// - Matches greater than or equal to version:
// 	//   - /user/repo@1.1.1+
// 	//   - /user/repo@1.1.1+
// 	// - Matches less than or equal to version:
// 	//   - /user/repo@1.1.1-
// 	// - Matches within patch level of version (carat in npm):
// 	//   - /user/repo@~1.1.1
// 	repoURLRegex = regexp.MustCompile(`\/([a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\/([a-zA-Z0-9\-_]+)(?:@(~?)([0-9]+)\.(?:([0-9]|x))(?:\.([0-9]|x))?(?:\-([a-zA-Z0-9\-_]+)(?:\.([0-9]|x+))?)?([+-]?))?`)
// )

func healthCheckHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("ok"))
	return
}
