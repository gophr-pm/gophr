
# Gophr - go package management
[![Go Report Card](https://goreportcard.com/badge/github.com/skeswa/gophr)](https://goreportcard.com/report/github.com/skeswa/gophr)
[![GoDoc](https://godoc.org/github.com/skeswa/gophr/common?status.svg)](https://godoc.org/github.com/skeswa/gophr/common)
<a href="https://zenhub.com"><img src="https://raw.githubusercontent.com/ZenHubIO/support/master/zenhub-badge.png"></a>

An end-to-end package management solution for Go. **No manifest or lock file and a fully versioned dependency graph.** Simply place the url in your import path and it's automatically fully versioned.

`
gophr.pm/author/repo@(semver or SHA)
`

#### Native go dependency management
`go get` can only retrieve the current master branch. If you ever need to re-download your dependency it could be totally different each time.
```go
  import (
      // Un-versioned github link
      "github.com/a/b"
  )
```

#### Gophr dependency management and versioning
Gophr allows you to version lock your dependencies by semver or SHA.
```go
  import (
      // Version current master branch
      "gophr.pm/a/b"
      // Version by semver
      "gophr.pm/a/b@1.0"
      // Version by semver logic
      "gophr.pm/a/b@>1.0.0"
      "gophr.pm/a/b@<1.3.2"
      // Version by SHA
      "gophr.pm/a/b@24638c6d1aaa1a39c14c704918e354fd3949b93c"
  )
```

### The problem with native Go dependency management
Go has **no** ability to version a specific SHA or tag for a repo. Anytime you pull down an import it grabs the current master branch. This not only bad practice but it could potentially silently break your code without you ever knowing why.

### There are plenty of Go versioning tools. What makes Gophr special?

Gophr doesn't require you to install _any_ tooling to use. Simply place the versioned url `gophr.pm/author/repo@(semver or SHA)` in your import path and you're done.

We give you the power of semver to reference tags and create logical equivalence operations just like in `gem` and `npm`.

```go
"gophr.pm/a/b@>1.0.4"
"gophr.pm/a/b@<1.0.0"
```

We also **fully** version the entire dependency graph. Meaning, we version lock every sub-dependency as well, so it's perfectly preserved, everytime. Something **no one else** does.

#### Contributing
- [Setting up dev environment](example.com)
