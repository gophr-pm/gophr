package io

import (
	"io"
	"io/ioutil"
	"os"
)

type ioImpl struct{}

// Mkdir calls os.Mkdir.
func (r *ioImpl) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// Create calls os.Create.
func (r *ioImpl) Create(name string) (*os.File, error) {
	return os.Create(name)
}

// Stat calls os.Stat
func (r *ioImpl) Stat(dirname string) (os.FileInfo, error) {
	return os.Stat(dirname)
}

// Copy calls io.Copy.
func (r *ioImpl) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

// ReadDir calls ioutil.ReadDir.
func (r *ioImpl) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

// NewIO creates a new IO.
func NewIO() IO {
	return &ioImpl{}
}
