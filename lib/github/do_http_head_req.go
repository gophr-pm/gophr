package github

import "net/http"

// TODO(skeswa): this should be migrated to lib/http @Shikkic.

// DoHTTPHeadReq makes a HEAD request and returns the corresponding
// response header.
func DoHTTPHeadReq(url string) (*http.Header, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		http.MethodHead,
		url,
		nil)
	if err != nil {
		return &http.Header{}, err
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode == 404 {
		return &http.Header{}, err
	}

	return &resp.Header, nil
}
