package github

import "net/http"

const (
	requestTypeHEAD = "HEAD"
)

// ParseStarCount TODO Won't need this after implementing FFJSON.
func ParseStarCount(responseBody map[string]interface{}) int {
	starCount := responseBody["stargazers_count"]
	if starCount == nil {
		return 0
	}

	return int(starCount.(float64))
}

// DoHTTPHeadReq makes a HEAD request and returns the corresponding
// response header.
func DoHTTPHeadReq(url string) (*http.Header, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		requestTypeHEAD,
		url,
		nil,
	)
	if err != nil {
		return &http.Header{}, err
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode == 404 {
		return &http.Header{}, err
	}

	return &resp.Header, nil
}
