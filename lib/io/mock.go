package io

import (
	"io"
	"os"

	"github.com/stretchr/testify/mock"
)

// MockIO is the mock of the IO interface.
type MockIO struct {
	mock.Mock
}

// Mkdir mocks os.Mkdir.
func (m *MockIO) Mkdir(name string, perm os.FileMode) error {
	args := m.Called(name, perm)
	return args.Error(0)
}

// Create mocks os.Create.
func (m *MockIO) Create(name string) (*os.File, error) {
	args := m.Called(name)
	return args.Get(0).(*os.File), args.Error(1)
}

// Stat mocks os.Stat.
func (m *MockIO) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)
	return args.Get(0).(os.FileInfo), args.Error(1)
}

// Copy mocks io.Copy.
func (m *MockIO) Copy(dst io.Writer, src io.Reader) (int64, error) {
	args := m.Called(dst, src)
	return args.Get(0).(int64), args.Error(1)
}

// ReadDir mocks io.ReadDir.
func (m *MockIO) ReadDir(dirname string) ([]os.FileInfo, error) {
	args := m.Called(dirname)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

// NewMockIO creates a new io mock.
func NewMockIO() *MockIO {
	return &MockIO{}
}
