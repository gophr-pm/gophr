package main

import (
	"errors"
	"testing"

	git "github.com/libgit2/git2go"
	"github.com/skeswa/gophr/common/config"
	g "github.com/skeswa/gophr/common/git"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPushtoDepot(t *testing.T) {
	mockGitClient := g.NewMockClient()
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, errors.New("this is an error"))
	args := packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
	}
	err := pushToDepot(args)
	assert.NotNil(t, err)

	mockGitClient = g.NewMockClient()
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, nil)
	mockGitClient.On("CreateIndex", &git.Repository{}).Return(&git.Index{}, errors.New("this is an error"))
	args = packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
	}
	err = pushToDepot(args)
	assert.NotNil(t, err)

	mockGitClient = g.NewMockClient()
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, nil)
	mockGitClient.On("CreateIndex", &git.Repository{}).Return(&git.Index{}, nil)
	mockGitClient.On("IndexAddAll", &git.Index{}).Return(errors.New("this is an error"))
	args = packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
	}
	err = pushToDepot(args)
	assert.NotNil(t, err)

	mockGitClient = g.NewMockClient()
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, nil)
	mockGitClient.On("CreateIndex", &git.Repository{}).Return(&git.Index{}, errors.New("this is an error"))
	args = packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
	}
	err = pushToDepot(args)
	assert.NotNil(t, err)

	mockGitClient = g.NewMockClient()
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, nil)
	mockGitClient.On("CreateIndex", &git.Repository{}).Return(&git.Index{}, nil)
	mockGitClient.On("IndexAddAll", &git.Index{}).Return(nil)
	mockGitClient.On("WriteToIndexTree", &git.Index{}, &git.Repository{}).Return(&git.Oid{}, errors.New("this is an error"))
	args = packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
	}
	err = pushToDepot(args)
	assert.NotNil(t, err)

	mockGitClient = g.NewMockClient()
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, nil)
	mockGitClient.On("CreateIndex", &git.Repository{}).Return(&git.Index{}, nil)
	mockGitClient.On("IndexAddAll", &git.Index{}).Return(nil)
	mockGitClient.On("WriteToIndexTree", &git.Index{}, &git.Repository{}).Return(&git.Oid{}, nil)
	mockGitClient.On("WriteIndex", &git.Index{}).Return(errors.New("this is an error"))
	args = packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
	}
	err = pushToDepot(args)
	assert.NotNil(t, err)

	mockGitClient = g.NewMockClient()
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, nil)
	mockGitClient.On("CreateIndex", &git.Repository{}).Return(&git.Index{}, nil)
	mockGitClient.On("IndexAddAll", &git.Index{}).Return(nil)
	mockGitClient.On("WriteToIndexTree", &git.Index{}, &git.Repository{}).Return(&git.Oid{}, nil)
	mockGitClient.On("WriteIndex", &git.Index{}).Return(nil)
	mockGitClient.On("LookUpTree", &git.Repository{}, &git.Oid{}).Return(&git.Tree{}, errors.New("this is an error"))
	args = packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
	}
	err = pushToDepot(args)
	assert.NotNil(t, err)

	mockGitClient = g.NewMockClient()
	sig := mock.AnythingOfType("*git.Signature")
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, nil)
	mockGitClient.On("CreateIndex", &git.Repository{}).Return(&git.Index{}, nil)
	mockGitClient.On("IndexAddAll", &git.Index{}).Return(nil)
	mockGitClient.On("WriteToIndexTree", &git.Index{}, &git.Repository{}).Return(&git.Oid{}, nil)
	mockGitClient.On("WriteIndex", &git.Index{}).Return(nil)
	mockGitClient.On("LookUpTree", &git.Repository{}, &git.Oid{}).Return(&git.Tree{}, nil)
	mockGitClient.On("CreateCommit", &git.Repository{}, "HEAD", sig, sig, "Gophr versioned repo authorName/repoName@repoSHA", &git.Tree{}).Return(errors.New("this is an error"))
	args = packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
	}
	err = pushToDepot(args)
	assert.NotNil(t, err)

	mockGitClient = g.NewMockClient()
	sig = mock.AnythingOfType("*git.Signature")
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, nil)
	mockGitClient.On("CreateIndex", &git.Repository{}).Return(&git.Index{}, nil)
	mockGitClient.On("IndexAddAll", &git.Index{}).Return(nil)
	mockGitClient.On("WriteToIndexTree", &git.Index{}, &git.Repository{}).Return(&git.Oid{}, nil)
	mockGitClient.On("WriteIndex", &git.Index{}).Return(nil)
	mockGitClient.On("LookUpTree", &git.Repository{}, &git.Oid{}).Return(&git.Tree{}, nil)
	mockGitClient.On("CreateCommit", &git.Repository{}, "HEAD", sig, sig, "Gophr versioned repo authorName/repoName@repoSHA", &git.Tree{}).Return(nil)
	mockGitClient.On("CreateRef", &git.Repository{}, "HEAD", "refs/heads/master", true, "headOne").Return(errors.New("this is an error"))
	args = packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
	}
	err = pushToDepot(args)
	assert.NotNil(t, err)

	mockGitClient = g.NewMockClient()
	sig = mock.AnythingOfType("*git.Signature")
	checkoutOpts := mock.AnythingOfType("*git.CheckoutOpts")
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, nil)
	mockGitClient.On("CreateIndex", &git.Repository{}).Return(&git.Index{}, nil)
	mockGitClient.On("IndexAddAll", &git.Index{}).Return(nil)
	mockGitClient.On("WriteToIndexTree", &git.Index{}, &git.Repository{}).Return(&git.Oid{}, nil)
	mockGitClient.On("WriteIndex", &git.Index{}).Return(nil)
	mockGitClient.On("LookUpTree", &git.Repository{}, &git.Oid{}).Return(&git.Tree{}, nil)
	mockGitClient.On("CreateCommit", &git.Repository{}, "HEAD", sig, sig, "Gophr versioned repo authorName/repoName@repoSHA", &git.Tree{}).Return(nil)
	mockGitClient.On("CreateRef", &git.Repository{}, "HEAD", "refs/heads/master", true, "headOne").Return(nil)
	mockGitClient.On("CheckoutHead", &git.Repository{}, checkoutOpts).Return(errors.New("this is an error"))
	args = packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
	}
	err = pushToDepot(args)
	assert.NotNil(t, err)

	mockGitClient = g.NewMockClient()
	sig = mock.AnythingOfType("*git.Signature")
	checkoutOpts = mock.AnythingOfType("*git.CheckoutOpts")
	remoteURL := mock.AnythingOfType("string")
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, nil)
	mockGitClient.On("CreateIndex", &git.Repository{}).Return(&git.Index{}, nil)
	mockGitClient.On("IndexAddAll", &git.Index{}).Return(nil)
	mockGitClient.On("WriteToIndexTree", &git.Index{}, &git.Repository{}).Return(&git.Oid{}, nil)
	mockGitClient.On("WriteIndex", &git.Index{}).Return(nil)
	mockGitClient.On("LookUpTree", &git.Repository{}, &git.Oid{}).Return(&git.Tree{}, nil)
	mockGitClient.On("CreateCommit", &git.Repository{}, "HEAD", sig, sig, "Gophr versioned repo authorName/repoName@repoSHA", &git.Tree{}).Return(nil)
	mockGitClient.On("CreateRef", &git.Repository{}, "HEAD", "refs/heads/master", true, "headOne").Return(nil)
	mockGitClient.On("CheckoutHead", &git.Repository{}, checkoutOpts).Return(nil)
	mockGitClient.On("CreateRemote", &git.Repository{}, "origin", remoteURL).Return(&git.Remote{}, errors.New("this is an error"))
	args = packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
	}
	err = pushToDepot(args)
	assert.NotNil(t, err)

	mockGitClient = g.NewMockClient()
	sig = mock.AnythingOfType("*git.Signature")
	checkoutOpts = mock.AnythingOfType("*git.CheckoutOpts")
	remoteURL = mock.AnythingOfType("string")
	refspec := mock.AnythingOfType("[]string")
	pushOpts := mock.AnythingOfType("*git.PushOptions")
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, nil)
	mockGitClient.On("CreateIndex", &git.Repository{}).Return(&git.Index{}, nil)
	mockGitClient.On("IndexAddAll", &git.Index{}).Return(nil)
	mockGitClient.On("WriteToIndexTree", &git.Index{}, &git.Repository{}).Return(&git.Oid{}, nil)
	mockGitClient.On("WriteIndex", &git.Index{}).Return(nil)
	mockGitClient.On("LookUpTree", &git.Repository{}, &git.Oid{}).Return(&git.Tree{}, nil)
	mockGitClient.On("CreateCommit", &git.Repository{}, "HEAD", sig, sig, "Gophr versioned repo authorName/repoName@repoSHA", &git.Tree{}).Return(nil)
	mockGitClient.On("CreateRef", &git.Repository{}, "HEAD", "refs/heads/master", true, "headOne").Return(nil)
	mockGitClient.On("CheckoutHead", &git.Repository{}, checkoutOpts).Return(nil)
	mockGitClient.On("CreateRemote", &git.Repository{}, "origin", remoteURL).Return(&git.Remote{}, nil)
	mockGitClient.On("Push", &git.Remote{}, refspec, pushOpts).Return(errors.New("This is an error"))
	args = packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
		creds: &config.Credentials{
			GithubPush: config.UserPass{
				User: "test",
				Pass: "testpassword",
			},
		},
	}
	err = pushToDepot(args)
	assert.NotNil(t, err)

	mockGitClient = g.NewMockClient()
	sig = mock.AnythingOfType("*git.Signature")
	checkoutOpts = mock.AnythingOfType("*git.CheckoutOpts")
	remoteURL = mock.AnythingOfType("string")
	refspec = mock.AnythingOfType("[]string")
	pushOpts = mock.AnythingOfType("*git.PushOptions")
	mockGitClient.On("InitRepo", "/archive/dir/path", false).Return(&git.Repository{}, nil)
	mockGitClient.On("CreateIndex", &git.Repository{}).Return(&git.Index{}, nil)
	mockGitClient.On("IndexAddAll", &git.Index{}).Return(nil)
	mockGitClient.On("WriteToIndexTree", &git.Index{}, &git.Repository{}).Return(&git.Oid{}, nil)
	mockGitClient.On("WriteIndex", &git.Index{}).Return(nil)
	mockGitClient.On("LookUpTree", &git.Repository{}, &git.Oid{}).Return(&git.Tree{}, nil)
	mockGitClient.On("CreateCommit", &git.Repository{}, "HEAD", sig, sig, "Gophr versioned repo authorName/repoName@repoSHA", &git.Tree{}).Return(nil)
	mockGitClient.On("CreateRef", &git.Repository{}, "HEAD", "refs/heads/master", true, "headOne").Return(nil)
	mockGitClient.On("CheckoutHead", &git.Repository{}, checkoutOpts).Return(nil)
	mockGitClient.On("CreateRemote", &git.Repository{}, "origin", remoteURL).Return(&git.Remote{}, nil)
	mockGitClient.On("Push", &git.Remote{}, refspec, pushOpts).Return(nil)
	args = packagePusherArgs{
		author: "authorName",
		repo:   "repoName",
		sha:    "repoSHA",
		packagePaths: packageDownloadPaths{
			archiveDirPath: "/archive/dir/path",
		},
		gitClient: mockGitClient,
		creds: &config.Credentials{
			GithubPush: config.UserPass{
				User: "test",
				Pass: "testpassword",
			},
		},
	}
	err = pushToDepot(args)
	assert.Nil(t, err)
}

func TestGenerateCredentialsCallback(t *testing.T) {
	fn := generateCredentialsCallback("test", "name")
	_, cred := fn("test", "name", 7)
	assert.True(t, cred.HasUsername(), "Checking if credential is valid")
}

func TestCertificateCheckCallback(t *testing.T) {
	errorCode := certificateCheckCallback(nil, false, "")
	assert.EqualValues(t, 0, errorCode, "Certificate callback should return 0")
}
