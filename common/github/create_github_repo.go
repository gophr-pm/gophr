package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	git "github.com/libgit2/git2go"
	"github.com/skeswa/gophr/common/dtos"
	"github.com/skeswa/gophr/common/models"
)

// CreateNewRepo if repo doesn't already exist will create a new
// repo on the GitHubGophrPackageOrgName repo
func (gitHubRequestService *RequestService) CreateNewRepo(packageModel models.PackageModel) error {
	log.Println("Creating New Repo")
	folderName := fmt.Sprintf(
		"%s.git",
		BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo),
	)
	log.Printf("Folder name %s \n", folderName)
	if err := checkIfFolderExists(folderName); err == nil {
		log.Println("Folder exists")
		return nil
	}

	log.Println("Mkdir")
	err := os.Mkdir(
		fmt.Sprintf("/repos/%s", folderName),
		0644,
	)

	if err != nil {
		if checkIfFolderExists(folderName); err != nil {
			return fmt.Errorf("Error, could not create folder or verify that it currently exists. %v", err)
		}
	}

	// Git init bare
	_, err = git.InitRepository(fmt.Sprintf("/repos/%s", folderName), true)
	if err != nil {
		return fmt.Errorf("Error, could not initialize new repository. %v", err)
	}

	return err
}

func checkIfFolderExists(folderName string) error {
	_, err := os.Stat(fmt.Sprintf("/repos/%s", folderName))
	return err
}

// CreateNewGitHubRepo if repo doesn't already exist will create a new
// repo on the GitHubGophrPackageOrgName repo
func (gitHubRequestService *RequestService) CreateNewGitHubRepo(author, repo string) error {
	err := gitHubRequestService.CheckGitHubRepoExists(author, repo)
	if err != nil {
		log.Println(err)
		return err
	}

	APIKeyModel := gitHubRequestService.APIKeyChain.getAPIKeyModel()
	log.Println(APIKeyModel)
	log.Printf("%+v \n", APIKeyModel)
	log.Printf("Determining APIKey %s \n", APIKeyModel.Key)

	JSONBody := buildNewGitHubRepoJSONBody(author, repo)
	gitHubURL := buildNewGitHubRepoAPIURL(author, repo, APIKeyModel)

	req, err := http.Post(gitHubURL, "application/json", JSONBody)
	defer req.Body.Close()

	if err != nil {
		log.Printf("Error occurred whilecreating new github repo %s \n", err)
		return err
	}
	if req.StatusCode != 201 {
		log.Printf("Error creating repo was not successful \n")
		return errors.New("Error creating repo was not successful")
	}

	return nil
}

func buildNewGitHubRepoJSONBody(author, repo string) *bytes.Buffer {
	newGitHubRepoName := BuildNewGitHubRepoName(author, repo)
	description := fmt.Sprintf("Auto generated and versioned go package for %s/%s", author, repo)
	homepage := fmt.Sprintf("https://github.com/%s/%s", author, repo)

	JSONStruct := dtos.NewGitHubRepoDTO{Name: newGitHubRepoName, Description: description, Homepage: homepage}
	JSONByteBuffer := new(bytes.Buffer)
	json.NewEncoder(JSONByteBuffer).Encode(JSONStruct)
	return JSONByteBuffer
}

func buildNewGitHubRepoAPIURL(
	author string,
	repo string,
	APIKeyModel *APIKeyModel) string {
	url := fmt.Sprintf("%s/orgs/%s/repos?access_token=%s",
		GitHubBaseAPIURL,
		GitHubGophrPackageOrgName,
		APIKeyModel.Key,
	)
	return url
}
