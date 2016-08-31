package github

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/skeswa/gophr/common/dtos"
	"github.com/skeswa/gophr/common/models"
)

// FetchCommitTimestamp fetches the timestamp of a commit from Github API
func (gitHubRequestService *RequestService) FetchCommitTimestamp(
	packageModel *models.PackageModel,
	commitSHA string,
) (time.Time, error) {
	APIKeyModel := gitHubRequestService.APIKeyChain.getAPIKeyModel()
	log.Println(APIKeyModel)
	log.Printf("%+v \n", APIKeyModel)
	log.Printf("Determining APIKey %s \n", APIKeyModel.Key)

	githubURL := buildGitHubCommitTimestampAPIURL(packageModel, *APIKeyModel, commitSHA)
	log.Printf("Fetching commit timestamp for %s \n", githubURL)

	resp, err := http.Get(githubURL)
	defer resp.Body.Close()

	if err != nil {
		return time.Time{}, errors.New("Request error.")
	}

	if resp.StatusCode == 404 {
		log.Println("PackageModel was not found on Github")
		return time.Time{}, nil
	}

	APIKeyModel.incrementUsageFromResponseHeader(resp.Header)
	APIKeyModel.print()

	timestamp, err := parseGitHubCommitLookUpResponseBody(resp)
	if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}

func buildGitHubCommitTimestampAPIURL(
	packageModel *models.PackageModel,
	APIKeyModel APIKeyModel,
	commitSHA string,
) string {
	author := *packageModel.Author
	repo := *packageModel.Repo

	url := fmt.Sprintf("%s/repos/%s/%s/commits/%s?&access_token=%s",
		GitHubBaseAPIURL,
		author,
		repo,
		commitSHA,
		APIKeyModel.Key,
	)
	return url
}

func parseGitHubCommitLookUpResponseBody(response *http.Response) (time.Time, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return time.Time{}, errors.New("Failed to parse response body")
	}

	var commitLookUpDTO dtos.GitCommitLookUpDTO
	err = ffjson.Unmarshal(body, &commitLookUpDTO)
	if err != nil {
		return time.Time{}, errors.New("Failed to unmarshal response body")
	}

	if commitLookUpDTO.Commit != nil && commitLookUpDTO.Commit.Committer != nil {
		return commitLookUpDTO.Commit.Committer.Date, nil
	}

	return time.Time{}, errors.New("No commit timestamp avilable for given SHA")
}