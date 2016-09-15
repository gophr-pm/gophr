package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	git "github.com/libgit2/git2go"
	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/github"
	"github.com/skeswa/gophr/common/verdeps"
)

var (
	folderName string
	folderPath string
	// TODO(skeswa): this needs to be configurable.
	commitAuthor = "gophrpm"
	// TODO(skeswa): this needs to be configurable.
	commitAuthorEmail  = "gophr.pm@gmail.com"
	gitHubRemoteOrigin = "git@github.com:gophr-packages/%s.git"
)

// versionAndArchivePackage creates a github repo for the packageModel on
// github.com/gophr/gophr-packages versioned a the specified args.sha.
func versionAndArchivePackage(args packageVersionerArgs) error {
	// If the given args.sha is empty or args.shaers to 'master' then we need to grab the current master SHA
	log.Printf("Preparing to sub-version %s/%s@%s \n", args.author, args.repo, args.sha)
	if len(args.sha) == 0 || args.sha == "master" {
		log.Println("Ref is empty or is 'master', fetching current master SHA")
		curretRef, err := common.FetchRefs(args.author, args.repo)
		if err != nil || len(curretRef.MasterRefHash) == 0 {
			return fmt.Errorf(
				"Error could not retrieve master args.sha of %s/%s, or packageModel does not exist \n",
				args.author,
				args.repo,
			)
		}
		args.sha = curretRef.MasterRefHash
	}

	// First check if this args.sha has already been versioned for this packageModel
	log.Printf("Checking if args.sha %s has been versioned before \n", args.sha)
	exists, err := github.CheckIfRefExists(args.author, args.repo, args.sha)
	if exists == true && err == nil {
		log.Println("That args.sha has already been versioned")
		// Since we wouldn't have gotten this far if this were already recorded,
		// make sure that we record it now.
		go args.recordPackageArchival(packageArchivalArgs{
			db:     args.db,
			sha:    args.sha,
			repo:   args.repo,
			author: args.author,
		})

		return nil
	}
	if err != nil {
		return fmt.Errorf("Error occurred in checking if args.sha exists. %s", err)
	}

	log.Printf("%s/%s@%s has not been versioned yet",
		github.GitHubGophrPackageOrgName,
		github.BuildNewGitHubRepoName(args.author, args.repo),
		github.BuildGitHubBranch(args.sha),
	)

	// Set working folderName and folderPath for package
	folderName = github.BuildNewGitHubRepoName(args.author, args.repo)
	folderPath = filepath.Join(args.constructionZonePath, folderName)

	// Fetch args.sha archive
	refZipURL := fmt.Sprintf("https://github.com/%s/%s/archive/%s.zip", args.author, args.repo, args.sha)
	resp, err := http.Get(refZipURL)
	if err != nil || resp.StatusCode == 404 {
		// TODO:(Shikkic) Better error description here
		return fmt.Errorf("Error 404, could not find args.sha archive for %s. %v \n", refZipURL, err)
	}
	defer resp.Body.Close()

	// Write Archive to filepath
	zipFilePath := fmt.Sprintf("%s/%s.zip", args.constructionZonePath, args.sha)
	out, err := os.Create(zipFilePath)
	if err != nil {
		if deletionErr := deleteAchriveFile(zipFilePath); deletionErr != nil {
			return fmt.Errorf("Error, could not write args.sha archive to file system or delete archive. %v, %v \n", err, deletionErr)
		}
		return fmt.Errorf("Error, could not write args.sha archive to file system. %v \n", err)
	}
	defer out.Close()
	io.Copy(out, resp.Body)

	// Unzip files
	if err = unzip(zipFilePath, args.constructionZonePath); err != nil {
		if deletionErr := deleteAchriveFile(zipFilePath); deletionErr != nil {
			return fmt.Errorf("Error, could not unzip args.sha archive or delete it. %v, %v. \n", err, deletionErr)
		}
		return fmt.Errorf("Error, could not unzip args.sha archive. %v \n", err)
	}

	// Delete The Archive File
	if deletionErr := deleteAchriveFile(zipFilePath); deletionErr != nil {
		return deletionErr
	}

	// Move files around
	targetFolder := fmt.Sprintf(
		"%s/%s-%s",
		args.constructionZonePath,
		args.repo,
		args.sha)
	newTargetFolder := fmt.Sprintf(
		"%s/%s",
		args.constructionZonePath,
		github.BuildNewGitHubRepoName(args.author, args.repo))
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
	gitHubRequestService := github.NewRequestService(
		github.RequestServiceParams{
			ForIndexer: false,
			Conf:       args.conf,
			Session:    args.db,
		},
	)

	// Prepare to Create a new Github repo for packageModel if DNE
	log.Println("Create new repo")
	err = gitHubRequestService.CreateNewRepo(args.author, args.repo)

	// Fetch the timestamp of the args.sha commit
	commitDate, err := gitHubRequestService.FetchCommitTimestamp(args.author, args.repo, args.sha)
	if err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error could not fetch commit timestamp or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error could not fetch commit timestamp %s \n", err)
	}

	// Version lock all of the Github dependencies in the packageModel
	if err = verdeps.VersionDeps(
		verdeps.VersionDepsArgs{
			SHA:           args.sha,
			Repo:          args.repo,
			Path:          folderPath,
			Date:          commitDate,
			Author:        args.author,
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
			args.author,
			args.repo,
			args.sha,
		),
		tree,
	)
	log.Println("Created commit", commitID)
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
	log.Println("Building the github branch")
	branchName := github.BuildGitHubBranch(args.sha)
	branch, err := repo.CreateBranch(branchName, headCommit, false)
	if err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error, could not create branch or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not create branch. %v \n", err)
	}

	log.Println("Setting the upstream")
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
		return fmt.Errorf("Error, could not create symbolic args.sha. %v \n", err)
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
			"http://%s:%s/%s.git",
			"depot-svc",
			"3000",
			github.BuildNewGitHubRepoName(args.author, args.repo),
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
			// TODO(skeswa): optimize out this closure.
			CredentialsCallback: func(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
				// TODO(shikkic): figure out how to get ssh working
				//ret, cred := git.NewCredSshKey("git", "/Users/shikkic/.ssh/id_rsa.pub", "/Users/shikkic/.ssh/id_rsa", "")
				ret, cred := git.NewCredUserpassPlaintext(
					args.creds.GithubPush.User,
					args.creds.GithubPush.Pass)

				return git.ErrorCode(ret), &cred
			},
			CertificateCheckCallback: certificateCheckCallback,
		},
	}

	log.Println("Doing the push")
	if err = remote.Push([]string{"refs/heads/" + branchName + ":refs/heads/" + branchName}, pushOptions); err != nil {
		if deletionErr := deleteFolder(folderPath); deletionErr != nil {
			return fmt.Errorf("Error,could not push to remote or delete repo folder. %v, %v \n", deletionErr, err)
		}
		return fmt.Errorf("Error, could not push to remote. %v \n", err)
	}

	// Delete work dir before returning
	log.Println("Deleting the folder")
	if deletionErr := deleteFolder(folderPath); deletionErr != nil {
		return fmt.Errorf("Error, could not delete repo folder and clean work dir. %v \n", deletionErr)
	}

	// Record that this package has been archived.
	go args.recordPackageArchival(packageArchivalArgs{
		db:     args.db,
		sha:    args.sha,
		repo:   args.repo,
		author: args.author,
	})

	return nil
}
