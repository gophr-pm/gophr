package dtos

//go:generate ffjson $GOFILE

// GithubRepo is the  used to parse out relevant metrics from the repo
// detail endpoint of the Github API.
type GithubRepo struct {
	Stars       int    `json:"stargazers_count"`
	Description string `json:"description"`
}
