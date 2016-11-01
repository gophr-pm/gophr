package awesome

import (
	"io/ioutil"
	"net/http"
)

// DoHTTPGet makes a GET request and returns the corresponding
// response header.
func DoHTTPGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil || (resp.StatusCode < 200 || resp.StatusCode > 300) {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
