package depot

import (
	"fmt"
	"net/http"
	"time"
)

const (
	depotReposPath = "/repos"
)

var (
	httpClient = &http.Client{Timeout: 10 * time.Second}
)

// BuildHashedRepoName creates a new repo name hash uses for repo creation
// and lookup. Eliminates collision between similar usernames and packages/versions
func BuildHashedRepoName(author string, repo string, ref string) string {
	return fmt.Sprintf("%d%s%d%s-%s", len(author), author, len(repo), repo, ref[:6])
}
