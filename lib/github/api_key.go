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
	defer key.requestLock.Unlock()

	// Make the request, then update accordingly.
	resp, err := http.Get(url)
	if resp != nil {
		key.update(resp.Header)
	}

	return resp, err
}

// TODO:(Shikkic) consider revising how we parse the values from header
func (key *apiKey) update(header http.Header) *apiKey {
	// TODO(skeswa): spell check.
	remaingRequests := header.Get(httpHeaderRequestsRemaining)
	rateLimitResetTime := header.Get(httpHeaderResetTime)

	remainingRequestsInt, _ := strconv.Atoi(remaingRequests)
	rateLimitResetInt, _ := strconv.ParseInt(rateLimitResetTime, 0, 64)
	rateLimitResetTimestamp := time.Unix(rateLimitResetInt, 0)

	key.dataLock.Lock()
	key.remainingUses = remainingRequestsInt
	key.rateLimitResetTime = rateLimitResetTimestamp
	key.dataLock.Unlock()

	if remainingRequestsInt < 1000 {
		// TODO(Shikkic): add a monitoring alert for stuff like this.
		log.Printf("Github rate limit reset time %s\n", rateLimitResetTime)
		log.Printf("Github rate limit remaining requests %s\n", remaingRequests)
	}

	return key
}

func (key *apiKey) canBeUsed() bool {
	key.dataLock.RLock()
	hasRemainingRequests := key.remainingUses < 1
	key.dataLock.RUnlock()

	return hasRemainingRequests
}

func (key *apiKey) waitUntilUseful() {
	key.dataLock.RLock()
	resetTime := key.rateLimitResetTime
	sleepTime := resetTime.Sub(time.Now())
	key.dataLock.RUnlock()

	log.Printf("Github API Key is sleeping until %s.\n", sleepTime.String())
	time.Sleep(sleepTime)
}

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
