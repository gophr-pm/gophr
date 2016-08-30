package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/skeswa/gophr/common/dtos"
	"github.com/skeswa/gophr/common/models"
)

// CreateNewGitHubRepo if repo doesn't already exist will create a new
// repo on the GitHubGophrPackageOrgName repo
func (gitHubRequestService *RequestService) CreateNewGitHubRepo(
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

func buildNewGitHubRepoJSONBody(
	packageModel models.PackageModel,
) *bytes.Buffer {
	author := *packageModel.Author
	repo := *packageModel.Repo
	newGitHubRepoName := BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo)
	description := fmt.Sprintf("Auto generated and versioned go package for %s/%s", author, repo)
	homepage := fmt.Sprintf("https://github.com/%s/%s", author, repo)

	JSONStruct := dtos.NewGitHubRepoDTO{Name: newGitHubRepoName, Description: description, Homepage: homepage}
	JSONByteBuffer := new(bytes.Buffer)
	json.NewEncoder(JSONByteBuffer).Encode(JSONStruct)
	return JSONByteBuffer
}

func buildNewGitHubRepoAPIURL(
	packageModel models.PackageModel,
	APIKeyModel *APIKeyModel,
) string {
	url := fmt.Sprintf("%s/orgs/%s/repos?access_token=%s",
		GitHubBaseAPIURL,
		GitHubGophrPackageOrgName,
		APIKeyModel.Key,
	)
	return url
}
