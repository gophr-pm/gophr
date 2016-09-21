package io

import (
	"os"

	"github.com/stretchr/testify/mock"
)

// MockIO is the mock of the IO interface.
type MockIO struct {
	mock.Mock
}

// Mkdir creates a new directory with the specified name and permission bits.
// If there is an error, it will be of type *PathError.
func (io *MockIO) Mkdir(name string, perm os.FileMode) error {
	args := io.Called(name, perm)
	return args.Error(0)
}

// NewMockIO creates a new IO mock.
func NewMockIO() *MockIO {
	return &MockIO{}
}
