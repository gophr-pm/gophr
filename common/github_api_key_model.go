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

func (apiKey *GitHubAPIKeyModel) incrementUsage(remainingRequests string, rateLimitResetTime string) {
	rateLimitResetInt, _ := strconv.ParseInt(rateLimitResetTime, 0, 64)
	rateLimitResetTimestamp := time.Unix(rateLimitResetInt, 0)
	remainingRequestsInt, _ := strconv.Atoi(remainingRequests)
	apiKey.RemainingUses = remainingRequestsInt
	apiKey.RateLimitResetTime = rateLimitResetTimestamp
	log.Printf("Decrementing APIKeyModel usage to %d uses \n", remainingRequestsInt)
}

func (apiKey *GitHubAPIKeyModel) prime() {
	githubTestURL := gitHubAPIBaseUrl + "repos/a/b" + "?access_token=" + apiKey.Key
	log.Printf("Preparing to prime APIKeyModel with key %s and url %s", apiKey.Key, githubTestURL)

	resp, err := http.Get(githubTestURL)
	if err != nil {
		log.Println("Could not prime APIKey, fatal error in github api request")
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseHeader := resp.Header
	remaingRequests := responseHeader.Get("X-RateLimit-Remaining")
	rateLimitResetTime := responseHeader.Get("X-RateLimit-Reset")
	apiKey.incrementUsage(remaingRequests, rateLimitResetTime)
}

func (apiKey *GitHubAPIKeyModel) reset() {
	apiKey.RemainingUses = 5000
}

func (apiKey *GitHubAPIKeyModel) print() {
	fmt.Printf("%+v", apiKey)
}
