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

// FetchCommitSHA mocks MockRequestService.FetchCommitSHA.
func (m *MockRequestService) FetchCommitSHA(
	author string,
	repo string,
	timestamp time.Time,
) (string, error) {
	args := m.Called(author, repo, timestamp)
	return args.String(0), args.Error(1)
}

// FetchCommitTimestamp mocks MockRequestService.FetchCommitTimestamp.
func (m *MockRequestService) FetchCommitTimestamp(
	author string,
	repo string,
	sha string,
) (time.Time, error) {
	args := m.Called(author, repo, sha)
	return args.Get(0).(time.Time), args.Error(1)
}

// FetchGitHubDataForPackageModel mocks
// MockRequestService.FetchGitHubDataForPackageModel.
func (m *MockRequestService) FetchGitHubDataForPackageModel(
	author string,
	repo string,
) (dtos.GithubRepo, error) {
	args := m.Called(author, repo)
	return args.Get(0).(dtos.GithubRepo), args.Error(1)
}
