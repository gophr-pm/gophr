package datadog

import (
	"github.com/DataDog/datadog-go/statsd"
	"github.com/stretchr/testify/mock"
)

// MockDataDogClient is a mock for MockDataDogClient.
type MockDataDogClient struct {
	mock.Mock
}

// NewMockDataDogClient creates a new MockDataDogClient.
func NewMockDataDogClient() *MockDataDogClient {
	return &MockDataDogClient{}
}

// Gauge mocks MockDataDogClient.Gauge.
func (m *MockDataDogClient) Gauge(name string, value float64, tags []string, rate float64) error {
	args := m.Called(name, value, tags, rate)
	return args.Error(1)
}

// Event mocks MockDataDogClient.Event.
func (m *MockDataDogClient) Event(e *statsd.Event) error {
	args := m.Called(e)
	return args.Error(0)
}

// Incr mocks MockDataDogClient.Incr.
func (m *MockDataDogClient) Incr(name string, tags []string, rate float64) error {
	args := m.Called(name, tags, rate)
	return args.Error(0)
}
