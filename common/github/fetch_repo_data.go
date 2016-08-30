package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/skeswa/gophr/common/models"
)

// FetchGitHubDataForPackageModel fetchs current repo data of a given packageModel
// TODO optimize this with FFJSON models
func (gitHubRequestService *GitHubRequestService) FetchGitHubDataForPackageModel(
	packageModel models.PackageModel,
) (map[string]interface{}, error) {
	APIKeyModel := gitHubRequestService.APIKeyChain.getAPIKeyModel()
	log.Println(APIKeyModel)
	log.Printf("%+v \n", APIKeyModel)
	log.Printf("Determining APIKey %s \n", APIKeyModel.Key)
	githubURL := buildGitHubRepoDataAPIURL(packageModel, *APIKeyModel)
	log.Printf("Fetching GitHub data for %s \n", githubURL)

	resp, err := http.Get(githubURL)
	defer resp.Body.Close()

	if err != nil {
		return nil, errors.New("Request error.")
	}

	if resp.StatusCode == 404 {
		log.Println("PackageModel was not found on Github")
		return nil, nil
	}

	APIKeyModel.incrementUsageFromResponseHeader(resp.Header)
	APIKeyModel.print()

	responseBodyMap, err := parseGitHubRepoDataResponseBody(resp)
	if err != nil {
		return nil, err
	}

	return responseBodyMap, nil
}

func buildGitHubRepoDataAPIURL(packageModel models.PackageModel, APIKeyModel GitHubAPIKeyModel) string {
	author := *packageModel.Author
	repo := *packageModel.Repo
	url := fmt.Sprintf("%s/repos/%s/%s?access_token=%s", GitHubBaseAPIURL, author, repo, APIKeyModel.Key)
	return url
}

// TODO Optimize this with ffjson struct!
func parseGitHubRepoDataResponseBody(response *http.Response) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("Failed to parse response body")
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return nil, errors.New("Failed to unmarshal response body")
	}

	return bodyMap, nil
}
