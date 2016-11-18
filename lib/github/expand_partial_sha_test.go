package github

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/stretchr/testify/assert"
)

func TestExpandPartialSHA(t *testing.T) {
	fakeRequestService := &requestServiceImpl{
		ddClient: datadog.NewFakeDataDogClient(),
	}

	fullSHA, err := fakeRequestService.ExpandPartialSHA(ExpandPartialSHAArgs{
		Author:   "test",
		Repo:     "testy",
		ShortSHA: "123456",
		DoHTTPHead: func(url string) (*http.Header, error) {
			assert.Equal(t, "https://github.com/test/testy/archive/123456.zip", url)

			header := http.Header{}
			header.Add("Etag", "\"1234567890123456789012345678901234567890\"")
			return &header, nil
		},
	},
	)
	assert.Nil(t, err)
	assert.Equal(t, "1234567890123456789012345678901234567890", fullSHA)

	fullSHA, err = fakeRequestService.ExpandPartialSHA(ExpandPartialSHAArgs{
		Author:   "test",
		Repo:     "testy",
		ShortSHA: "123456",
		DoHTTPHead: func(url string) (*http.Header, error) {
			return &http.Header{}, errors.New("This is an error")
		},
	})
	assert.Equal(t, "", fullSHA)
	assert.NotNil(t, err)

	fullSHA, err = fakeRequestService.ExpandPartialSHA(ExpandPartialSHAArgs{
		Author:   "test",
		Repo:     "testy",
		ShortSHA: "123456",
		DoHTTPHead: func(url string) (*http.Header, error) {
			assert.Equal(t, "https://github.com/test/testy/archive/123456.zip", url)

			header := http.Header{}
			header.Add("Etag", "\"123456\"")
			return &header, nil
		},
	})
	assert.Equal(t, "", fullSHA)
	assert.NotNil(t, err)

	fullSHA, err = fakeRequestService.ExpandPartialSHA(ExpandPartialSHAArgs{
		Author:   "test",
		Repo:     "testy",
		ShortSHA: "123456",
		DoHTTPHead: func(url string) (*http.Header, error) {
			assert.Equal(t, "https://github.com/test/testy/archive/123456.zip", url)

			header := http.Header{}
			header.Add("Etag", "\"ThisisnotashasoyoushouldnotbeabletousethisasaSHA\"")
			return &header, nil
		},
	})
	assert.Equal(t, "", fullSHA)
	assert.NotNil(t, err)
}
