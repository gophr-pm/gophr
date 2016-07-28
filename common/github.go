package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/skeswa/gophr/common/models"
)

type GitHubRequestService struct {
	APIKeyChain *GitHubAPIKeyChain
}

func NewGitHubRequestService() *GitHubRequestService {
	newGitHubRequestService := GitHubRequestService{}
	newGitHubRequestService.APIKeyChain = NewGitHubAPIKeyChain()

	return &newGitHubRequestService
}

func (gitHubRequestService *GitHubRequestService) FetchGitHubDataForPackageModel(packageModel models.PackageModel) (map[string]interface{}, error) {
	APIKeyModel := gitHubRequestService.APIKeyChain.getAPIKeyModel()
	log.Println(APIKeyModel)
	fmt.Printf("%+v \n", APIKeyModel)
	log.Printf("Determining APIKey %s \n", APIKeyModel.Key)
	githubURL := buildGithubAPIURL(packageModel, *APIKeyModel)
	log.Printf("Fetching GitHub data for %s \n", githubURL)

	resp, err := http.Get(githubURL)
	defer resp.Body.Close()

	if err != nil {
		return nil, errors.New("Request error.")
	}

	// TODO Special action for 404s
	if resp.StatusCode == 404 {
		// TODO use package not found error here
		return nil, errors.New("Package does not exist")
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
		return nil, err
	}

	return responseBodyMap, nil
}

// TODO potentially return error here
func buildGithubAPIURL(packageModel models.PackageModel, APIKeyModel GitHubAPIKeyModel) string {
	author := *packageModel.Author
	repo := *packageModel.Repo
	url := "https://api.github.com/repos/" + author + "/" + repo + "?access_token=" + APIKeyModel.Key
	return url
}

func parseResponseBody(response *http.Response) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("Failed to parse response body.")
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, errors.New("Failed to unmarshal response body.")
	}

	return bodyMap, nil
}

func ParseStarCount(responseBody map[string]interface{}) int {
	starCount := responseBody["stargazers_count"]
	if starCount == nil {
		return 0
	}

	return int(starCount.(float64))
}

type GitHubPackageModelDTO struct {
	Package      models.PackageModel
	ResponseBody map[string]interface{}
}
