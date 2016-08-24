package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/skeswa/gophr/common/dtos"
	"github.com/skeswa/gophr/common/models"
)

var (
	github_base_api_url     = "https://api.github.com"
	github_gophr_org_name   = "gophr-pm"
	commits_until_parameter = "until"
	commits_after_parameter = "after'"
)

type GitHubRequestService struct {
	APIKeyChain *GitHubAPIKeyChain
}

func NewGitHubRequestService() *GitHubRequestService {
	newGitHubRequestService := GitHubRequestService{}
	newGitHubRequestService.APIKeyChain = NewGitHubAPIKeyChain()

	return &newGitHubRequestService
}

// TODO optimize this
func (gitHubRequestService *GitHubRequestService) FetchGitHubDataForPackageModel(
	packageModel models.PackageModel,
) (map[string]interface{}, error) {
	APIKeyModel := gitHubRequestService.APIKeyChain.getAPIKeyModel()
	log.Println(APIKeyModel)
	fmt.Printf("%+v \n", APIKeyModel)
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
	url := fmt.Sprintf("%s/repos/%s/%s?access_token=%s", github_base_api_url, author, repo, APIKeyModel.Key)
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

func (gitHubRequestService *GitHubRequestService) CheckGitHubRepoExists(
	packageModel models.PackageModel,
) error {
	repoName := buildNewGitHubRepoName(packageModel)
	url := fmt.Sprintf("https://github.com/%s/%s", github_gophr_org_name, repoName)
	resp, err := http.Get(url)

	if err != nil {
		log.Println("Error, a repo with that name already exists")
		return err
	}

	if resp.StatusCode != 404 {
		log.Println("Error, a repo with that name already exists")
		return errors.New("Error, a repo with that name already exists.")
	}

	return nil
}

func (gitHubRequestService *GitHubRequestService) CreateNewGitHubRepo(
	packageModel models.PackageModel,
) error {
	APIKeyModel := gitHubRequestService.APIKeyChain.getAPIKeyModel()
	log.Println(APIKeyModel)
	fmt.Printf("%+v \n", APIKeyModel)
	log.Printf("Determining APIKey %s \n", APIKeyModel.Key)

	JSONBody := buildNewGitHubRepoJSONBody(packageModel)
	gitHubURL := buildNewGitHubRepoAPIURL(packageModel, APIKeyModel)

	req, err := http.Post(gitHubURL, "application/json", JSONBody)
	defer req.Body.Close()

	if err != nil {
		log.Printf("Error occured whilecreating new github repo %s", err)
		return err
	}
	if req.StatusCode != 201 {
		log.Printf("Error creating repo was not successful")
		return errors.New("Error creating repo was not successful")
	}

	return nil
}

func buildNewGitHubRepoAPIURL(
	packageModel models.PackageModel,
	APIKeyModel *GitHubAPIKeyModel,
) string {
	url := fmt.Sprintf("%s/orgs/%s/repos?access_token=%s",
		github_base_api_url,
		github_gophr_org_name,
		APIKeyModel.Key,
	)
	return url
}

func buildNewGitHubRepoJSONBody(
	packageModel models.PackageModel,
) *bytes.Buffer {
	author := *packageModel.Author
	repo := *packageModel.Repo
	newGitHubRepoName := buildNewGitHubRepoName(packageModel)
	description := fmt.Sprintf("Auto generated and versioned go package for %s/%s", author, repo)
	homepage := fmt.Sprintf("https://github.com/%s/%s", author, repo)

	JSONStruct := NewGitHubRepo{Name: newGitHubRepoName, Description: description, Homepage: homepage}
	JSONByteBuffer := new(bytes.Buffer)
	json.NewEncoder(JSONByteBuffer).Encode(JSONStruct)
	return JSONByteBuffer
}

func buildNewGitHubRepoName(packageModel models.PackageModel) string {
	author := *packageModel.Author
	repo := *packageModel.Repo
	newGitHubRepoName := fmt.Sprintf("%s-%s", author, repo)
	return newGitHubRepoName
}

func (gitHubRequestService *GitHubRequestService) FetchCommitSHA(
	packageModel models.PackageModel,
	timestamp time.Time,
) (string, error) {
	commitSHA, err := gitHubRequestService.fetchCommitSHAByTimeSelector(packageModel, timestamp, commits_until_parameter)
	if err == nil {
		return commitSHA, nil
	}

	log.Printf("%s \n", err)
	commitSHA, err = gitHubRequestService.fetchCommitSHAByTimeSelector(packageModel, timestamp, commits_after_parameter)
	if err == nil {
		return commitSHA, nil
	}

	log.Printf("%s \n", err)
	refs, err := FetchRefs(*packageModel.Author, *packageModel.Repo)
	if err != nil {
		return refs.MasterRefHash, nil
	}

	return "", err
}

func (gitHubRequestService *GitHubRequestService) fetchCommitSHAByTimeSelector(
	packageModel models.PackageModel,
	timestamp time.Time,
	timeSelector string,
) (string, error) {
	APIKeyModel := gitHubRequestService.APIKeyChain.getAPIKeyModel()
	log.Println(APIKeyModel)
	fmt.Printf("%+v \n", APIKeyModel)
	log.Printf("Determining APIKey %s \n", APIKeyModel.Key)

	githubURL := buildGitHubRepoCommitsFromTimestampAPIURL(packageModel, *APIKeyModel, timestamp, timeSelector)
	log.Printf("Fetching GitHub data for %s \n", githubURL)

	resp, err := http.Get(githubURL)
	defer resp.Body.Close()

	if err != nil {
		return "", errors.New("Request error.")
	}

	if resp.StatusCode == 404 {
		log.Println("PackageModel was not found on Github")
		return "", nil
	}

	APIKeyModel.incrementUsageFromResponseHeader(resp.Header)
	APIKeyModel.print()

	commitSHA, err := parseGitHubCommitLookUpResponseBody(resp)
	if err != nil {
		return "", err
	}

	return commitSHA, nil
}

func buildGitHubRepoCommitsFromTimestampAPIURL(
	packageModel models.PackageModel,
	APIKeyModel GitHubAPIKeyModel,
	timestamp time.Time,
	timeSelector string,
) string {
	author := *packageModel.Author
	repo := *packageModel.Repo

	url := fmt.Sprintf("%s/repos/%s/%s/commits?%s=%s&access_token=%s",
		github_base_api_url,
		author,
		repo,
		timeSelector,
		strings.Replace(timestamp.String(), " ", "", -1),
		APIKeyModel.Key,
	)
	return url
}

func parseGitHubCommitLookUpResponseBody(response *http.Response) (string, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("Failed to parse response body")
	}

	var commitSHAArray []dtos.GitCommitDTO
	err = json.Unmarshal(body, &commitSHAArray)
	if err != nil {
		return "", errors.New("Failed to unmarshal response body")
	}

	if len(commitSHAArray) >= 1 {
		return commitSHAArray[0].SHA, nil
	}

	return "", errors.New("No commit SHAs available for timestamp given")
}

// ==== Misc ====

type NewGitHubRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}

// TODO Optimize this
type GitHubPackageModelDTO struct {
	Package      models.PackageModel
	ResponseBody map[string]interface{}
}

// TODO Won't need this after implementing FFJSON
func ParseStarCount(responseBody map[string]interface{}) int {
	starCount := responseBody["stargazers_count"]
	if starCount == nil {
		return 0
	}

	return int(starCount.(float64))
}
