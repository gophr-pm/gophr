package main

import (
	"errors"
	"sync"
	"testing"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

func TestIsPackageArchived(t *testing.T) {
	db := &gocql.Session{}

	args := packageArchivalCheckerArgs{
		db:     db,
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		isPackageArchivedInDB: func(db *gocql.Session, author string, repo string, sha string) (bool, error) {
			return false, errors.New("this is an error of some kind")
		},
	}
	archived, err := isPackageArchived(args)
	assert.NotNil(t, err)
	assert.Equal(t, false, archived)

	args = packageArchivalCheckerArgs{
		db:     db,
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		isPackageArchivedInDB: func(db *gocql.Session, author string, repo string, sha string) (bool, error) {
			return true, nil
		},
	}
	archived, err = isPackageArchived(args)
	assert.Nil(t, err)
	assert.Equal(t, true, archived)

	args = packageArchivalCheckerArgs{
		db:     db,
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		isPackageArchivedInDB: func(_db *gocql.Session, author string, repo string, sha string) (bool, error) {
			assert.Equal(t, "myauthor", author)
			assert.Equal(t, "myrepo", repo)
			assert.Equal(t, "mysha", sha)
			assert.Equal(t, db, _db)
			return false, nil
		},
		packageExistsInDepot: func(author, repo, sha string) (bool, error) {
			assert.Equal(t, "myauthor", author)
			assert.Equal(t, "myrepo", repo)
			assert.Equal(t, "mysha", sha)
			return false, errors.New("this a big scary error")
		},
	}
	archived, err = isPackageArchived(args)
	assert.NotNil(t, err)
	assert.Equal(t, false, archived)

	args = packageArchivalCheckerArgs{
		db:     db,
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		isPackageArchivedInDB: func(_db *gocql.Session, author string, repo string, sha string) (bool, error) {
			assert.Equal(t, "myauthor", author)
			assert.Equal(t, "myrepo", repo)
			assert.Equal(t, "mysha", sha)
			assert.Equal(t, db, _db)
			return false, nil
		},
		packageExistsInDepot: func(author, repo, sha string) (bool, error) {
			assert.Equal(t, "myauthor", author)
			assert.Equal(t, "myrepo", repo)
			assert.Equal(t, "mysha", sha)
			return false, nil
		},
	}
	archived, err = isPackageArchived(args)
	assert.Nil(t, err)
	assert.Equal(t, false, archived)

	wg := sync.WaitGroup{}
	wg.Add(1)
	args = packageArchivalCheckerArgs{
		db:     db,
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		isPackageArchivedInDB: func(_db *gocql.Session, author string, repo string, sha string) (bool, error) {
			assert.Equal(t, "myauthor", author)
			assert.Equal(t, "myrepo", repo)
			assert.Equal(t, "mysha", sha)
			assert.Equal(t, db, _db)
			return false, nil
		},
		packageExistsInDepot: func(author, repo, sha string) (bool, error) {
			assert.Equal(t, "myauthor", author)
			assert.Equal(t, "myrepo", repo)
			assert.Equal(t, "mysha", sha)
			return true, nil
		},
		recordPackageArchival: func(args packageArchivalRecorderArgs) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			assert.Equal(t, db, args.db)
			wg.Done()
		},
	}
	archived, err = isPackageArchived(args)
	wg.Wait()
	assert.Nil(t, err)
	assert.Equal(t, true, archived)
}
