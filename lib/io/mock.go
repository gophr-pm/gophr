package io

import (
	"io"
	"os"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockIO is the mock of the IO interface.
type MockIO struct {
	mock.Mock
}

// Mkdir calls os.Mkdir.
func (m *MockIO) Mkdir(name string, perm os.FileMode) error {
	args := m.Called(name, perm)
	return args.Error(0)
}

// Create calls os.Create.
func (m *MockIO) Create(name string) (*os.File, error) {
	args := m.Called(name)
	return args.Get(0).(*os.File), args.Error(1)
}

// Copy calls io.Copy.
func (m *MockIO) Copy(dst io.Writer, src io.Reader) (int64, error) {
	args := m.Called(dst, src)
	return args.Get(0).(int64), args.Error(1)
}

// ReadDir calls ioutil.ReadDir.
func (m *MockIO) ReadDir(dirname string) ([]os.FileInfo, error) {
	args := m.Called(dirname)
	return args.Get(0).([]os.FileInfo), args.Error(1)
}

// NewMockIO creates a new IO mock.
func NewMockIO() *MockIO {
	return &MockIO{}
}

// MockFileInfo is a struct used to mock os.FileInfo.
type MockFileInfo struct {
	NameProp     string
	SizeProp     int64
	FileModeProp os.FileMode
	ModTimeProp  time.Time
	IsDirProp    bool
	SysProp      interface{}
}

// NewMockFileInfo creates a new, shallowly defined MockFileInfo.
func NewMockFileInfo(name string, size int64) MockFileInfo {
	return MockFileInfo{
		NameProp: name,
		SizeProp: size,
	}
}

// Name returns the name.
func (m MockFileInfo) Name() string {
	return m.NameProp
}

// Size returns the size.
func (m MockFileInfo) Size() int64 {
	return m.SizeProp
}

// Mode returns the mode.
func (m MockFileInfo) Mode() os.FileMode {
	return m.FileModeProp
}

// ModTime returns the modTime.
func (m MockFileInfo) ModTime() time.Time {
	return m.ModTimeProp
}

// IsDir returns the isDir.
func (m MockFileInfo) IsDir() bool {
	return m.IsDirProp
}

// Sys returns the sys.
func (m MockFileInfo) Sys() interface{} {
	return m.SysProp
}
