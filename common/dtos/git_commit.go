package dtos

//go:generate ffjson $GOFILE

type GitCommitDTO struct {
	SHA string `json:"sha"`
}
