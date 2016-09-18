package main

import (
	"fmt"
	"time"

	git "github.com/libgit2/git2go"
	"github.com/skeswa/gophr/common/depot"
)

const (
	commitAuthor        = "Gophrs Archiver"
	commitAuthorEmail   = "archiver@gophr.pm"
	masterBranchName    = "master"
	masterBranchRef     = "refs/heads/master"
	masterPushDirective = "refs/heads/master:refs/heads/master"
)

func pushToDepot(args packagePusherArgs) error {
	// Initialize Git Repo
	repo, err := git.InitRepository(args.packagePaths.archiveDirPath, false)
	if err != nil {
		return fmt.Errorf("Could not initialize new repository: %v.", err)
	}

	// Git add all.
	index, err := repo.Index()
	if err = index.AddAll([]string{}, git.IndexAddDefault, nil); err != nil {
		return fmt.Errorf("Could not add files to git repo: %v.", err)
	}

	treeID, err := index.WriteTreeTo(repo)
	if err != nil {
		return fmt.Errorf("Could not write tree: %v.", err)
	}

	// Write the index
	if err = index.Write(); err != nil {
		return fmt.Errorf("Could not write index: %v.", err)
	}

	tree, err := repo.LookupTree(treeID)
	if err != nil {
		return fmt.Errorf("Could not retrieve repo tree: %v.", err)
	}

	// Create commit Signature
	sig := &git.Signature{
		Name:  commitAuthor,
		Email: commitAuthorEmail,
		When:  time.Now(),
	}

	if _, err = repo.CreateCommit(
		"HEAD",
		sig,
		sig,
		fmt.Sprintf("Gophr versioned repo %s/%s@%s",
			args.author,
			args.repo,
			args.sha,
		),
		tree,
	); err != nil {
		return fmt.Errorf("Could not commit data: %v.", err)
	}

	// Create ref for master
	if _, err = repo.References.CreateSymbolic(
		"HEAD",
		masterBranchRef,
		true,
		"headOne",
	); err != nil {
		return fmt.Errorf("Could not create ref for master: %v.", err)
	}

	// Check out master.
	if err = repo.CheckoutHead(&git.CheckoutOpts{
		Strategy: git.CheckoutSafe | git.CheckoutRecreateMissing,
	}); err != nil {
		return fmt.Errorf("Could not checkout master: %v.", err)
	}

	// Create remote origin.
	remote, err := repo.Remotes.Create(
		"origin",
		fmt.Sprintf(
			"http://%s/%s.git",
			depot.DepotInternalServiceAddress,
			depot.BuildHashedRepoName(args.author, args.repo, args.sha),
		),
	)
	if err != nil {
		return fmt.Errorf("Could not create remote origin: %v.", err)
	}

	// Define push options.
	pushOptions := &git.PushOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback: generateCredentialsCallback(
				args.creds.GithubPush.User,
				args.creds.GithubPush.Pass,
			),
			CertificateCheckCallback: certificateCheckCallback,
		},
	}

	if err = remote.Push([]string{masterPushDirective}, pushOptions); err != nil {
		return fmt.Errorf("Could not push to master: %v.", err)
	}

	return nil
}

func generateCredentialsCallback(user, pass string) func(string, string, git.CredType) (git.ErrorCode, *git.Cred) {
	return func(url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
		ret, cred := git.NewCredUserpassPlaintext(user, pass)

		return git.ErrorCode(ret), &cred
	}
}

func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	return 0
}
