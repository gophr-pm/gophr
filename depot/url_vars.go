package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	urlVarAuthor = "author"
	urlVarRepo   = "repo"
	urlVarSHA    = "sha"
)

type urlVars struct {
	sha    string
	repo   string
	author string
}

// readURLVars reads author, repo & sha from the URL.
func readURLVars(r *http.Request) (urlVars, error) {
	var (
		vars = mux.Vars(r)

		sha    = vars[urlVarSHA]
		repo   = vars[urlVarRepo]
		author = vars[urlVarAuthor]
	)

	if len(sha) < 40 {
		return urlVars{}, fmt.Errorf(`Invalid value "%v" specified for URL variable "%s".`, urlVarSHA, sha)
	}
	if len(repo) < 1 {
		return urlVars{}, fmt.Errorf(`Invalid value "%v" specified for URL variable "%s".`, urlVarRepo, repo)
	}
	if len(author) < 1 {
		return urlVars{}, fmt.Errorf(`Invalid value "%v" specified for URL variable "%s".`, urlVarAuthor, author)
	}

	return urlVars{
		sha:    sha,
		repo:   repo,
		author: author,
	}, nil
}
