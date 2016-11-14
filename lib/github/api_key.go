package github

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	httpHeaderResetTime            = "X-RateLimit-Reset"
	httpHeaderRequestsRemaining    = "X-RateLimit-Remaining"
	githubAPIUsageEndpointTemplate = "https://api.github.com/repos/a/b?access_token=%s"
)

// apiKey represents a single Github API keys,
// it's responsible for keeping track of API usage via that key
type apiKey struct {
	token              string
	dataLock           sync.RWMutex
	requestLock        sync.Mutex
	remainingUses      int
	rateLimitResetTime time.Time
}

// newAPIKey creates a new Github API key.
func newAPIKey(token string) (*apiKey, error) {
	newKey := &apiKey{token: token}
	if err := newKey.updateByRequest(); err != nil {
		return nil, err
	}

	return newKey, nil
}

// getFromGithub issues an HTTP GET request against the specified URL, but with
// the authentication of this key.
func (key *apiKey) getFromGithub(url string) (*http.Response, error) {
	// Add the token to the URL differently depending on whether there is already
	// a query string.
	if strings.IndexByte(url, '?') == -1 {
		url = url + "?access_token=" + key.token
	} else {
		url = url + "&access_token=" + key.token
	}

	// This is to make sure that only one request at a time happens on this token.
	key.requestLock.Lock()
	// Make the request, then update accordingly.
	resp, err := http.Get(url)
	// Allow other requests to go through on this key.
	key.requestLock.Unlock()

	// Attempt to update the key.
	if resp != nil {
		key.update(resp.Header)
	}

	return resp, err
}

// update modifies the usage metadata for this API key using the header of a
// Github API response.
func (key *apiKey) update(header http.Header) *apiKey {
	var (
		rawRemainingRequests  = header.Get(httpHeaderRequestsRemaining)
		rawRateLimitResetTime = header.Get(httpHeaderResetTime)
	)

	var (
		remainingRequests, _    = strconv.Atoi(rawRemainingRequests)
		rateLimitResetTimeMS, _ = strconv.ParseInt(rawRateLimitResetTime, 10, 64)
		rateLimitResetTimeStamp = time.Unix(rateLimitResetTimeMS, 0)
	)

	key.dataLock.Lock()
	key.remainingUses = remainingRequests
	key.rateLimitResetTime = rateLimitResetTimeStamp
	key.dataLock.Unlock()

	log.Printf(
		"Recently used Github API key now has %d remaining requests.\n",
		remainingRequests)

	return key
}

// canBeUsed returns true if this key has remaining usages.
func (key *apiKey) canBeUsed() bool {
	key.dataLock.RLock()
	hasRemainingRequests := key.remainingUses < 1
	key.dataLock.RUnlock()

	return hasRemainingRequests
}

// waitUntilUseful blocks until this key can be used.
func (key *apiKey) waitUntilUseful() {
	key.dataLock.RLock()
	resetTime := key.rateLimitResetTime
	sleepTime := resetTime.Sub(time.Now())
	key.dataLock.RUnlock()

	log.Printf("Github API Key is sleeping until %s.\n", sleepTime.String())
	time.Sleep(sleepTime)
}

// updateByRequest updates usage metadata by calling the Github API.
func (key *apiKey) updateByRequest() error {
	resp, err := http.Get(fmt.Sprintf(githubAPIUsageEndpointTemplate, key.token))
	if err != nil {
		return fmt.Errorf(
			"Failed to update key usage metadata by request: %v.",
			err)
	}
	resp.Body.Close()
	key.update(resp.Header)

	return nil
}
