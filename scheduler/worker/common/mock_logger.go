package common

import "github.com/stretchr/testify/mock"

// MockJobLogger is the mock implementation of the job logger.
type MockJobLogger struct {
	mock.Mock
}

// NewMockJobLogger creates a new MockJobLogger.
func NewMockJobLogger() *MockJobLogger {
	return &MockJobLogger{}
}

// Info logs the args at the info log level.
func (m *MockJobLogger) Info(args ...interface{}) {
	m.Called(args...)
}

// Infof logs the template at the info log level.
func (m *MockJobLogger) Infof(template string, args ...interface{}) {
	var adaptedArgs []interface{}
	adaptedArgs = append(adaptedArgs, template)
	adaptedArgs = append(adaptedArgs, args...)

	m.Called(adaptedArgs...)
}

// Error logs the args at the error log level.
func (m *MockJobLogger) Error(args ...interface{}) {
	m.Called(args...)
}

// Errorf logs the template at the error log level.
func (m *MockJobLogger) Errorf(template string, args ...interface{}) {
	var adaptedArgs []interface{}
	adaptedArgs = append(adaptedArgs, template)
	adaptedArgs = append(adaptedArgs, args...)

	m.Called(adaptedArgs...)
}

// Start logs that the job started.
func (m *MockJobLogger) Start() {
	m.Called()
}

// Finish logs that the job finished.
func (m *MockJobLogger) Finish() {
	m.Called()
}
