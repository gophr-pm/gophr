package db

import "github.com/stretchr/testify/mock"

// MockBatch mocks Batch.
type MockBatch struct {
	mock.Mock
}

// NewMockBatch creates a new MockBatch.
func NewMockBatch() *MockBatch {
	return &MockBatch{}
}

// Query adds the query to the batch operation
func (m *MockBatch) Query(stmt string, values ...interface{}) {
	var all []interface{}
	all = append(all, stmt)
	all = append(all, values...)

	m.Called(all...)
}

// Execute executes the batch operation.
func (m *MockBatch) Execute() error {
	args := m.Called()
	return args.Error(0)
}
