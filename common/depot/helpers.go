package depot

import "fmt"

// DepotServiceAddress is the address for which all public requests will be resolved
const (
	DepotPublicServiceAddress = "depot-svc"
)

// DepotPrivateServiceAddress is the address for which all internal requests will be resolved
const (
	DepotInternalServiceAddress = "depot-svc:3000"
)

const (
	depotReposPath = "/repos"
)

// BuildHashedRepoName creates a new repo name hash uses for repo creation
// and lookup. Eliminates collision between similar usernames and packages/versions
func BuildHashedRepoName(author string, repo string, ref string) string {
	return fmt.Sprintf("%d%s%d%s-%s", len(author), author, len(repo), repo, ref[:6])
}
