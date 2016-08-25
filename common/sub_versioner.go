package common

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/skeswa/gophr/common/models"
)

var (
	folderName              string
	gitHubRemoteOrigin      = "git@github.com:gophr-pm/%s.git"
	navigateToPackageFolder = "cd /tmp/%s"
)

var (
	initalizeRepo    = "cd /tmp && mkdir %s && cd %s && git init"
	createBranch     = "%s && git checkout -b %s"
	setRemoteCommand = "%s && git remote add origin %s"
	fetchRepoArchive = "%s && wget https://github.com/%s/%s/archive/%s.zip"
	unzipRepoArchive = "%s && unzip %s.zip && rm %s.zip"
	addFiles         = "%s && git add . "
	commitFiles      = "%s && git commit -m \" %s \""
	pushFiles        = "%s && git push --set-upstream origin %s"
)

//  SubVersionPackageModel TODO (@Shikkic): Possibly add Channel as a param
func SubVersionPackageModel(packageModel *models.PackageModel, ref string) {
	// Set working folderName for package
	folderName = BuildNewGitHubRepoName(packageModel)

	// Instantiate New Github Request Service
	log.Println("Initializing gitHub component")
	gitHubRequestService := NewGitHubRequestService()

	log.Printf("Creating new Github repo for %s/%s at %s",
		*packageModel.Author,
		*packageModel.Repo,
		ref,
	)
	err := gitHubRequestService.CreateNewGitHubRepo(*packageModel)
	log.Printf("%s", err)

	log.Printf("Initializing folder and initializing git repo for %s \n", folderName)
	err = initializeRepoCMD(packageModel)
	checkError(err, folderName)

	log.Printf("Creating branch %s \n", buildGitHubBranch(packageModel))
	err = createBranchCMD(packageModel, ref)
	checkError(err, folderName)

	log.Printf("Setting remote branch url %s \n", buildGitHubBranch(packageModel))
	err = setRemoteCMD(packageModel, ref)
	checkError(err, folderName)

	log.Printf("Fetching github archive for %s/%s with tag %s \n",
		*packageModel.Author,
		*packageModel.Repo,
		ref,
	)
	err = fetchArchiveCMD(packageModel, ref)
	checkError(err, folderName)

	log.Printf("Fetching github archive for %s/%s with tag %s \n",
		*packageModel.Author,
		*packageModel.Repo,
		ref,
	)
	err = unzipArchiveCMD(packageModel, ref)
	checkError(err, folderName)

	// TODO subverisoning
	// Create array of all sub-dependencies

	log.Println("Adding unarchived repo data to branch")
	err = addFilesCMD()
	checkError(err, folderName)

	log.Println("Commiting repo data to branch")
	err = commitFilesCMD(packageModel)
	checkError(err, folderName)

	log.Printf("Pushing files to branch %s \n", buildRemoteURL(packageModel, ref))
	err = pushFilesCMD(packageModel, ref)
	checkError(err, folderName)

	cleanUpExitHook(folderName)
}

func initializeRepoCMD(packageModel *models.PackageModel) error {
	log.Println("Initializing folder and repo commmand")
	cmd := fmt.Sprintf(initalizeRepo, folderName, folderName)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)
	return err
}

func createBranchCMD(packageModel *models.PackageModel, ref string) error {
	log.Println("Initializing folder and repo commmand")
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	cmd := fmt.Sprintf(createBranch, navigateFolderCMD, buildGitHubBranch(packageModel))
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)
	return err
}

func setRemoteCMD(packageModel *models.PackageModel, ref string) error {
	log.Println("Initializing folder and repo commmand")
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	remoteURL := buildRemoteURL(packageModel, ref)
	cmd := fmt.Sprintf(setRemoteCommand, navigateFolderCMD, remoteURL)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)
	return err
}
func fetchArchiveCMD(packageModel *models.PackageModel, ref string) error {
	log.Println("Fetching and Unzipping Archive for tag")
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)

	cmd := fmt.Sprintf(fetchRepoArchive, navigateFolderCMD, *packageModel.Author, *packageModel.Repo, ref)
	log.Printf("%s FETCH ARCHIVE COMMAND", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)

	return err
}

func unzipArchiveCMD(packageModel *models.PackageModel, ref string) error {
	log.Println("Fetching and Unzipping Archive for tag")
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	cmd := fmt.Sprintf(unzipRepoArchive, navigateFolderCMD, ref, ref)
	log.Printf("%s UNZIP ARCHIVE COMMAND", cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)

	return err
}

func addFilesCMD() error {
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	cmd := fmt.Sprintf(addFiles, navigateFolderCMD)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)

	return err
}

// TODO add version
func commitFilesCMD(packageModel *models.PackageModel) error {
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	commitMessage := fmt.Sprintf("Created Versin Repo of %s / %s @ %s", *packageModel.Author, *packageModel.Repo, "1.0")
	cmd := fmt.Sprintf(commitFiles, navigateFolderCMD, commitMessage)
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)

	return err
}

func pushFilesCMD(packageModel *models.PackageModel, ref string) error {
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	remoteURL := buildRemoteURL(packageModel, ref)
	cmd := fmt.Sprintf(pushFiles, navigateFolderCMD, remoteURL)
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s 		\n", out)
	return err
}

// Helper functions

func buildRemoteURL(packageModel *models.PackageModel, ref string) string {
	repoHash := buildGitHubBranch(packageModel)
	remoteURL := fmt.Sprintf(gitHubRemoteOrigin, repoHash)
	return remoteURL
}

func buildGitHubBranch(packageModel *models.PackageModel) string {
	repoName := BuildNewGitHubRepoName(packageModel)
	repoHash := repoName[:len(repoName)-1]
	return repoHash
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
