package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common"
)

var (
	api_key_pool      = 1
	requests_per_hour = 5000
	GITHUB_API_KEYS   []*GitHubAPIKeyModel
)

func ReIndexPackageGitHubStats(session *gocql.Session) {
	packageModels, err := common.ScanAllPackageModels(session)
	log.Println(err)
	log.Printf("%d packages found", len(packageModels))

	totalNumPackages := len(packageModels)
	log.Printf("Total num packages found = %d ", totalNumPackages)

	var wg sync.WaitGroup
	nbConcurrentInserts := 20
	packageChan := make(chan *common.PackageModel, 20)

	// TODO PREP API KEYS
	// FOR EACH API KEY IN DB GENERATE NEW APIKeyModel (create a struct for this)
	// APIKeyModel = [
	//   remaining_uses: 5,000 (default number)
	//   requests_per_hour: 5,000
	//   rate_limit_reset: timestamp or 00:00:00:00
	//

	log.Println("Initializing GitHub Key Model")
	// TODO create initialize function
	apiKeyModel := GitHubAPIKeyModel{
		"24a474ffd3e884a7cb8a14d0408d654b361806b9",
		5000,
		5000,
		time.Time{},
	}

	log.Println("GitHub APIKey collection Initialized")
	GITHUB_API_KEYS = append(GITHUB_API_KEYS, &apiKeyModel)

	// TODO Will consume the new packages with updated stars and store them in the database
	log.Printf("Spinning up %d consumers", nbConcurrentInserts)
	for i := 0; i < nbConcurrentInserts; i++ {
		wg.Add(1)
		go func() {
			for packageModel := range packageChan {
				log.Println(packageModel)
			}
			wg.Done()
		}()
	}

	log.Printf("Preparing to fetch stars for %d repos", totalNumPackages)
	var count int = 0
	for _, packageModel := range packageModels {
		log.Printf("PROCESSING PACKAGE %s/%s #%d \n", *packageModel.Author, *packageModel.Repo, count)
		packageStars := fetchStarsForPackageModel(packageModel)
		log.Printf("This package has %d stars \n", packageStars)
		packageChan <- packageModel
		count++
	}

	close(packageChan)
	wg.Wait()
	log.Println("Finished testing star fetching")
}

func fetchStarsForPackageModel(packageModel *common.PackageModel) int {
	APIKey := getAPIKey()
	for APIKeyIsValid(APIKey) != true {
		timeNow := time.Now()
		log.Println("Time now = ", timeNow)
		resetTime := APIKey.RateLimitResetTime
		log.Println("Reset time = ", resetTime)

		sleepTime := resetTime.Sub(timeNow)
		log.Println("Sleep time is = ", sleepTime)
		time.Sleep(sleepTime)

		// Check if we're past the rate refresh time
		if timeNow.After(resetTime) {
			log.Println("IT'S OGRE WE DID IT!")
			APIKey.reset()
		}
	}

	// Prepare URL
	author := packageModel.Author
	repo := packageModel.Repo
	url := "https://api.github.com/repos/" + *author + "/" + *repo + "?access_token=" + APIKey.Key
	log.Println("Fetching GitHub data for ", url)

	// Make request
	resp, _ := http.Get(url)
	statusCode := resp.StatusCode
	log.Println("STATUS CODE = ", statusCode)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	// Unmarshal the request body
	var bodyMap map[string]interface{}
	err := json.Unmarshal(body, &bodyMap)
	log.Println("json error = ", err)

	// Parse Body Values
	starCount := int(bodyMap["stargazers_count"].(float64))
	log.Println("Star count = ", starCount)

	// Parse Header Values
	rateLimitResetTime := resp.Header.Get("X-RateLimit-Reset")
	log.Println("Rate limit reset time = ", rateLimitResetTime)
	remaingUses := resp.Header.Get("X-RateLimit-Remaining")
	log.Println("Rate limit remaining uses = ", remaingUses)

	// Increment APIKey Usage here
	APIKey.incrementUsage(remaingUses, rateLimitResetTime)
	APIKey.print()

	// TODO add artificial delay
	time.Sleep(600 * time.Millisecond)

	return starCount
}

// TODO pops the first key in the Priority Queue
// TODO create priority queue
func getAPIKey() *GitHubAPIKeyModel {
	return GITHUB_API_KEYS[0]
}

// TODO returns if we should wait or not
func APIKeyIsValid(APIKey *GitHubAPIKeyModel) bool {
	log.Printf("Remaining API calls %d", APIKey.RemainingUses)
	if APIKey.RemainingUses > 0 {
		return true
	}

	return false
}
