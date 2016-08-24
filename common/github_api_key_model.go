package common

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var gitHubAPIBaseUrl = "https://api.github.com/"

type GitHubAPIKeyModel struct {
	Key                string
	RemainingUses      int
	RequestsPerHour    int
	RateLimitResetTime time.Time
}

func (apiKey *GitHubAPIKeyModel) incrementUsageFromResponseHeader(respHeader http.Header) {
	remaingRequests := respHeader.Get("X-RateLimit-Remaining")
	rateLimitResetTime := respHeader.Get("X-RateLimit-Reset")

	remainingRequestsInt, _ := strconv.Atoi(remaingRequests)
	rateLimitResetInt, _ := strconv.ParseInt(rateLimitResetTime, 0, 64)
	rateLimitResetTimestamp := time.Unix(rateLimitResetInt, 0)

	apiKey.RemainingUses = remainingRequestsInt
	apiKey.RateLimitResetTime = rateLimitResetTimestamp

	log.Printf("Rate limit remaining requests %s \n", remaingRequests)
	log.Printf("Rate limit reset time %s \n", rateLimitResetTime)
	log.Printf("Decrementing APIKeyModel usage to %d uses \n", remainingRequestsInt)
}

// TODO consider passing url endpoint to prime, or maybe a enum for more accuracy when pinging GH
func (apiKey *GitHubAPIKeyModel) prime() {
	githubTestURL := gitHubAPIBaseUrl + "repos/a/b" + "?access_token=" + apiKey.Key
	log.Printf("Preparing to prime APIKeyModel with key %s and url %s \n", apiKey.Key, githubTestURL)

	resp, err := http.Get(githubTestURL)
	if err != nil {
		log.Println("Could not prime APIKey, fatal error in github api request")
		log.Fatal(err)
	}
	defer resp.Body.Close()

	apiKey.incrementUsageFromResponseHeader(resp.Header)
}

func (apiKey *GitHubAPIKeyModel) reset() {
	apiKey.RemainingUses = 5000
}

func (apiKey *GitHubAPIKeyModel) print() {
	fmt.Printf("%+v", apiKey)
}
