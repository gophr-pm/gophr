package io

import (
	"io"
	"io/ioutil"
	"os"
)

type realIO struct{}

// Mkdir calls os.Mkdir.
func (r *realIO) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// Create calls os.Create.
func (r *realIO) Create(name string) (*os.File, error) {
	return os.Create(name)
}

// Stat calls os.Stat
func (r *realIO) Stat(dirname string) (os.FileInfo, error) {
	return os.Stat(dirname)
}

// Copy calls io.Copy.
func (r *realIO) Copy(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}

// ReadDir calls ioutil.ReadDir.
func (r *realIO) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

// NewIO creates a new IO.
func NewIO() IO {
	return &realIO{}
}
