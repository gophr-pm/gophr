package main

import (
	git "github.com/libgit2/git2go"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mock of client
type MockClient struct {
	mock.Mock
}

// NewMockClient initialies a new mock client that implements a Client interface.
func NewMockClient() *MockClient {
	return &MockClient{}
}

// InitRepo mocks GitClint#InitRepo.
func (m *MockClient) InitRepo(archiveDirPath string, bare bool) (*git.Repository, error) {
	args := m.Called(archiveDirPath, bare)
	return args.Get(0).(*git.Repository), args.Error(1)
}

// CreateIndex mocks GitClint#CreateIndex.
func (m *MockClient) CreateIndex(repo *git.Repository) (*git.Index, error) {
	args := m.Called(repo)
	return args.Get(0).(*git.Index), args.Error(1)
}

// IndexAddAll mocks GitClint#IndexAddAll.
func (m *MockClient) IndexAddAll(index *git.Index) error {
	args := m.Called(index)
	return args.Error(0)
}

// WriteToIndexTree mocks GitClint#WriteToIndexTree.
func (m *MockClient) WriteToIndexTree(index *git.Index, repo *git.Repository) (*git.Oid, error) {
	args := m.Called(index, repo)
	return args.Get(0).(*git.Oid), args.Error(1)
}

// WriteIndex mocks GitClint#WriteIndex.
func (m *MockClient) WriteIndex(index *git.Index) error {
	args := m.Called(index)
	return args.Error(0)
}

// LookUpTree mocks GitClint#LookUpTree.
func (m *MockClient) LookUpTree(repo *git.Repository, treeID *git.Oid) (*git.Tree, error) {
	args := m.Called(repo, treeID)
	return args.Get(0).(*git.Tree), args.Error(1)
}

// CreateCommit mocks GitClint#CreateCommit.
func (m *MockClient) CreateCommit(
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

// CreateRef mocks GitClint#CreateRef.
func (m *MockClient) CreateRef(
	repo *git.Repository,
	name string,
	target string,
	force bool,
	message string,
) error {
	args := m.Called(repo, name, target, force, message)
	return args.Error(0)
}

// CheckoutHead mocks GitClint#CheckoutHead.
func (m *MockClient) CheckoutHead(repo *git.Repository, opts *git.CheckoutOpts) error {
	args := m.Called(repo, opts)
	return args.Error(0)
}

// CreateRemote mocks GitClint#CreateRemote.
func (m *MockClient) CreateRemote(
	repo *git.Repository,
	name string,
	url string,
) (*git.Remote, error) {
	args := m.Called(repo, name, url)
	return args.Get(0).(*git.Remote), args.Error(1)
}

// Push mocks GitClint#Push.
func (m *MockClient) Push(
	remote *git.Remote,
	refspec []string,
	opts *git.PushOptions,
) error {
	args := m.Called(remote, refspec, opts)
	return args.Error(0)
}
