package db

import "github.com/stretchr/testify/mock"

// MockClient mocks Client.
type MockClient struct {
	mock.Mock
}

// NewMockClient creates a new MockClient.
func NewMockClient() *MockClient {
	return &MockClient{}
}

// Close closes all connections. The client is unusable after this operation.
func (m *MockClient) Close() {
	m.Called()
}

// Query generates a new query object for interacting with the database.
// Further details of the query may be tweaked using the resulting query value
// before the query is executed. Query is automatically prepared if it has not
// previously been executed.
func (m *MockClient) Query(stmt string, values ...interface{}) Query {
	var all []interface{}
	all = append(all, stmt)
	all = append(all, values...)

	args := m.Called(all...)
	return args.Get(0).(Query)
}

// NewLoggedBatch creates a new logged batch.
func (m *MockClient) NewLoggedBatch() Batch {
	args := m.Called()
	return args.Get(0).(Batch)
}

// NewUnloggedBatch creates a new unlogged batch.
func (m *MockClient) NewUnloggedBatch() Batch {
	args := m.Called()
	return args.Get(0).(Batch)
}
