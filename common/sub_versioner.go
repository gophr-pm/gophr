package common

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/models"
)

var (
	folder_name                string
	github_remote_origin       = "git@github.com:gophr-pm/%s.git"
	navigate_to_package_folder = "cd /tmp/%s"
)

var (
	initialize_repo    = "cd /tmp && mkdir %s && cd %s && git init"
	create_branch      = "%s && git checkout -b %s"
	set_remote_command = "%s && git remote add origin %s"
	fetch_repo_archive = "%s && wget https://github.com/%s/%s/archive/%s.zip"
	unzip_repo_archive = "%s && unzip %s.zip && rm %s.zip"
	add_files          = "%s && git add . "
	commit_files       = "%s && git commit -m \" %s \""
	push_files         = "%s && git push --set-upstream origin %s"
)

// TODO (@Shikkic): Possibly add Channel as a param
func SubVersionPackageModel(packageModel models.PackageModel, version string) {
	// Set working folder_name for package
	// TODO (Shikkic): Build the folder name based with algo
	folder_name = *packageModel.Author + "-" + *packageModel.Repo

	// Instantiate New Github Request Service
	log.Println("Initializing gitHub component")
	gitHubRequestService := common.NewGitHubRequestService()

	log.Printf("Creating new Github repo for %s/%s at %s",
		*packageModel.Author,
		*packageModel.Repo,
		buildTagName(&packageModel),
	)
	err := gitHubRequestService.CreateNewGitHubRepo(packageModel)
	log.Printf("%s", err)

	log.Printf("Initializing folder and initializing git repo for %s \n", folder_name)
	err = initializeRepoCMD(&packageModel)
	checkError(err, folder_name)

	log.Printf("Creating branch %s \n", buildBranchName(&packageModel))
	err = createBranchCMD(&packageModel)
	checkError(err, folder_name)

	log.Printf("Setting remote branch name %s \n", buildRemoteName(&packageModel))
	err = setRemoteCMD(&packageModel)
	checkError(err, folder_name)

	log.Printf("Fetching github archive for %s/%s with tag %s \n",
		*packageModel.Author,
		*packageModel.Repo,
		buildTagName(&packageModel),
	)
	err = fetchArchiveCMD(&packageModel)
	checkError(err, folder_name)

	log.Printf("Fetching github archive for %s/%s with tag %s \n",
		*packageModel.Author,
		*packageModel.Repo,
		buildTagName(&packageModel),
	)
	err = unzipArchiveCMD(&packageModel)
	checkError(err, folder_name)

	// TODO subverisoning
	// Create array of all sub-dependencies

	log.Println("Adding unarchived repo data to branch")
	err = addFilesCMD()
	checkError(err, folder_name)

	log.Println("Commiting repo data to branch")
	err = commitFilesCMD(&packageModel)
	checkError(err, folder_name)

	log.Printf("Pushing files to branch %s \n", buildBranchName(&packageModel))
	err = pushFilesCMD(&packageModel)
	checkError(err, folder_name)

	cleanUpExitHook(folder_name)
}

func initializeRepoCMD(packageModel *models.PackageModel) error {
	log.Println("Initializing folder and repo commmand")
	cmd := fmt.Sprintf(initialize_repo, folder_name, folder_name)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)
	return err
}

func createBranchCMD(packageModel *models.PackageModel) error {
	log.Println("Initializing folder and repo commmand")
	navigateFolderCMD := fmt.Sprintf(navigate_to_package_folder, folder_name)
	branchName := buildBranchName(packageModel)
	cmd := fmt.Sprintf(create_branch, navigateFolderCMD, branchName)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)
	return err
}

func setRemoteCMD(packageModel *models.PackageModel) error {
	log.Println("Initializing folder and repo commmand")
	navigateFolderCMD := fmt.Sprintf(navigate_to_package_folder, folder_name)
	remoteName := buildRemoteName(packageModel)

	cmd := fmt.Sprintf(set_remote_command, navigateFolderCMD, remoteName)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)
	return err
}
func fetchArchiveCMD(packageModel *models.PackageModel) error {
	log.Println("Fetching and Unzipping Archive for tag")
	tag := buildTagName(packageModel)
	navigateFolderCMD := fmt.Sprintf(navigate_to_package_folder, folder_name)

	cmd := fmt.Sprintf(fetch_repo_archive, navigateFolderCMD, *packageModel.Author, *packageModel.Repo, tag)
	log.Printf("%s FETCH ARCHIVE COMMAND", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)

	return err
}

func unzipArchiveCMD(packageModel *models.PackageModel) error {
	log.Println("Fetching and Unzipping Archive for tag")
	navigateFolderCMD := fmt.Sprintf(navigate_to_package_folder, folder_name)
	tag := buildTagName(packageModel)

	cmd := fmt.Sprintf(unzip_repo_archive, navigateFolderCMD, tag, tag)
	log.Printf("%s UNZIP ARCHIVE COMMAND", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)

	return err
}

func addFilesCMD() error {
	navigateFolderCMD := fmt.Sprintf(navigate_to_package_folder, folder_name)
	cmd := fmt.Sprintf(add_files, navigateFolderCMD)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)

	return err
}

// TODO add version
func commitFilesCMD(packageModel *models.PackageModel) error {
	navigateFolderCMD := fmt.Sprintf(navigate_to_package_folder, folder_name)
	commitMessage := fmt.Sprintf("Created Versin Repo of %s / %s @ %s", *packageModel.Author, *packageModel.Repo, "1.0")
	cmd := fmt.Sprintf(commit_files, navigateFolderCMD, commitMessage)
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)

	return err
}

func pushFilesCMD(packageModel *models.PackageModel) error {
	navigateFolderCMD := fmt.Sprintf(navigate_to_package_folder, folder_name)
	branchName := buildBranchName(packageModel)
	cmd := fmt.Sprintf(push_files, navigateFolderCMD, branchName)
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s 		\n", out)

	return err
}

// Helper functions

// TODO Fix this
func buildBranchName(packageModel *models.PackageModel) string {
	log.Println("Generating Branch Name")
	return "master"
}

// TODO Fix this
func buildRemoteName(packageModel *models.PackageModel) string {
	remoteURL := fmt.Sprintf(github_remote_origin, *packageModel.Author+"-"+*packageModel.Repo)
	log.Printf("Generating Remote URL %s \n", remoteURL)
	return remoteURL
}

func buildTagName(packageModel *models.PackageModel) string {
	return "master"
}

// Exit Hook to clean up files
func cleanUpExitHook(folderName string) {
	log.Printf("Exit hook initiated deleting folder %s \n", folderName)
	if strings.ContainsAny(folderName, "/") == true {
		cmd := fmt.Sprintf("cd /tmp && rm -rf %s", folderName)
		out, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			log.Println("Could not properly engage exit hook")
			log.Fatalln(err)
		}
		log.Printf("Output: %s", out)
	} else {
		log.Println("Cowardly refusing to initiate exit hook. Will not rm -rf folder names that contains any leading '/'")
	}
}

// Check Error and Engage Exit Hook if fatal error occured
func checkError(err error, folderName string) {
	if err != nil {
		cleanUpExitHook(folderName)
		log.Fatalf("Error occured %s \n", err)
	}
}

/*

Package - skeswa / gophr @ 1.0

push to len(author)+skeswa+len(repo)+gophr

*/
