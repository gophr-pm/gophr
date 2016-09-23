package git

import git "github.com/libgit2/git2go"

// Client is the external interface of client.
type Client interface {
	InitRepo(archiveDirPath string, bare bool) (*git.Repository, error)
	CreateIndex(repo *git.Repository) (*git.Index, error)
	IndexAddAll(index *git.Index) error
	WriteToIndexTree(index *git.Index, repo *git.Repository) (*git.Oid, error)
	WriteIndex(index *git.Index) error
	LookUpTree(repo *git.Repository, treeID *git.Oid) (*git.Tree, error)
	CreateCommit(
		repo *git.Repository,
		refname string,
		author *git.Signature,
		committer *git.Signature,
		message string,
		tree *git.Tree,
	) error
	CreateRef(
		repo *git.Repository,
		name string,
		target string,
		force bool,
		message string,
	) error
	CheckoutHead(repo *git.Repository, opts *git.CheckoutOpts) error
	CreateRemote(repo *git.Repository, name string, url string) (*git.Remote, error)
	Push(remote *git.Remote, refspec []string, opts *git.PushOptions) error
}
