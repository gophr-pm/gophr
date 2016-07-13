package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

// TODO Include default values

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

func (apiKey *GitHubAPIKeyModel) reset() {
	apiKey.RemainingUses = 5000
}

func (apiKey *GitHubAPIKeyModel) print() {
	fmt.Printf("%+v", apiKey)
}
