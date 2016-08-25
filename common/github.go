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

// GitHubGophrPackageOrgName is the  Github organization name for all versioned packages
var (
	GitHubGophrPackageOrgName = "gophr-packages"
	GitHubBaseAPIURL          = "https://api.github.com"
)

var (
	commitsUntilParameter = "until"
	commitsAfterParameter = "after'"
)

// GitHubRequestService is the library responsible for managing all outbound
// requests to GitHub
type GitHubRequestService struct {
	APIKeyChain *GitHubAPIKeyChain
}

// NewGitHubRequestService initialies a new GitHubRequestService and APIKeyChain
func NewGitHubRequestService() *GitHubRequestService {
	newGitHubRequestService := GitHubRequestService{}
	newGitHubRequestService.APIKeyChain = NewGitHubAPIKeyChain()

	return &newGitHubRequestService
}

// FetchGitHubDataForPackageModel fetchs current repo data of a given packageModel
// TODO optimize this with FFJSON models
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

// CheckGitHubRepoExists returns whether a repo exists
// TODO(Shikkic): Instead of pinging try downloading refs, might be more sustainable?
func (gitHubRequestService *GitHubRequestService) CheckGitHubRepoExists(
	packageModel models.PackageModel,
) error {
	repoName := BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo)
	url := fmt.Sprintf("https://github.com/%s/%s", GitHubGophrPackageOrgName, repoName)
	resp, err := http.Get(url)

	if err != nil {
		log.Println("Error occured during request")
		return err
	}

	if resp.StatusCode == 404 {
		log.Printf("No Github repo exists in %s org with the name %s \n", GitHubGophrPackageOrgName, repoName)
		return nil
	}

	return fmt.Errorf("Error status code %d, a repo with that name already exists.", resp.StatusCode)
}

// CreateNewGitHubRepo if repo doesn't already exist will create a new
// repo on the GitHubGophrPackageOrgName repo
func (gitHubRequestService *GitHubRequestService) CreateNewGitHubRepo(
	packageModel models.PackageModel,
) error {
	err := gitHubRequestService.CheckGitHubRepoExists(packageModel)
	if err != nil {
		log.Println(err)
		return err
	}

	APIKeyModel := gitHubRequestService.APIKeyChain.getAPIKeyModel()
	log.Println(APIKeyModel)
	fmt.Printf("%+v \n", APIKeyModel)
	log.Printf("Determining APIKey %s \n", APIKeyModel.Key)

	JSONBody := buildNewGitHubRepoJSONBody(packageModel)
	gitHubURL := buildNewGitHubRepoAPIURL(packageModel, APIKeyModel)

	req, err := http.Post(gitHubURL, "application/json", JSONBody)
	defer req.Body.Close()

	if err != nil {
		log.Printf("Error occured whilecreating new github repo %s \n", err)
		return err
	}
	if req.StatusCode != 201 {
		log.Printf("Error creating repo was not successful \n")
		return errors.New("Error creating repo was not successful")
	}

	return nil
}

func buildNewGitHubRepoAPIURL(
	packageModel models.PackageModel,
	APIKeyModel *GitHubAPIKeyModel,
) string {
	url := fmt.Sprintf("%s/orgs/%s/repos?access_token=%s",
		GitHubBaseAPIURL,
		GitHubGophrPackageOrgName,
		APIKeyModel.Key,
	)
	return url
}

func buildNewGitHubRepoJSONBody(
	packageModel models.PackageModel,
) *bytes.Buffer {
	author := *packageModel.Author
	repo := *packageModel.Repo
	newGitHubRepoName := BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo)
	description := fmt.Sprintf("Auto generated and versioned go package for %s/%s", author, repo)
	homepage := fmt.Sprintf("https://github.com/%s/%s", author, repo)

	JSONStruct := NewGitHubRepoDTO{Name: newGitHubRepoName, Description: description, Homepage: homepage}
	JSONByteBuffer := new(bytes.Buffer)
	json.NewEncoder(JSONByteBuffer).Encode(JSONStruct)
	return JSONByteBuffer
}

// BuildNewGitHubRepoName creates a new repo name hash uses for repo creation
// and lookup. Eliminates collision between similiar usernames and packages
func BuildNewGitHubRepoName(author string, repo string) string {
	return fmt.Sprintf("%d%s%d%s", len(author), author, len(repo), repo)
}

// FetchCommitSHA Fetches a commitSHA closest to a given timestamp
func (gitHubRequestService *GitHubRequestService) FetchCommitSHA(
	packageModel models.PackageModel,
	timestamp time.Time,
) (string, error) {
	commitSHA, err := gitHubRequestService.fetchCommitSHAByTimeSelector(packageModel, timestamp, commitsUntilParameter)
	if err == nil {
		return commitSHA, nil
	}

	log.Printf("%s \n", err)
	commitSHA, err = gitHubRequestService.fetchCommitSHAByTimeSelector(packageModel, timestamp, commitsAfterParameter)
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
		GitHubBaseAPIURL,
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

// NewGitHubRepoDTO used as a DTO for building POST requests to Github
// to create new repos
type NewGitHubRepoDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}

// GitHubPackageModelDTO TODO Optimize this
type GitHubPackageModelDTO struct {
	Package      models.PackageModel
	ResponseBody map[string]interface{}
}

// ParseStarCount TODO Won't need this after implementing FFJSON
func ParseStarCount(responseBody map[string]interface{}) int {
	starCount := responseBody["stargazers_count"]
	if starCount == nil {
		return 0
	}

	return int(starCount.(float64))
}

// BuildGitHubBranch creates a new ref based on a hash of the old ref
func BuildGitHubBranch(ref string) string {
	repoHash := ref[:len(ref)-1]
	return repoHash
}

// BuildRemoteURL creates a remote url for a packageModel based on it's ref
func BuildRemoteURL(packageModel *models.PackageModel, ref string) string {
	repoName := BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo)
	remoteURL := fmt.Sprintf(gitHubRemoteOrigin, repoName)
	return remoteURL
}
