package main

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/skeswa/gophr/common/semver"
)

const (
	at                          = '@'
	dot                         = '.'
	hyphen                      = '-'
	slash                       = '/'
	shaLength                   = 40
	shortSHALength              = 6
	semverSelectorRegexTemplate = `([\%c\%c]?)([0-9]+)(?:\.([0-9]+|%c))?(?:\.([0-9]+|%c))?(?:\-([a-zA-Z0-9\-_]+[a-zA-Z0-9])(?:\.([0-9]+|%c))?)?([\%c\%c]?)`
)

var (
	// semverSelectorRegex is the regular expression used to parse semver package
	// version selectors.
	semverSelectorRegex = regexp.MustCompile(fmt.Sprintf(
		semverSelectorRegexTemplate,
		semver.SemverSelectorTildeChar,
		semver.SemverSelectorCaratChar,
		semver.SemverSelectorWildcardChar,
		semver.SemverSelectorWildcardChar,
		semver.SemverSelectorWildcardChar,
		semver.SemverSelectorLessThanChar,
		semver.SemverSelectorGreaterThanChar,
	))
)

// packageRequestParts represents the piecewise breakdown of a package request.
// TODO: Remove shaSelector and replace it with fullSHA and shortSHA
type packageRequestParts struct {
	url                   string
	repo                  string
	author                string
	subpath               string
	selector              string
	shaSelector           string
	semverSelector        semver.SemverSelector
	hasFullSHASelector    bool
	hasShortSHASelector   bool
	semverSelectorDefined bool
}

// hasSHASelector returns true if this parts struct has a sha selector.
func (parts *packageRequestParts) hasSHASelector() bool {
	return len(parts.shaSelector) > 0
}

// hasSemverSelector returns true if this parts struct has a semver selector.
func (parts *packageRequestParts) hasSemverSelector() bool {
	return parts.semverSelectorDefined
}

// String returns a string representation of this struct. This function returns
// a JSON-like representation of this struct.
func (parts *packageRequestParts) String() string {
	var b bytes.Buffer
	b.WriteString("{ ")
	b.WriteString("url: \"")
	if parts != nil {
		b.WriteString(parts.url)
	}
	b.WriteString("\", repo: \"")
	if parts != nil {
		b.WriteString(parts.repo)
	}
	b.WriteString("\", author: \"")
	if parts != nil {
		b.WriteString(parts.author)
	}
	b.WriteString("\", subpath: \"")
	if parts != nil {
		b.WriteString(parts.subpath)
	}
	b.WriteString("\", selector: \"")
	if parts != nil {
		b.WriteString(parts.selector)
	}
	b.WriteString("\", shaSelector: \"")
	if parts != nil {
		b.WriteString(parts.shaSelector)
	}
	b.WriteString("\", semverSelector: ")
	if parts != nil {
		b.WriteString(parts.semverSelector.String())
	}
	b.WriteString(" }")

	return b.String()
}

// readPackageRequestParts reads an http request in the format of a package
// request, and breaks down the URL of the request into parts. Lastly, the parts
// are composed into a parts struct and returned.
func readPackageRequestParts(req *http.Request) (*packageRequestParts, error) {
	var (
		i      = 0
		url    = strings.TrimSpace(req.URL.Path)
		urlLen = len(url)

		repoEndIndex       = -1 // Exclusive
		repoStartIndex     = -1 // Inclusive
		authorEndIndex     = -1
		authorStartIndex   = -1
		subpathStartIndex  = -1
		selectorStartIndex = -1

		selector              string
		shaSelector           string
		semverSelector        semver.SemverSelector
		semverSelectorDefined bool
	)

	// Exit if the the url is empty or just a slash or doesn't start with a slash.
	// If so, return an empty parts.
	if len(url) < 2 || url[0] != slash {
		return nil, NewInvalidPackageVersionRequestURLError(url)
	}

	// So, the first character should be a slash. That means the next character is
	// the beginning of the author.
	authorStartIndex = 1

	// Next step is to scan to the next slash to find the beginning of the repo.
	for i = authorStartIndex + 1; i < urlLen && url[i] != slash; i = i + 1 {
	}
	// Make sure the author exists (it is required).
	if (i - authorStartIndex) < 2 {
		return nil, NewInvalidPackageVersionRequestURLError(url)
	}
	// If we ran out bytes, then this is an invalid pac
	if i == urlLen {
		return nil, NewInvalidPackageVersionRequestURLError(url)
	}

	// So, we have arrived at the slash that prefixes the repo.
	authorEndIndex = i
	repoStartIndex = i + 1

	// Next step is to scan to the next slash OR at OR end of the string.
	for i = repoStartIndex; i < urlLen; i = i + 1 {
		char := url[i]
		// Time to read the subpath!
		if char == slash {
			repoEndIndex = i
			subpathStartIndex = i
			break
		}
		// Time to read the selector!
		if char == at {
			repoEndIndex = i
			selectorStartIndex = i + 1
			break
		}
	}
	// Make sure the repo exists (it is required).
	if (i - repoStartIndex) < 2 {
		return nil, NewInvalidPackageVersionRequestURLError(url)
	}
	// If we're out of url bytes, exit with just the author and the repo.
	if i == urlLen {
		return &packageRequestParts{
			url:    url,
			repo:   url[repoStartIndex:urlLen],
			author: url[authorStartIndex:authorEndIndex],
		}, nil
	}

	// Only read the selector if there is evidence that it exists.
	if selectorStartIndex != -1 {
		// If we got here, then we are reading until the end of the url or the
		// beginning of the subpath.
		for i = selectorStartIndex; i < urlLen && url[i] != slash; i = i + 1 {
		}
		// Make sure the selector exists (it is required if we see an @).
		if (i - selectorStartIndex) < 2 {
			return nil, NewInvalidPackageVersionRequestURLError(url)
		}
		// Whatever the case may be, this is where the selector ends.
		selector = url[selectorStartIndex:i]
		hasShortSHASelector := false
		hasFullSHASelector := false

		// Read the selector to figure out what it is.
		if isShortSHASelector(selector) {
			hasShortSHASelector = true
			shaSelector = selector
		} else if isFullSHASelector(selector) {
			hasFullSHASelector = true
			shaSelector = selector
		} else {
			var err error
			if semverSelector, err = readSemverSelector(selector); err != nil {
				return nil, NewInvalidPackageVersionRequestURLError(url, err)
			}

			// If we got here, the semver selector exists.
			semverSelectorDefined = true
		}

		// If we're out of url bytes then there is no subpath.
		if i == urlLen {
			return &packageRequestParts{
				url:                   url,
				repo:                  url[repoStartIndex:repoEndIndex],
				author:                url[authorStartIndex:authorEndIndex],
				selector:              selector,
				shaSelector:           shaSelector,
				semverSelector:        semverSelector,
				hasShortSHASelector:   hasShortSHASelector,
				hasFullSHASelector:    hasFullSHASelector,
				semverSelectorDefined: semverSelectorDefined,
			}, nil
		}

		// Otherwise, this is the subpath start index.
		subpathStartIndex = i
	}

	// Make sure the subpath is valid.
	if (urlLen - subpathStartIndex) < 2 {
		return nil, NewInvalidPackageVersionRequestURLError(url)
	}

	// If there are still unexplored bytes, there is a subpath as well.
	return &packageRequestParts{
		url:                   url,
		repo:                  url[repoStartIndex:repoEndIndex],
		author:                url[authorStartIndex:authorEndIndex],
		subpath:               url[subpathStartIndex:urlLen],
		selector:              selector,
		shaSelector:           shaSelector,
		semverSelector:        semverSelector,
		semverSelectorDefined: semverSelectorDefined,
	}, nil
}

// isFullSHASelector returns true if the selector is in a full sha
func isFullSHASelector(selector string) bool {
	// If it isn't a sha hash (which is 40 characters long), then its a semver
	// selector for out purposes.
	return len(selector) == shaLength && strings.IndexByte(selector, dot) == -1
}

func isShortSHASelector(selector string) bool {
	// If it isn't a 6 character short SHA return false, it might be semvar
	return len(selector) == shortSHALength && strings.IndexByte(selector, dot) == -1 && strings.IndexByte(selector, hyphen) == -1
}

// isSemverSelector converts a semver selector string into a semver selector.
func readSemverSelector(selector string) (semver.SemverSelector, error) {
	match := semverSelectorRegex.FindStringSubmatch(selector)
	if match == nil {
		return semver.SemverSelector{}, fmt.Errorf("Invalid version selector \"%s\"", selector)
	}

	semverSelector, err := semver.NewSemverSelector(
		match[1], // Prefix
		match[2], // Major Version
		match[3], // Minor Version
		match[4], // Patch Version
		match[5], // Pre-release Label
		match[6], // Pre-release Version
		match[7], // Suffix
	)
	if err != nil {
		return semver.SemverSelector{}, err
	}

	return semverSelector, nil
}
