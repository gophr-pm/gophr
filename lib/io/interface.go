package io

import (
	"io"
	"os"
)

// IO is an abstraction layer for all os-level I/O in gophr. This interface
// makes testing a great deal easier.
type IO interface {
	// Mkdir creates a new directory with the specified name and permission bits.
	// If there is an error, it will be of type *PathError.
	Mkdir(name string, perm os.FileMode) error
	// Create creates the named file with mode 0666 (before umask), truncating it
	// if it already exists. If successful, methods on the returned File can be
	// used for I/O; the associated file descriptor has mode O_RDWR. If there is
	// an error, it will be of type *PathError.
	Create(name string) (*os.File, error)
	// Copy copies from src to dst until either EOF is reached on src or an error
	// occurs. It returns the number of bytes copied and the first error
	// encountered while copying, if any.
	//
	// A successful Copy returns err == nil, not err == EOF. Because Copy is
	// defined to read from src until EOF, it does not treat an EOF from Read as
	// an error to be reported.
	//
	// If src implements the WriterTo interface, the copy is implemented by
	// calling src.WriteTo(dst). Otherwise, if dst implements the ReaderFrom
	// interface, the copy is implemented by calling dst.ReadFrom(src).
	Copy(dst io.Writer, src io.Reader) (written int64, err error)
	// ReadDir reads the directory named by dirname and returns a list of
	// directory entries sorted by filename.
	ReadDir(dirname string) ([]os.FileInfo, error)
	// Stat returns the FileInfo of a file
	Stat(dirname string) (os.FileInfo, error)
	// ReadFile reads the file named by filename and returns the contents. A
	// successful call returns err == nil, not err == EOF. Because ReadFile reads
	// the whole file, it does not treat an EOF from Read as an error to be
	// reported.
	ReadFile(filename string) ([]byte, error)
	// WriteFile writes data to a file named by filename. If the file does not
	// exist, WriteFile creates it with permissions perm; otherwise WriteFile
	// truncates it before writing.
	WriteFile(filename string, data []byte, perm os.FileMode) error
	// Rename renames (moves) oldpath to newpath. If newpath already exists,
	// Rename replaces it. OS-specific restrictions may apply when oldpath and
	// newpath are in different directories. If there is an error, it will be of
	// type *LinkError.
	Rename(oldpath, newpath string) error
}
