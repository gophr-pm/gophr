package query

import (
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/mock"
)

// MockBatchingQueryable is the mock of the BatchingQueryable interface.
type MockBatchingQueryable struct {
	mock.Mock
}

// NewMockBatchingQueryable creates a new BatchingQueryable mock.
func NewMockBatchingQueryable() *MockBatchingQueryable {
	return &MockBatchingQueryable{}
}

// Query mocks BatchingQueryable.Query.
func (m *MockBatchingQueryable) Query(
	stmt string,
	values ...interface{},
) *gocql.Query {
	argsSlice := []interface{}{stmt}
	argsSlice = append(argsSlice, values...)
	args := m.Called(argsSlice...)
	return args.Get(0).(*gocql.Query)
}

// NewBatch mocks BatchingQueryable.NewBatch.
func (m *MockBatchingQueryable) NewBatch(bt gocql.BatchType) *gocql.Batch {
	args := m.Called(bt)
	return args.Get(0).(*gocql.Batch)
}

// ExecuteBatch mocks BatchingQueryable.ExecuteBatch.
func (m *MockBatchingQueryable) ExecuteBatch(b *gocql.Batch) error {
	args := m.Called(b)
	return args.Error(0)
}
