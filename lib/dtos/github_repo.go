package dtos

//go:generate ffjson $GOFILE

// GithubRepoDTO is the DTO used to parse out relevant metrics from the repo
// detail endpoint of the Github API.
type GithubRepoDTO struct {
	Stars       string `json:"stargazers_count"`
	Description string `json:"description"`
}
