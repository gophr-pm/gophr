package github

import "fmt"

// ParseStarCount TODO Won't need this after implementing FFJSON
func ParseStarCount(responseBody map[string]interface{}) int {
	starCount := responseBody["stargazers_count"]
	if starCount == nil {
		return 0
	}

	return int(starCount.(float64))
}

// BuildGitHubBranch creates a new ref based on a hash of the old ref
// TODO Delete, but router depends on it, fix it
func BuildGitHubBranch(ref string) string {
	repoHash := ref[:len(ref)-1]
	return repoHash
}

// BuildNewGitHubRepoName creates a new repo name hash uses for repo creation
// and lookup. Eliminates collision between similar usernames and packages
// TODO Delete, but router depends on it, fix it
func BuildNewGitHubRepoName(author string, repo string) string {
	return fmt.Sprintf("%d%s%d%s", len(author), author, len(repo), repo)
}
