package io

import (
	"os"
	"time"
)

// FakeFileInfo is a struct used to mock os.FileInfo.
type FakeFileInfo struct {
	SysProp      interface{}
	NameProp     string
	SizeProp     int64
	IsDirProp    bool
	ModTimeProp  time.Time
	FileModeProp os.FileMode
}

// NewFakeFileInfo creates a new, shallowly defined FakeFileInfo.
func NewFakeFileInfo(name string, size int64, isDir bool) FakeFileInfo {
	return FakeFileInfo{
		NameProp:  name,
		SizeProp:  size,
		IsDirProp: isDir,
	}
}

// Name returns the name.
func (m FakeFileInfo) Name() string {
	return m.NameProp
}

// Size returns the size.
func (m FakeFileInfo) Size() int64 {
	return m.SizeProp
}

// Mode returns the mode.
func (m FakeFileInfo) Mode() os.FileMode {
	return m.FileModeProp
}

// ModTime returns the modTime.
func (m FakeFileInfo) ModTime() time.Time {
	return m.ModTimeProp
}

// IsDir returns the isDir.
func (m FakeFileInfo) IsDir() bool {
	return m.IsDirProp
}

// Sys returns the sys.
func (m FakeFileInfo) Sys() interface{} {
	return m.SysProp
}
