package subv

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/github"
	"github.com/skeswa/gophr/common/models"
	"github.com/skeswa/gophr/common/verdeps"
)

var (
	folderName              string
	gitHubRemoteOrigin      = "git@github.com:gophr-packages/%s.git"
	navigateToPackageFolder = "cd /tmp/%s"
)

var (
	initalizeRepo    = "cd /tmp && mkdir %s && cd %s && git init"
	createBranch     = "%s && git checkout -b %s"
	setRemoteCommand = "%s && git remote add origin %s"
	fetchRepoArchive = "%s && wget https://github.com/%s/%s/archive/%s.zip"
	unzipRepoArchive = "%s && unzip %s.zip && cd %s && mv * ../ && cd .. && rm %s.zip && rm -rf %s"
	addFiles         = "%s && git add . "
	commitFiles      = "%s && git commit -m \" %s \""
	pushFiles        = "%s && git push --set-upstream origin %s"
)

// SubVersionPackageModel creates a github repo for the packageModel on github.com/gophr/gophr-packages
// versioned a the speicifed ref
func SubVersionPackageModel(packageModel *models.PackageModel, ref string) error {
	log.Printf("Preparing to sub-version %s/%s@%s \n", *packageModel.Author, *packageModel.Repo, ref)
	// If the given ref is empty or refers to 'master' then we need to grab the current master SHA
	if len(ref) == 0 || ref == "master" {
		log.Println("Ref is empty or is 'master', fetching current master SHA")
		curretRef, err := common.FetchRefs(*packageModel.Author, *packageModel.Repo)
		if err != nil || len(curretRef.MasterRefHash) == 0 {
			return fmt.Errorf(
				"Error could not retrieve master ref of %s/%s, or packageModel does not exist",
				*packageModel.Author,
				*packageModel.Repo,
			)
		}
		ref = curretRef.MasterRefHash
	}

	// First check if this ref has already been versioned for this packageModel
	log.Printf("Checking if ref %s has been versioned before \n", ref)
	exists, err := github.CheckIfRefExists(*packageModel.Author, *packageModel.Repo, ref)
	if exists == true && err == nil {
		log.Println("That ref has already been versioned")
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error occured in checking if ref exists. %s", err)
	}

	log.Printf("%s/%s@%s has not been versioned yet",
		github.GitHubGophrPackageOrgName,
		github.BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo),
		github.BuildGitHubBranch(ref),
	)

	// Set working folderName for package
	folderName = github.BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo)

	// Instantiate New Github Request Service
	log.Println("Initializing gitHub component")
	gitHubRequestService := github.NewGitHubRequestService()

	// Prepare to Create a new Github repo for packageModel if DNE
	log.Printf("Creating new Github repo for %s/%s at %s",
		*packageModel.Author,
		*packageModel.Repo,
		ref,
	)
	err = gitHubRequestService.CreateNewGitHubRepo(*packageModel)
	// TODO:(Shikkic) figure out
	log.Printf("%s", err)

	log.Printf("Initializing folder and initializing git repo for %s \n", folderName)
	if err = initializeRepoCMD(packageModel); err != nil {
		checkError(err, folderName)
		return fmt.Errorf("Error occured in initializing the git repo. %s", err)
	}

	log.Printf("Creating branch %s \n", github.BuildGitHubBranch(ref))
	if err = createBranchCMD(packageModel, ref); err != nil {
		checkError(err, folderName)
		return fmt.Errorf("Error occured in creating git branch. %s", err)
	}

	log.Printf("Setting remote branch url %s \n", github.BuildGitHubBranch(ref))
	if err = setRemoteCMD(packageModel, ref); err != nil {
		checkError(err, folderName)
		return fmt.Errorf("Error occured in setting remote url. %s", err)
	}

	log.Printf("Fetching github archive for %s/%s with tag %s \n",
		*packageModel.Author,
		*packageModel.Repo,
		ref,
	)
	if err = fetchArchiveCMD(packageModel, ref); err != nil {
		checkError(err, folderName)
		return fmt.Errorf("Error occured in fetching ref archive. %s", err)
	}

	log.Printf("Fetching github archive for %s/%s with tag %s \n",
		*packageModel.Author,
		*packageModel.Repo,
		ref,
	)
	if err = unzipArchiveCMD(packageModel, ref); err != nil {
		checkError(err, folderName)
		return fmt.Errorf("Error occured in unziping ref archive. %s", err)
	}

	// Fetch the timestamp of the ref commit
	commitDate, err := gitHubRequestService.FetchCommitTimestamp(packageModel, ref)
	if err != nil {
		return fmt.Errorf("Error occured in retrieving commit timestamp %s \n", err)
	}

	// Version lock all of the Github dependencies in the packageModel
	if err = verdeps.VersionDeps(
		verdeps.VersionDepsArgs{
			SHA:           ref,
			Path:          fmt.Sprintf("/tmp/%s", folderName),
			Date:          commitDate,
			Model:         packageModel,
			GithubService: gitHubRequestService,
		}); err != nil {
		return fmt.Errorf("Error occured in versioning deps %s \n", err)
	}

	log.Println("Adding unarchived repo data to branch")
	if err = addFilesCMD(); err != nil {
		checkError(err, folderName)
		return fmt.Errorf("Error occured in git add. %s", err)
	}

	log.Println("Commiting repo data to branch")
	if err = commitFilesCMD(packageModel, ref); err != nil {
		checkError(err, folderName)
		return fmt.Errorf("Error occured in commit repo data. %s", err)
	}

	log.Printf("Pushing files to branch %s \n", github.BuildRemoteURL(packageModel, ref))
	if err = pushFilesCMD(packageModel, ref); err != nil {
		checkError(err, folderName)
		return fmt.Errorf("Error occured in pushing files to branch. %s", err)
	}

	cleanUpExitHook(folderName)

	return nil
}

func initializeRepoCMD(packageModel *models.PackageModel) error {
	log.Println("Initializing folder and repo commmand")
	cmd := fmt.Sprintf(initalizeRepo, folderName, folderName)
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)
	return err
}

func createBranchCMD(packageModel *models.PackageModel, ref string) error {
	log.Println("Initializing folder and repo commmand")
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	cmd := fmt.Sprintf(createBranch, navigateFolderCMD, github.BuildGitHubBranch(ref))
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)
	return err
}

func setRemoteCMD(packageModel *models.PackageModel, ref string) error {
	log.Println("Initializing folder and repo commmand")
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	remoteURL := github.BuildRemoteURL(packageModel, ref)
	cmd := fmt.Sprintf(setRemoteCommand, navigateFolderCMD, remoteURL)
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)
	return err
}
func fetchArchiveCMD(packageModel *models.PackageModel, ref string) error {
	log.Println("Fetching and Unzipping Archive for tag")
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	cmd := fmt.Sprintf(fetchRepoArchive, navigateFolderCMD, *packageModel.Author, *packageModel.Repo, ref)
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)

	return err
}

func unzipArchiveCMD(packageModel *models.PackageModel, ref string) error {
	log.Println("Fetching and Unzipping Archive for tag")
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	zipFolder := *packageModel.Repo + "-" + ref
	cmd := fmt.Sprintf(unzipRepoArchive, navigateFolderCMD, ref, zipFolder, ref, zipFolder)
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)

	return err
}

func addFilesCMD() error {
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	cmd := fmt.Sprintf(addFiles, navigateFolderCMD)
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)
	return err
}

func commitFilesCMD(packageModel *models.PackageModel, ref string) error {
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	commitMessage := fmt.Sprintf("Gophr versioned repo of %s/%s@%s", *packageModel.Author, *packageModel.Repo, ref)
	cmd := fmt.Sprintf(commitFiles, navigateFolderCMD, commitMessage)
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s \n", out)
	return err
}

func pushFilesCMD(packageModel *models.PackageModel, ref string) error {
	navigateFolderCMD := fmt.Sprintf(navigateToPackageFolder, folderName)
	cmd := fmt.Sprintf(pushFiles, navigateFolderCMD, github.BuildGitHubBranch(ref))
	log.Println(cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	log.Printf("Output: %s 		\n", out)
	return err
}

// Helper functions

// Exit Hook to clean up files
func cleanUpExitHook(folderName string) {
	log.Printf("Exit hook initiated deleting folder %s \n", folderName)
	if strings.HasPrefix(folderName, "/") == true {
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
	}
}
