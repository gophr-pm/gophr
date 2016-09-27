package github

import (
	"errors"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchFullSHAFromPartialSHA(t *testing.T) {
	fullSHA, err := FetchFullSHAFromPartialSHA(
		FetchFullSHAArgs{
			Author:   "test",
			Repo:     "testy",
			ShortSHA: "123456",
			DoHTTPHead: func(url string) (*http.Header, error) {
				assert.Equal(t, "https://github.com/test/testy/archive/123456.zip", url)

				header := http.Header{}
				header.Add("Etag", "\"1234567890123456789012345678901234567890\"")
				log.Println(header.Get("Etag"))
				return &header, nil
			},
		},
	)
	assert.Nil(t, err)
	assert.Equal(t, "1234567890123456789012345678901234567890", fullSHA)

	fullSHA, err = FetchFullSHAFromPartialSHA(
		FetchFullSHAArgs{
			Author:   "test",
			Repo:     "testy",
			ShortSHA: "123456",
			DoHTTPHead: func(url string) (*http.Header, error) {
				return &http.Header{}, errors.New("This is an error")
			},
		},
	)
	assert.Equal(t, "", fullSHA)
	assert.NotNil(t, err)

	fullSHA, err = FetchFullSHAFromPartialSHA(
		FetchFullSHAArgs{
			Author:   "test",
			Repo:     "testy",
			ShortSHA: "123456",
			DoHTTPHead: func(url string) (*http.Header, error) {
				assert.Equal(t, "https://github.com/test/testy/archive/123456.zip", url)

				header := http.Header{}
				header.Add("Etag", "\"\"")
				log.Println(header.Get("Etag"))
				return &header, nil
			},
		},
	)
	assert.Equal(t, "", fullSHA)
	assert.NotNil(t, err)
}
