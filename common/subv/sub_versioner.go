package subv

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	git "github.com/libgit2/git2go"
	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/github"
	"github.com/skeswa/gophr/common/models"
	"github.com/skeswa/gophr/common/verdeps"
)

var (
	folderName              string
	gitHubRemoteOrigin      = "git@github.com:gophr-packages/%s.git"
	navigateToPackageFolder = "cd /tmp/%s"
	filePath                = "/tmp"
	commitAuthor            = "gophrpm"
	commitAuthorEmail       = "gophr.pm@gmail.com"
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

	log.Printf("%s/%s@%s has not been versioned yet",
		github.GitHubGophrPackageOrgName,
		github.BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo),
		github.BuildGitHubBranch(ref),
	)

	// Set working folderName for package
	folderName = github.BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo)

	// Fetch ref archive
	refZipURL := fmt.Sprintf("https://github.com/%s/%s/archive/%s.zip", *packageModel.Author, *packageModel.Repo, ref)
	resp, err := http.Get(refZipURL)
	if err != nil || resp.StatusCode == 404 {
		// TODO:(Shikkic) Better error description here
		return fmt.Errorf("Error 404, could not find ref archive for %s. %v", refZipURL, err)
	}
	defer resp.Body.Close()

	// Write Archive to filepath
	zipFilePath := fmt.Sprintf("%s/%s.zip", filePath, ref)
	out, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("Error, could not write ref archive to file system. %v", err)
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	// Unzip files
	if err = unzip(zipFilePath, filePath); err != nil {
		return fmt.Errorf("Error, could not unzip ref archive. %v", err)
	}

	// Delete The Archive File
	if err = os.Remove(zipFilePath); err != nil {
		return fmt.Errorf("Error, could not delete ref archive file. %v", err)
	}

	// Move files around
	targetFolder := fmt.Sprintf("%s/%s-%s", filePath, *packageModel.Repo, ref)
	newTargetFolder := fmt.Sprintf("%s/%s", filePath, github.BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo))
	if err = os.Rename(targetFolder, newTargetFolder); err != nil {
		return fmt.Errorf("Error, could not rename archive folder to target folder. %v", err)
	}

	// Git init
	repo, err := git.InitRepository(newTargetFolder, false)
	if err != nil {
		log.Fatalln(err)
		return fmt.Errorf("Error, could not initialize new repository. %v", err)
	}

	// Instantiate New Github Request Service
	gitHubRequestService := github.NewRequestService()

	// Prepare to Create a new Github repo for packageModel if DNE
	err = gitHubRequestService.CreateNewGitHubRepo(*packageModel)
	// TODO:(Shikkic) figure out error handling here

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
		return fmt.Errorf("Error occured in versioning deps. %v \n", err)
	}

	// Git add all
	index, err := repo.Index()
	if err = index.AddAll([]string{}, git.IndexAddDefault, nil); err != nil {
		return fmt.Errorf("Error, could not add files to git repo. %v", err)
	}

	// Write tree
	treeID, err := index.WriteTreeTo(repo)
	if err != nil {
		return fmt.Errorf("Error, could not write tree. %v", err)
	}

	// Write the index
	if err = index.Write(); err != nil {
		return fmt.Errorf("Error, could not write index. %v", err)
	}
	// TODO is this necessary here?
	tree, err := repo.LookupTree(treeID)

	// Create commit Signature
	sig := &git.Signature{
		Name:  commitAuthor,
		Email: commitAuthorEmail,
		When:  time.Now(),
	}

	// Create commit
	// TODO is commitID necessary here?
	commitID, err := repo.CreateCommit(
		"HEAD",
		sig,
		sig,
		fmt.Sprintf("Gophr versioned repo of %s/%s@%s",
			*packageModel.Author,
			*packageModel.Repo,
			ref,
		),
		tree,
	)
	log.Println(commitID)
	if err != nil {
		return fmt.Errorf("Error, could not commit data. %v", err)
	}

	// Lookup Current Commit
	// TODO dont think this is necessary
	head, err := repo.Head()
	if err != nil {
		panic(err)
	}
	headCommit, err := repo.LookupCommit(head.Target())
	if err != nil {
		panic(err)
	}

	// Creating branch
	branchName := github.BuildGitHubBranch(ref)
	branch, err := repo.CreateBranch(branchName, headCommit, false)
	if err != nil {
		return fmt.Errorf("Error, could not create branch. %v", err)
	}

	if err = branch.SetUpstream(branchName); err != nil {
		return fmt.Errorf("Error, could not set upstream branch. %v", err)
	}

	_, err = repo.References.CreateSymbolic("HEAD", "refs/heads/"+branchName, true, "headOne")
	if err != nil {
		return fmt.Errorf("Error, could not create symbolic ref. %v", err)
	}

	// Check out Branch
	opts := &git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing,
	}
	if err = repo.CheckoutHead(opts); err != nil {
		return fmt.Errorf("Error, could not checkout branch. %v", err)
	}

	// Creating remote origin
	remote, err := repo.Remotes.Create(
		"origin",
		fmt.Sprintf(
			"https://github.com/gophr-packages/%s.git",
			github.BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo),
		),
	)
	if err != nil {
		return fmt.Errorf("Error, could not create remote origin. %v", err)
	}

	// Define push options
	pushOptions := &git.PushOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:      credentialsCallback,
			CertificateCheckCallback: certificateCheckCallback,
		},
	}

	if err := remote.Push([]string{"refs/heads/" + branchName + ":refs/heads/" + branchName}, pushOptions); err != nil {
		return fmt.Errorf("Error, could not push to remote. %v", err)
	}

	// TODO re-implement cleanUpExitHook
	//cleanUpExitHook(folderName)

	return nil
}
