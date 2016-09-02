package subv

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gocql/gocql"
	git "github.com/libgit2/git2go"
	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/github"
	"github.com/skeswa/gophr/common/models"
	"github.com/skeswa/gophr/common/verdeps"
)

var (
	folderName         string
	folderPath         string
	gitHubRemoteOrigin = "git@github.com:gophr-packages/%s.git"
	commitAuthor       = "gophrpm"
	commitAuthorEmail  = "gophr.pm@gmail.com"
)

// SubVersionPackageModel creates a github repo for the packageModel on github.com/gophr/gophr-packages
// versioned a the speicifed ref.
func SubVersionPackageModel(session *gocql.Session, packageModel *models.PackageModel, ref string, fileDir string) error {
	log.Printf("Preparing to sub-version %s/%s@%s \n", *packageModel.Author, *packageModel.Repo, ref)
	// If the given ref is empty or refers to 'master' then we need to grab the current master SHA
	if len(ref) == 0 || ref == "master" {
		log.Println("Ref is empty or is 'master', fetching current master SHA")
		curretRef, err := common.FetchRefs(*packageModel.Author, *packageModel.Repo)
		if err != nil || len(curretRef.MasterRefHash) == 0 {
			return fmt.Errorf(
				"Error could not retrieve master ref of %s/%s, or packageModel does not exist \n",
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

	// Set working folderName and folderPath for package
	folderName = github.BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo)
	folderPath = fmt.Sprintf("%s/%s", fileDir, folderName)

	// Fetch ref archive
	refZipURL := fmt.Sprintf("https://github.com/%s/%s/archive/%s.zip", *packageModel.Author, *packageModel.Repo, ref)
	resp, err := http.Get(refZipURL)
	if err != nil || resp.StatusCode == 404 {
		// TODO:(Shikkic) Better error description here
		return fmt.Errorf("Error 404, could not find ref archive for %s. %v \n", refZipURL, err)
	}
	defer resp.Body.Close()

	// Write Archive to filepath
	zipFilePath := fmt.Sprintf("%s/%s.zip", fileDir, ref)
	out, err := os.Create(zipFilePath)
	if err != nil {
		if deletionErr := deleteAchriveFile(zipFilePath); deletionErr != nil {
			return fmt.Errorf("Error, could not write ref archive to file system or delete archive. %v, %v \n", err, deletionErr)
		}
		return fmt.Errorf("Error, could not write ref archive to file system. %v \n", err)
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	// Unzip files
	if err = unzip(zipFilePath, fileDir); err != nil {
		if deletionErr := deleteAchriveFile(zipFilePath); deletionErr != nil {
			return fmt.Errorf("Error, could not unzip ref archive or delete it. %v, %v. \n", err, deletionErr)
		}
		return fmt.Errorf("Error, could not unzip ref archive. %v \n", err)
	}

	// Delete The Archive File
	if deletionErr := deleteAchriveFile(zipFilePath); deletionErr != nil {
		return deletionErr
	}

	// Move files around
	targetFolder := fmt.Sprintf("%s/%s-%s", fileDir, *packageModel.Repo, ref)
	newTargetFolder := fmt.Sprintf("%s/%s", fileDir, github.BuildNewGitHubRepoName(*packageModel.Author, *packageModel.Repo))
	if err = os.Rename(targetFolder, newTargetFolder); err != nil {
		if deletionErr := deleteFolder(newTargetFolder); deletionErr != nil {
			return fmt.Errorf("Error, could not rename archive folder to target folder or delete it. %v %v. \n", err, deletionErr)
		}
		return fmt.Errorf("Error, could not rename archive folder to target folder. %v \n", err)
	}

	// Git init
	repo, err := git.InitRepository(newTargetFolder, false)
	if err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not initialize new repository or delete repo folder. %v, %v \n", deletionErr, err)
		}
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
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error could not fetch commit timestamp or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error could not fetch commit timestamp %s \n", err)
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
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error could not version deps properly or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error could not version deps properly. %v \n", err)
	}

	// Git add all
	index, err := repo.Index()
	if err = index.AddAll([]string{}, git.IndexAddDefault, nil); err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not add files to git repo or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not add files to git repo. %v \n", err)
	}

	// Write tree
	treeID, err := index.WriteTreeTo(repo)
	if err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not write tree or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not write tree. %v \n", err)
	}

	// Write the index
	if err = index.Write(); err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not write index or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not write index. %v \n", err)
	}
	// TODO(Shikkic): is this necessary here?
	tree, err := repo.LookupTree(treeID)

	// Create commit Signature
	sig := &git.Signature{
		Name:  commitAuthor,
		Email: commitAuthorEmail,
		When:  time.Now(),
	}

	// Create commit
	// TODO(Shikkic): is commitID necessary here?
	// We dont use it
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
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not commit data or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not commit data. %v \n", err)
	}

	// Lookup Current Commit
	// TODO(Shikkic): dont think this is necessary
	head, err := repo.Head()
	if err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not look up repo HEAD or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not look up repo HEAD. %v \n", err)
	}
	headCommit, err := repo.LookupCommit(head.Target())
	if err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not get HEAD commit or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not get HEAD commit. %v \n", err)
	}

	// Creating branch
	branchName := github.BuildGitHubBranch(ref)
	branch, err := repo.CreateBranch(branchName, headCommit, false)
	if err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not create branch or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not create branch. %v \n", err)
	}

	if err = branch.SetUpstream(branchName); err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not set upstream branch or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not set upstream branch. %v \n", err)
	}

	_, err = repo.References.CreateSymbolic("HEAD", fmt.Sprintf("refs/heads/%s", branchName), true, "headOne")
	if err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not create symbolic ref or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not create symbolic ref. %v \n", err)
	}

	// Check out Branch
	opts := &git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing,
	}
	if err = repo.CheckoutHead(opts); err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not checkout branch or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not checkout branch. %v \n", err)
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
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not create remote origin or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not create remote origin. %v \n", err)
	}

	// Define push options
	pushOptions := &git.PushOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback:      credentialsCallback,
			CertificateCheckCallback: certificateCheckCallback,
		},
	}

	if err = remote.Push([]string{"refs/heads/" + branchName + ":refs/heads/" + branchName}, pushOptions); err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error,could not push to remote or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not push to remote. %v \n", err)
	}

	// Delete work dir before returning
	if deletionErr := deleteFolder(folderPath); deletionErr != nil {
		return fmt.Errorf("Error, could not delete repo folder and clean work dir. %v \n", deletionErr)
	}

	// Record that this package has been archived.
	go models.RecordPackageArchival(
		session,
		*(packageModel.Author),
		*(packageModel.Repo),
		ref,
	)

	return nil
}
