package dtos

import "time"

//go:generate ffjson $GOFILE

// GithubCommit is the response to a Github API commit detail request.
type GithubCommit struct {
	SHA string `json:"sha"`
}

// GithubCommitLookUp is the response to a Github API commit detail request.
type GithubCommitLookUp struct {
	Commit *GithubCommitDetail `json:"commit"`
}

// GithubCommitDetail is part of GithubCommitLookUp.
type GithubCommitDetail struct {
	Committer *GithubCommitCommitter `json:"committer"`
}

// GithubCommitCommitter is part of GithubCommitDetail.
type GithubCommitCommitter struct {
	Date time.Time `json:"date"`
}
