package dtos

//go:generate ffjson $GOFILE

// NewGitHubRepoDTO used as a DTO for building POST requests to Github
// to create new repos
type NewGitHubRepoDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
}
