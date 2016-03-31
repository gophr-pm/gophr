package main

// Repo represents a source code repository on GitHub.
type Repo struct {
	User         string
	Name         string
	SubPath      string
	OldFormat    bool // The old /v2/pkg format.
	MajorVersion Version

	// FullVersion is the best version in AllVersions that matches MajorVersion.
	// It defaults to InvalidVersion if there are no matches.
	FullVersion Version

	// AllVersions holds all versions currently available in the repository,
	// either coming from branch names or from tag names. Version zero (v0)
	// is only present in the list if it really exists in the repository.
	AllVersions VersionList
}

// SetVersions records in the relevant fields the details about which
// package versions are available in the repository.
func (repo *Repo) SetVersions(all []Version) {
	repo.AllVersions = all
	for _, v := range repo.AllVersions {
		if v.Major == repo.MajorVersion.Major && v.Unstable == repo.MajorVersion.Unstable && repo.FullVersion.Less(v) {
			repo.FullVersion = v
		}
	}
}

// GitHubRoot returns the repository root at GitHub, without a schema.
func (repo *Repo) GitHubRoot() string {
	if repo.User == "" {
		return "github.com/go-" + repo.Name + "/" + repo.Name
	} else {
		return "github.com/" + repo.User + "/" + repo.Name
	}
}

// GitHubTree returns the repository tree name at GitHub for the selected version.
func (repo *Repo) GitHubTree() string {
	if repo.FullVersion == InvalidVersion {
		return "master"
	}
	return repo.FullVersion.String()
}

// GopkgRoot returns the package root at gopkg.in, without a schema.
func (repo *Repo) GopkgRoot() string {
	return repo.GopkgVersionRoot(repo.MajorVersion)
}

// GopkgPath returns the package path at gopkg.in, without a schema.
func (repo *Repo) GopkgPath() string {
	return repo.GopkgVersionRoot(repo.MajorVersion) + repo.SubPath
}

// GopkgVersionRoot returns the package root in gopkg.in for the
// provided version, without a schema.
func (repo *Repo) GopkgVersionRoot(version Version) string {
	version.Minor = -1
	version.Patch = -1
	v := version.String()
	if repo.OldFormat {
		if repo.User == "" {
			return "gopkg.in/" + v + "/" + repo.Name
		}
		return "gopkg.in/" + repo.User + "/" + v + "/" + repo.Name
	} else {
		if repo.User == "" {
			return "gopkg.in/" + repo.Name + "." + v
		}
		return "gopkg.in/" + repo.User + "/" + repo.Name + "." + v
	}
}
