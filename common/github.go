package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type GitHubRequestService struct {
	APIKeyChain *GitHubAPIKeyChain
}

func NewGitHubRequestService() *GitHubRequestService {
	newGitHubRequestService := GitHubRequestService{}
	newGitHubRequestService.APIKeyChain = NewGitHubAPIKeyChain()

	return &newGitHubRequestService
}

func (gitHubRequestService *GitHubRequestService) FetchGitHubDataForPackageModel(packageModel PackageModel) (map[string]interface{}, error) {
	APIKeyModel := gitHubRequestService.APIKeyChain.getAPIKeyModel()
	log.Println(APIKeyModel)
	fmt.Printf("%+v \n", APIKeyModel)
	log.Printf("Determining APIKey %s \n", APIKeyModel.Key)
	githubURL := buildGithubAPIURL(packageModel, *APIKeyModel)
	log.Printf("Fetching GitHub data for %s \n", githubURL)

	resp, err := http.Get(githubURL)
	defer resp.Body.Close()
	// TODO Drop 404s
	// resp.StatusCode != 200 check for 404s
	if err != nil {
		log.Printf("PANIC %v \n", err)
		log.Printf("STATUS CODE!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! %d", resp.StatusCode)
	}

	responseHeader := resp.Header
	remaingRequests := responseHeader.Get("X-RateLimit-Remaining")
	rateLimitResetTime := responseHeader.Get("X-RateLimit-Reset")
	APIKeyModel.incrementUsage(remaingRequests, rateLimitResetTime)
	log.Printf("Rate limit reset time %s \n", rateLimitResetTime)
	log.Printf("Rate limit remaining requests %s \n", remaingRequests)
	APIKeyModel.print()

	// Parse Body Values
	responseBodyMap, err := parseResponseBody(resp)
	if err != nil {
		log.Printf("PANIC %v \n", err)
		os.Exit(3)
	}

	return responseBodyMap, nil
}

func buildGithubAPIURL(packageModel PackageModel, APIKeyModel GitHubAPIKeyModel) string {
	author := packageModel.Author
	repo := packageModel.Repo
	url := "https://api.github.com/repos/" + *author + "/" + *repo + "?access_token=" + APIKeyModel.Key
	return url
}

type GitHubPackageModelDTO struct {
	Package      PackageModel
	ResponseBody map[string]interface{}
}

func parseResponseBody(response *http.Response) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("PANIC %v \n", err)
		os.Exit(3)
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		log.Printf("PANIC %v \n", err)
		os.Exit(3)
	}

	return bodyMap, nil
}
