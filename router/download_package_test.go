package main

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDownloadPackage(t *testing.T) {
	mockIO := io.NewMockIO()
	mockIO.
		On("Mkdir", mock.AnythingOfType("string"), os.FileMode(0644)).
		Return(errors.New("this is an error"))
	args := packageDownloaderArgs{
		io:                   mockIO,
		author:               "myauthor",
		repo:                 "myrepo",
		sha:                  "mysha",
		constructionZonePath: "/my/cons/zone",
	}
	_, err := downloadPackage(args)
	assert.NotNil(t, err)
	mockIO.AssertExpectations(t)

	mockIO = io.NewMockIO()
	mockIO.
		On("Mkdir", mock.AnythingOfType("string"), os.FileMode(0644)).
		Return(nil)
	zipResp := &http.Response{
		StatusCode: 500,
		Body:       common.NewMockHTTPResponseBody(nil),
	}
	deleteWorkDirCalled := false
	args = packageDownloaderArgs{
		io:                   mockIO,
		author:               "myauthor",
		repo:                 "myrepo",
		sha:                  "mysha",
		constructionZonePath: "/my/cons/zone",

		doHTTPGet: func(url string) (*http.Response, error) {
			assert.Equal(t, "https://github.com/myauthor/myrepo/archive/mysha.zip", url)
			return zipResp, errors.New("this is an error")
		},
		deleteWorkDir: func(folderPath string) {
			deleteWorkDirCalled = true
		},
	}
	_, err = downloadPackage(args)
	assert.NotNil(t, err)
	mockIO.AssertExpectations(t)
	assert.True(t, zipResp.Body.(*common.MockHTTPResponseBody).WasClosed())
	assert.True(t, deleteWorkDirCalled)

	mockIO = io.NewMockIO()
	mockIO.
		On("Mkdir", mock.AnythingOfType("string"), os.FileMode(0644)).
		Return(nil)
	zipResp = &http.Response{
		StatusCode: 404,
		Body:       common.NewMockHTTPResponseBody(nil),
	}
	deleteWorkDirCalled = false
	args = packageDownloaderArgs{
		io:                   mockIO,
		author:               "myauthor",
		repo:                 "myrepo",
		sha:                  "mysha",
		constructionZonePath: "/my/cons/zone",

		doHTTPGet: func(url string) (*http.Response, error) {
			assert.Equal(t, "https://github.com/myauthor/myrepo/archive/mysha.zip", url)
			return zipResp, nil
		},
		deleteWorkDir: func(folderPath string) {
			deleteWorkDirCalled = true
		},
	}
	_, err = downloadPackage(args)
	assert.NotNil(t, err)
	mockIO.AssertExpectations(t)
	assert.True(t, zipResp.Body.(*common.MockHTTPResponseBody).WasClosed())
	assert.True(t, deleteWorkDirCalled)

	mockIO = io.NewMockIO()
	mockIO.
		On("Mkdir", mock.AnythingOfType("string"), os.FileMode(0644)).
		Return(nil)
	mockIO.
		On("Create", mock.AnythingOfType("string")).
		Return((*os.File)(nil), errors.New("oh no"))
	zipResp = &http.Response{
		StatusCode: 200,
		Body:       common.NewMockHTTPResponseBody([]byte("this is a zip")),
	}
	deleteWorkDirCalled = false
	args = packageDownloaderArgs{
		io:                   mockIO,
		author:               "myauthor",
		repo:                 "myrepo",
		sha:                  "mysha",
		constructionZonePath: "/my/cons/zone",

		doHTTPGet: func(url string) (*http.Response, error) {
			assert.Equal(t, "https://github.com/myauthor/myrepo/archive/mysha.zip", url)
			return zipResp, nil
		},
		deleteWorkDir: func(folderPath string) {
			deleteWorkDirCalled = true
		},
	}
	_, err = downloadPackage(args)
	assert.NotNil(t, err)
	mockIO.AssertExpectations(t)
	assert.True(t, zipResp.Body.(*common.MockHTTPResponseBody).WasClosed())
	assert.True(t, deleteWorkDirCalled)
}
