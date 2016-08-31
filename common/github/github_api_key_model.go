package github

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// APIKeyModel represents a single Github API keys,
// it's responsible for keeping track of API usage via that key
type APIKeyModel struct {
	Key                string
	RemainingUses      int
	RequestsPerHour    int
	RateLimitResetTime time.Time
}

// TODO:(Shikkic) consider revising how we parse the values from header
func (apiKey *APIKeyModel) incrementUsageFromResponseHeader(header http.Header) {
	remaingRequests := header.Get("X-RateLimit-Remaining")
	rateLimitResetTime := header.Get("X-RateLimit-Reset")

	remainingRequestsInt, _ := strconv.Atoi(remaingRequests)
	rateLimitResetInt, _ := strconv.ParseInt(rateLimitResetTime, 0, 64)
	rateLimitResetTimestamp := time.Unix(rateLimitResetInt, 0)

	apiKey.RemainingUses = remainingRequestsInt
	apiKey.RateLimitResetTime = rateLimitResetTimestamp

	log.Printf("Rate limit remaining requests %s \n", remaingRequests)
	log.Printf("Rate limit reset time %s \n", rateLimitResetTime)
	log.Printf("Decrementing APIKeyModel usage to %d uses \n", remainingRequestsInt)
}

// TODO:(Shikkic) consider passing url endpoint to prime, or maybe an enum for more accuracy when pinging GH
func (apiKey *APIKeyModel) prime() {
	gitHubTestURL := fmt.Sprintf("%s/repos/a/b?access_token=%s", GitHubBaseAPIURL, apiKey.Key)
	log.Printf("Preparing to prime APIKeyModel with key %s and url %s \n", apiKey.Key, gitHubTestURL)

	resp, err := http.Get(gitHubTestURL)
	if err != nil {
		log.Println("Could not prime APIKey, fatal error in Github API request")
		log.Fatal(err)
	}
	defer resp.Body.Close()

	apiKey.incrementUsageFromResponseHeader(resp.Header)
}

func (apiKey *APIKeyModel) reset() {
	apiKey.RemainingUses = 5000
}

func (apiKey *APIKeyModel) print() {
	log.Printf("%+v", apiKey)
}
