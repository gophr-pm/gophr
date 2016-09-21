package depot

import "fmt"

const (
	// DepotPublicServiceAddress is the address for which all public requests will
	// be resolved.
	DepotPublicServiceAddress = "depot-svc"
	// DepotInternalServiceAddress is the address for which all internal requests
	// will be resolved.
	DepotInternalServiceAddress = "depot-svc:3000"
	// UploadPackURLTemplate is the address for all interal clone requests
	UploadPackURLTemplate = "https://%s/depot/%s"
	depotReposPath        = "/repos"
)

func BuildHashedFolderName(author string, repo string, sha string) string {
	repoName := BuildHashedRepoName(author, repo, sha)
	return repoName + ".git"
}

// BuildHashedRepoName creates a new repo name hash uses for repo creation
// and lookup. Eliminates collision between similar usernames and packages/versions
func BuildHashedRepoName(author string, repo string, sha string) string {
	return fmt.Sprintf(
		"%d%s%d%s-%s",
		len(author),
		author,
		len(repo),
		repo,
		sha[:6])
}
