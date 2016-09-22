package depot

import git "github.com/libgit2/git2go"

// NewGitClient initialies a new implementation of a GitClient interface.
func NewGitClient() GitClient {
	return &gitClient{}
}

// gitClient is responsible for managing all interal git2go functionality
// for depot.
type gitClient struct{}

// TODO(Shikkic): add comments
func (gc *gitClient) InitRepo(archiveDirPath string, bare bool) (*git.Repository, error) {
	repo, err := git.InitRepository(archiveDirPath, bare)
	return repo, err
}

func (gc *gitClient) CreateIndex(repo *git.Repository) (*git.Index, error) {
	index, err := repo.Index()
	return index, err
}

func (gc *gitClient) IndexAddAll(index *git.Index) error {
	err := index.AddAll([]string{}, git.IndexAddDefault, nil)
	return err
}

func (gc *gitClient) WriteToIndexTree(index *git.Index, repo *git.Repository) (*git.Oid, error) {
	treeID, err := index.WriteTreeTo(repo)
	return treeID, err
}

func (gc *gitClient) WriteIndex(index *git.Index) error {
	err := index.Write()
	return err
}

func (gc *gitClient) LookUpTree(repo *git.Repository, treeID *git.Oid) (*git.Tree, error) {
	tree, err := repo.LookupTree(treeID)
	return tree, err
}

func (gc *gitClient) CreateCommit(
	repo *git.Repository,
	refname string,
	author *git.Signature,
	committer *git.Signature,
	message string,
	tree *git.Tree,
) error {
	_, err := repo.CreateCommit(
		"HEAD",
		author,
		committer,
		message,
		tree,
	)
	return err
}

func (gc *gitClient) CreateRef(
	repo *git.Repository,
	name string,
	target string,
	force bool,
	message string,
) error {
	_, err := repo.References.CreateSymbolic(
		name,
		target,
		force,
		message,
	)
	return err
}

func (gc *gitClient) CheckoutHead(repo *git.Repository, opts *git.CheckoutOpts) error {
	err := repo.CheckoutHead(opts)
	return err
}

func (gc *gitClient) CreateRemote(
	repo *git.Repository,
	name string,
	url string,
) (*git.Remote, error) {
	remote, err := repo.Remotes.Create(
		"origin",
		url,
	)
	return remote, err
}

func (gc *gitClient) Push(
	remote *git.Remote,
	refspec []string,
	opts *git.PushOptions,
) error {
	err := remote.Push(refspec, opts)
	return err
}
