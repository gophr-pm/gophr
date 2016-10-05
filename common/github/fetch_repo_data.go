package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gophr-pm/gophr/common/models"
)

// FetchGitHubDataForPackageModel fetchs current repo data of a given packageModel
// TODO optimize this with FFJSON models
func (svc *requestService) FetchGitHubDataForPackageModel(
	packageModel models.PackageModel,
) (map[string]interface{}, error) {
	APIKeyModel := svc.APIKeyChain.getAPIKeyModel()
	log.Println(APIKeyModel)
	githubURL := buildGitHubRepoDataAPIURL(*packageModel.Author, *packageModel.Repo, *APIKeyModel)
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

	responseBodyMap, err := parseGitHubRepoDataResponseBody(resp)
	if err != nil {
		return nil, err
	}

	return responseBodyMap, nil
}

func buildGitHubRepoDataAPIURL(
	author string,
	repo string,
	keyModel APIKeyModel,
) string {
	url := fmt.Sprintf(
		"%s/repos/%s/%s?access_token=%s",
		GitHubBaseAPIURL,
		author,
		repo,
		keyModel.Key)
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
