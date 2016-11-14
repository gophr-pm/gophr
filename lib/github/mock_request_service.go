package github

import (
	"time"

	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/stretchr/testify/mock"
)

// MockRequestService is a mock for RequestService.
type MockRequestService struct {
	mock.Mock
}

// NewMockRequestService creates a new MockRequestService.
func NewMockRequestService() *MockRequestService {
	return &MockRequestService{}
}

// FetchRepoData mocks RequestService.FetchCommitSHA.
func (m *MockRequestService) FetchRepoData(
	author string,
	repo string,
) (dtos.GithubRepo, error) {
	args := m.Called(author, repo)
	return args.Get(0).(dtos.GithubRepo), args.Error(1)
}

// FetchCommitSHA mocks RequestService.FetchCommitSHA.
func (m *MockRequestService) FetchCommitSHA(
	author string,
	repo string,
	timestamp time.Time,
) (string, error) {
	args := m.Called(author, repo, timestamp)
	return args.String(0), args.Error(1)
}

// ExpandPartialSHA mocks RequestService.ExpandPartialSHA.
func (m *MockRequestService) ExpandPartialSHA(
	args ExpandPartialSHAArgs,
) (string, error) {
	a := m.Called(args)
	return a.String(0), a.Error(1)
}

// FetchCommitTimestamp mocks RequestService.FetchCommitTimestamp.
func (m *MockRequestService) FetchCommitTimestamp(
	author string,
	repo string,
	sha string,
) (time.Time, error) {
	args := m.Called(author, repo, sha)
	return args.Get(0).(time.Time), args.Error(1)
}
