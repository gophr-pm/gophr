package github

import "net/http"

// HTTPHeadReq executes an HTTP `HEAD` to the specified URL and returns the
// corresponding response.
type HTTPHeadReq func(url string) (*http.Header, error)

// FetchFullSHAArgs is the arguments struct for packagePusher.
type FetchFullSHAArgs struct {
	Author     string
	Repo       string
	ShortSHA   string
	DoHTTPHead HTTPHeadReq
}
