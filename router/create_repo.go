package main

import (
	"fmt"
	"os"

	"github.com/skeswa/gophr/common/github"
	"github.com/skeswa/gophr/common/models"
)

// CreateNewRepo if repo doesn't already exist will create a new
// repo on the GitHubGophrPackageOrgName repo
func CreateNewRepo(packageModel models.PackageModel) error {
	folderName := fmt.Sprintf(
		"%s.git",
		github.BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo),
	)
	err := os.Mkdir(
		fmt.Sprintf("/data/repos/%s", folderName),
		0644,
	)
	return err
}
