package common

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var gitHubAPIBaseUrl = "https://api.github.com/"

// GitHubAPIKeyModel is a struct representing on individual github key in the database
type GitHubAPIKeyModel struct {
	Key                string
	RemainingUses      int
	RequestsPerHour    int
	RateLimitResetTime time.Time
}

// Should pass the number of remaining here instead of -1
// Pass remaining
func (apiKey *GitHubAPIKeyModel) incrementUsage(remainingRequests string, rateLimitResetTime string) {
	n, _ := strconv.ParseInt(rateLimitResetTime, 0, 64)
	resetTime := time.Unix(n, 0)
	l, _ := strconv.Atoi(remainingRequests)
	log.Println("Remaining uses in save ", l)
	apiKey.RemainingUses = l
	apiKey.RateLimitResetTime = resetTime
}

// TODO sent request to api route thats passed in
// TODO should return error
func (apiKey *GitHubAPIKeyModel) prime() {
	githubURL := gitHubAPIBaseUrl + "repos/a/b" + "?access_token=" + apiKey.Key
	log.Println(githubURL)
	resp, err := http.Get(githubURL)
	if err != nil {
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
