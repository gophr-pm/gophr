package query

import (
	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/stretchr/testify/mock"
)

// MockBatchingQueryable is the mock of the db.BatchingQueryable interface.
type MockBatchingQueryable struct {
	mock.Mock
}

// NewMockBatchingQueryable creates a new db.BatchingQueryable mock.
func NewMockBatchingQueryable() *MockBatchingQueryable {
	return &MockBatchingQueryable{}
}

// Query mocks db.BatchingQueryable.Query.
func (m *MockBatchingQueryable) Query(
	stmt string,
	values ...interface{},
) db.Query {
	argsSlice := []interface{}{stmt}
	argsSlice = append(argsSlice, values...)
	args := m.Called(argsSlice...)
	return args.Get(0).(db.Query)
}

// NewBatch mocks db.BatchingQueryable.NewBatch.
func (m *MockBatchingQueryable) NewBatch(bt gocql.BatchType) db.Batch {
	args := m.Called(bt)
	return args.Get(0).(db.Batch)
}

// ExecuteBatch mocks db.BatchingQueryable.ExecuteBatch.
func (m *MockBatchingQueryable) ExecuteBatch(b db.Batch) error {
	args := m.Called(b)
	return args.Error(0)
}
