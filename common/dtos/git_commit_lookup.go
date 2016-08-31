package dtos

import "time"

//go:generate ffjson $GOFILE

type GitCommitLookUpDTO struct {
	Commit *GitCommitDetailDTO `json:"commit"`
}

type GitCommitDetailDTO struct {
	Committer *GitComitCommitterDTO `json:"committer"`
}

type GitComitCommitterDTO struct {
	Date time.Time `json:"date"`
}
