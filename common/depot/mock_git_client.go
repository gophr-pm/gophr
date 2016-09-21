package depot

import (
	git "github.com/libgit2/git2go"
	"github.com/stretchr/testify/mock"
)

// MockGitClient is a mock of gitClient
type MockGitClient struct {
	mock.Mock
}

// NewMockGitClient initialies a new mock gitClient that implements a GitClient interface.
func NewMockGitClient() *MockGitClient {
	return &MockGitClient{}
}

func (m *MockGitClient) InitRepo(archiveDirPath string, bare bool) (*git.Repository, error) {
	args := m.Called(archiveDirPath, bare)
	return args.Get(0).(*git.Repository), args.Error(1)
}

func (m *MockGitClient) CreateIndex(repo *git.Repository) (*git.Index, error) {
	args := m.Called(repo)
	return args.Get(0).(*git.Index), args.Error(1)
}

func (m *MockGitClient) IndexAddAll(index *git.Index) error {
	args := m.Called(index)
	return args.Error(0)
}

func (m *MockGitClient) WriteToIndexTree(index *git.Index, repo *git.Repository) (*git.Oid, error) {
	args := m.Called(index, repo)
	return args.Get(0).(*git.Oid), args.Error(1)
}

func (m *MockGitClient) WriteIndex(index *git.Index) error {
	args := m.Called(index)
	return args.Error(0)
}

func (m *MockGitClient) LookUpTree(repo *git.Repository, treeID *git.Oid) (*git.Tree, error) {
	args := m.Called(repo, treeID)
	return args.Get(0).(*git.Tree), args.Error(1)
}

func (m *MockGitClient) CreateCommit(
	repo *git.Repository,
	refname string,
	author *git.Signature,
	committer *git.Signature,
	message string,
	tree *git.Tree,
) error {
	args := m.Called(repo, refname, author, committer, message, tree)
	return args.Error(0)
}

func (m *MockGitClient) CreateRef(
	repo *git.Repository,
	name string,
	target string,
	force bool,
	message string,
) error {
	args := m.Called(repo, name, target, force, message)
	return args.Error(0)
}

func (m *MockGitClient) CheckoutHead(repo *git.Repository, opts *git.CheckoutOpts) error {
	args := m.Called(repo, opts)
	return args.Error(0)
}

func (m *MockGitClient) CreateRemote(
	repo *git.Repository,
	name string,
	url string,
) (*git.Remote, error) {
	args := m.Called(repo, name, url)
	return args.Get(0).(*git.Remote), args.Error(1)
}

func (m *MockGitClient) Push(
	remote *git.Remote,
	refspec []string,
	opts *git.PushOptions,
) error {
	args := m.Called(remote, refspec, opts)
	return args.Error(0)
}
