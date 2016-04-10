package common

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	selectorRegex = regexp.MustCompile(
		fmt.Sprintf(
			`([%c%c]?)([0-9]+)\.(?:([0-9]+|%c))(?:\.([0-9]+|%c))?(?:\-([a-zA-Z0-9\-_]+[a-zA-Z0-9])(?:\.([0-9]+|%c))?)?([%c%c]?)`,
			SemverSelectorTildeChar,
			SemverSelectorCaratChar,
			SemverSelectorWildcardChar,
			SemverSelectorWildcardChar,
			SemverSelectorWildcardChar,
			SemverSelectorLessThanChar,
			SemverSelectorGreaterThanChar,
		),
	)

	candidateRegex = regexp.MustCompile(
		fmt.Sprintf(
			`([0-9]+)\.(?:([0-9]+|%c))(?:\.([0-9]+|%c))?(?:\-([a-zA-Z0-9\-_]+)(?:\.([0-9]+|%c))?)?`,
			SemverSelectorWildcardChar,
			SemverSelectorWildcardChar,
			SemverSelectorWildcardChar,
		),
	)

	matchTuples = []*semverMatchTuple{
		// Tilde
		selector("~1.2.2").bounds("1.2.2"),
		selector("~1.3.2").bounds("1.3.3"),
		selector("~2.5.4").doesntBound("3.5.3"),
		selector("~2.5.4").doesntBound("2.6.3"),
		selector("~2.5.4").doesntBound("2.5.3"),
		selector("~2.5.9").doesntBound("2.5.9-beta"),
		// Carat
		selector("^1.2.2").bounds("1.2.2"),
		selector("^1.2.2").bounds("1.2.9"),
		selector("^1.2.2").bounds("1.9.11"),
		selector("^8.2.2").bounds("8.3.11"),
		selector("^8.2.2").bounds("8.2.5"),
		selector("^8.1.2").doesntBound("9.0.0"),
		selector("^8.2.2").doesntBound("8.1.11"),
		selector("^8.2.2").doesntBound("8.2.1"),
		selector("^8.1.2").doesntBound("8.1.2-alpha"),
		// Greater
		selector("1.2.2+").bounds("10.10.10"),
		selector("1.2.2+").bounds("1.3.5"),
		selector("1.2.2+").bounds("1.3.5"),
		selector("1.2.2+").bounds("1.2.3"),
		selector("1.2.2-alpha+").bounds("1.2.2"),
		selector("1.2.2-beta+").bounds("1.2.2-beta.5"),
		selector("1.2.2-beta.1+").bounds("1.2.2-beta.5"),
		selector("1.2.2-beta+").doesntBound("1.2.2-alpha"),
		selector("1.2.2-beta.2+").doesntBound("1.2.2-alpha.4"),
		selector("1.2.2+").doesntBound("1.2.2-beta"),
		selector("1.2.2+").doesntBound("1.2.1"),
		selector("1.2.2+").doesntBound("0.3.5"),
		selector("1.2.2+").doesntBound("1.1.2"),
		// Lesser
		selector("1.2.2-").bounds("1.2.1"),
		selector("1.2.2-").bounds("0.0.9"),
		selector("1.2.2-").bounds("0.3.5"),
		selector("1.2.2-").bounds("1.1.5"),
		selector("1.2.2-alpha-").bounds("1.2.1"),
		selector("1.2.2-thing.8-").bounds("1.2.2-thing.2"),
		selector("1.2.2-alpha-").doesntBound("1.2.2"),
		selector("1.2.2-beta-").doesntBound("1.2.2-alpha"),
		selector("1.2.2-thing-").doesntBound("1.2.2-thing.6"),
		selector("1.2.2-beta.1-").doesntBound("1.2.2-beta.5"),
		selector("1.2.2-beta.2-").doesntBound("1.2.2-alpha.4"),
		selector("1.5.5-beta-").doesntBound("1.5.5-beta.9"),
		selector("1.2.2-").doesntBound("1.2.3"),
		selector("1.2.2-").doesntBound("3.2.2"),
		selector("1.2.2-").doesntBound("1.3.2"),
		selector("1.2.2-").doesntBound("1.2.2-beta"),
		selector("1.2.2-").doesntBound("1.10.2"),
		// Wildcard
		selector("1.x").bounds("1.11.9"),
		selector("1.2.x").bounds("1.2.1"),
		selector("1.2.3-alpha.x").bounds("1.2.3-alpha.5"),
		selector("1.x").doesntBound("2.11.9"),
		// Vanilla
		selector("1.1.1").bounds("1.1.1"),
		selector("1.1.1-beta").bounds("1.1.1-beta"),
		selector("1.1.1-beta.2").bounds("1.1.1-beta.2"),
		selector("1.1.1-beta").doesntBound("1.1.1-alpha"),
	}
)

type semverMatchTuple struct {
	correct           bool
	selector          string
	compiledSelector  SemverSelector
	candidate         string
	compiledCandidate SemverCandidate
}

func selector(selector string) *semverMatchTuple {
	return &semverMatchTuple{selector: selector}
}

func (tuple *semverMatchTuple) bounds(candidate string) *semverMatchTuple {
	tuple.correct = true
	tuple.candidate = candidate
	return tuple
}

func (tuple *semverMatchTuple) doesntBound(candidate string) *semverMatchTuple {
	tuple.correct = false
	tuple.candidate = candidate
	return tuple
}

func (tuple *semverMatchTuple) compile() {
	matches := selectorRegex.FindStringSubmatch(tuple.selector)
	if matches != nil {
		compiledSelector, err := NewSemverSelector(
			matches[1],
			matches[2],
			matches[3],
			matches[4],
			matches[5],
			matches[6],
			matches[7],
		)

		if err == nil {
			tuple.compiledSelector = compiledSelector
		} else {
			panic(fmt.Sprint("A test selector could not be initialized:", tuple.selector, "(", err, ")"))
		}
	} else {
		panic(fmt.Sprint("A test selector was invalid:", tuple.selector))
	}

	matches = candidateRegex.FindStringSubmatch(tuple.candidate)
	if matches != nil {
		compiledCandidate, err := NewSemverCandidate(
			"fakeHash",
			"fakeName",
			"fakeLabel",
			matches[1],
			matches[2],
			matches[3],
			matches[4],
			matches[5],
		)

		if err == nil {
			tuple.compiledCandidate = compiledCandidate
		} else {
			panic(fmt.Sprint("A test candidate could not be initialized:", tuple.candidate, "(", err, ")"))
		}
	} else {
		panic(fmt.Sprint("A test candidate was invalid:", tuple.candidate))
	}
}

func TestNewSemverSelector(t *testing.T) {
	var (
		err    error
		semver SemverSelector
	)

	semver, err = NewSemverSelector("6", "1", "", "", "", "", "")
	assert.NotNil(t, err, "should fail on illegal prefixes")

	semver, err = NewSemverSelector("", "1", "", "", "", "", "?")
	assert.NotNil(t, err, "should fail on illegal suffixes")

	semver, err = NewSemverSelector("", "c", "", "", "", "", "")
	assert.NotNil(t, err, "should fail on illegal major segment")

	semver, err = NewSemverSelector("", "", "", "", "", "", "")
	assert.NotNil(t, err, "should fail on no major segment provided")

	semver, err = NewSemverSelector("", "1", "", "1", "", "", "")
	assert.NotNil(t, err, "should fail on gap between version segments")

	semver, err = NewSemverSelector("", "1", "z", "", "", "", "")
	assert.NotNil(t, err, "should fail on illegal minor segment")

	semver, err = NewSemverSelector("", "1", "x", "1", "", "", "")
	assert.NotNil(t, err, "should fail on an segment trailing a wildcard")

	semver, err = NewSemverSelector("", "1", "x", "x", "", "", "")
	assert.NotNil(t, err, "should fail on an segment trailing a wildcard")

	semver, err = NewSemverSelector("~", "1", "x", "", "", "", "")
	assert.NotNil(t, err, "should fail when prefix is mixed with minor wildcard")

	semver, err = NewSemverSelector("~", "1", "1", "x", "", "", "")
	assert.NotNil(t, err, "should fail when prefix is mixed with patch wildcard")

	semver, err = NewSemverSelector("", "1", "1", "x", "alpha", "", "")
	assert.NotNil(t, err, "should fail on an segment trailing a wildcard")

	semver, err = NewSemverSelector("~", "1", "1", "z", "", "", "")
	assert.NotNil(t, err, "should fail on illegal patch segment")

	semver, err = NewSemverSelector("~", "1", "1", "", "alpha", "", "")
	assert.NotNil(t, err, "should fail on gap between version segments")

	semver, err = NewSemverSelector("~", "1", "1", "1", "alpha", "x", "")
	assert.NotNil(t, err, "should fail when prefix is mixed with prerelease wildcard")

	semver, err = NewSemverSelector("~", "1", "1", "", "", "x", "")
	assert.NotNil(t, err, "should fail on gap between version segments")

	semver, err = NewSemverSelector("~", "1", "1", "1", "alpha", "z", "")
	assert.NotNil(t, err, "should fail on illegal prelease segment")

	semver, err = NewSemverSelector("~", "1", "2", "", "", "", "+")
	assert.NotNil(t, err, "should fail when prefix is mixed with suffix")

	semver, err = NewSemverSelector("", "1", "2", "x", "", "", "+")
	assert.NotNil(t, err, "should fail when wildcard is mixed with suffix")

	semver, err = NewSemverSelector("", "1", "2", "x", "", "", "x")
	assert.NotNil(t, err, "should fail when wildcard is mixed with suffix")

	// semver, err = NewSemverSelector("~", "1", "", "", "", "", "")
	// assert.NotNil(t, err)

	semver, err = NewSemverSelector("", "1", "", "", "", "", "")
	assert.Nil(t, err)
	assert.Equal(t, SemverSelectorPrefixNone, semver.Prefix, "prefix should be unspecified")
	assert.Equal(t, SemverSegmentTypeNumber, semver.MajorVersion.Type, "major should be type number")
	assert.Equal(t, 1, semver.MajorVersion.Number, "major should be the correct number")
	assert.Equal(t, SemverSegmentTypeUnspecified, semver.MinorVersion.Type, "minor should be type number")
	assert.Equal(t, SemverSegmentTypeUnspecified, semver.PatchVersion.Type, "patch should be type unspecified")
	assert.Equal(t, "", semver.PrereleaseLabel, "prerelease label should be empty")
	assert.Equal(t, SemverSegmentTypeUnspecified, semver.PrereleaseVersion.Type, "prerelease should be type unspecified")
	assert.Equal(t, SemverSelectorSuffixNone, semver.Suffix, "suffix should be unspecified")

	semver, err = NewSemverSelector("", "2", "", "", "", "", "-")
	assert.Nil(t, err)
	assert.Equal(t, SemverSelectorPrefixNone, semver.Prefix, "prefix should be unspecified")
	assert.Equal(t, SemverSegmentTypeNumber, semver.MajorVersion.Type, "major should be type number")
	assert.Equal(t, 2, semver.MajorVersion.Number, "major should be the correct number")
	assert.Equal(t, SemverSegmentTypeUnspecified, semver.MinorVersion.Type, "minor should be type number")
	assert.Equal(t, SemverSegmentTypeUnspecified, semver.PatchVersion.Type, "patch should be type unspecified")
	assert.Equal(t, "", semver.PrereleaseLabel, "prerelease label should be empty")
	assert.Equal(t, SemverSegmentTypeUnspecified, semver.PrereleaseVersion.Type, "prerelease should be type unspecified")
	assert.Equal(t, SemverSelectorSuffixLessThan, semver.Suffix, "suffix should be less than")

	semver, err = NewSemverSelector("~", "1", "2", "", "", "", "")
	assert.Nil(t, err)
	assert.Equal(t, SemverSelectorPrefixTilde, semver.Prefix, "prefix should be a tilde")
	assert.Equal(t, SemverSegmentTypeNumber, semver.MajorVersion.Type, "major should be type number")
	assert.Equal(t, 1, semver.MajorVersion.Number, "major should be the correct number")
	assert.Equal(t, SemverSegmentTypeNumber, semver.MinorVersion.Type, "minor should be type number")
	assert.Equal(t, 2, semver.MinorVersion.Number, "minor should be the correct number")
	assert.Equal(t, SemverSegmentTypeUnspecified, semver.PatchVersion.Type, "patch should be type unspecified")
	assert.Equal(t, "", semver.PrereleaseLabel, "prerelease label should be empty")
	assert.Equal(t, SemverSegmentTypeUnspecified, semver.PrereleaseVersion.Type, "prerelease should be type unspecified")
	assert.Equal(t, SemverSelectorSuffixNone, semver.Suffix, "suffix should be unspecified")

	semver, err = NewSemverSelector("^", "1", "2", "3", "", "", "")
	assert.Nil(t, err)
	assert.Equal(t, SemverSelectorPrefixCarat, semver.Prefix, "prefix should be a carat")
	assert.Equal(t, SemverSegmentTypeNumber, semver.MajorVersion.Type, "major should be type number")
	assert.Equal(t, 1, semver.MajorVersion.Number, "major should be the correct number")
	assert.Equal(t, SemverSegmentTypeNumber, semver.MinorVersion.Type, "minor should be type number")
	assert.Equal(t, 2, semver.MinorVersion.Number, "minor should be the correct number")
	assert.Equal(t, SemverSegmentTypeNumber, semver.PatchVersion.Type, "patch should be type number")
	assert.Equal(t, 3, semver.PatchVersion.Number, "patch should be the correct number")
	assert.Equal(t, "", semver.PrereleaseLabel, "prerelease label should be empty")
	assert.Equal(t, SemverSegmentTypeUnspecified, semver.PrereleaseVersion.Type, "prerelease should be type unspecified")
	assert.Equal(t, SemverSelectorSuffixNone, semver.Suffix, "suffix should be unspecified")

	semver, err = NewSemverSelector("", "1", "2", "3", "alpha", "x", "")
	assert.Nil(t, err)
	assert.Equal(t, SemverSelectorPrefixNone, semver.Prefix, "prefix should be unspecified")
	assert.Equal(t, SemverSegmentTypeNumber, semver.MajorVersion.Type, "major should be type number")
	assert.Equal(t, 1, semver.MajorVersion.Number, "major should be the correct number")
	assert.Equal(t, SemverSegmentTypeNumber, semver.MinorVersion.Type, "minor should be type number")
	assert.Equal(t, 2, semver.MinorVersion.Number, "minor should be the correct number")
	assert.Equal(t, SemverSegmentTypeNumber, semver.PatchVersion.Type, "patch should be type number")
	assert.Equal(t, 3, semver.PatchVersion.Number, "patch should be the correct number")
	assert.Equal(t, "alpha", semver.PrereleaseLabel, "prerelease label should be alpha")
	assert.Equal(t, SemverSegmentTypeWildcard, semver.PrereleaseVersion.Type, "prerelease should be type wildcard")
	assert.Equal(t, SemverSelectorSuffixNone, semver.Suffix, "suffix should be unspecified")

	semver, err = NewSemverSelector("", "1", "2", "3", "beta", "43", "+")
	assert.Nil(t, err)
	assert.Equal(t, SemverSelectorPrefixNone, semver.Prefix, "prefix should be unspecified")
	assert.Equal(t, SemverSegmentTypeNumber, semver.MajorVersion.Type, "major should be type number")
	assert.Equal(t, 1, semver.MajorVersion.Number, "major should be the correct number")
	assert.Equal(t, SemverSegmentTypeNumber, semver.MinorVersion.Type, "minor should be type number")
	assert.Equal(t, 2, semver.MinorVersion.Number, "minor should be the correct number")
	assert.Equal(t, SemverSegmentTypeNumber, semver.PatchVersion.Type, "patch should be type number")
	assert.Equal(t, 3, semver.PatchVersion.Number, "patch should be the correct number")
	assert.Equal(t, "beta", semver.PrereleaseLabel, "prerelease label should be alpha")
	assert.Equal(t, SemverSegmentTypeNumber, semver.PrereleaseVersion.Type, "prerelease should be type number")
	assert.Equal(t, 43, semver.PrereleaseVersion.Number, "prerelease should the correct number")
	assert.Equal(t, SemverSelectorSuffixGreaterThan, semver.Suffix, "suffix should be greater than")
}

func TestSemverMatches(t *testing.T) {
	// Compile all the things first
	for _, tuple := range matchTuples {
		tuple.compile()
	}

	// Run all the tests fam
	for _, tuple := range matchTuples {
		if tuple.correct {
			assert.True(
				t,
				tuple.compiledSelector.Matches(tuple.compiledCandidate),
				fmt.Sprintf(`"%s" should match "%s"`, tuple.selector, tuple.candidate),
			)
		} else {
			assert.False(
				t,
				tuple.compiledSelector.Matches(tuple.compiledCandidate),
				fmt.Sprintf(
					`"%s" shouldn't match "%s"`,
					tuple.selector,
					tuple.candidate,
				),
			)
		}
	}
}

func TestSemverString(t *testing.T) {
	var (
		semver SemverSelector
	)

	semver, _ = NewSemverSelector("", "1", "", "", "", "", "")
	assert.Equal(t, "1", semver.String(), "serialized semver should match expectations")

	semver, _ = NewSemverSelector("~", "1", "", "", "", "", "")
	assert.Equal(t, "~1", semver.String(), "serialized semver should match expectations")

	semver, _ = NewSemverSelector("^", "1", "", "", "", "", "")
	assert.Equal(t, "^1", semver.String(), "serialized semver should match expectations")

	semver, _ = NewSemverSelector("~", "1", "2", "", "", "", "")
	assert.Equal(t, "~1.2", semver.String(), "serialized semver should match expectations")

	semver, _ = NewSemverSelector("", "1", "2", "3", "", "", "")
	assert.Equal(t, "1.2.3", semver.String(), "serialized semver should match expectations")

	semver, _ = NewSemverSelector("", "1", "x", "", "", "", "")
	assert.Equal(t, "1.x", semver.String(), "serialized semver should match expectations")

	semver, _ = NewSemverSelector("", "1", "2", "x", "", "", "")
	assert.Equal(t, "1.2.x", semver.String(), "serialized semver should match expectations")

	semver, _ = NewSemverSelector("", "1", "2", "3", "alpha", "", "")
	assert.Equal(t, "1.2.3-alpha", semver.String(), "serialized semver should match expectations")

	semver, _ = NewSemverSelector("", "1", "2", "3", "alpha", "x", "")
	assert.Equal(t, "1.2.3-alpha.x", semver.String(), "serialized semver should match expectations")

	semver, _ = NewSemverSelector("", "1", "2", "3", "alpha", "4", "")
	assert.Equal(t, "1.2.3-alpha.4", semver.String(), "serialized semver should match expectations")

	semver, _ = NewSemverSelector("", "1", "2", "3", "alpha", "4", "+")
	assert.Equal(t, "1.2.3-alpha.4+", semver.String(), "serialized semver should match expectations")

	semver, _ = NewSemverSelector("", "1", "2", "3", "alpha", "4", "-")
	assert.Equal(t, "1.2.3-alpha.4-", semver.String(), "serialized semver should match expectations")

	semver = SemverSelector{
		MajorVersion: SemverSelectorSegment{
			Type: SemverSegmentTypeWildcard,
		},
	}
	assert.Panics(t, func() {
		_ = semver.String()
	}, "semver serialization should fail since major segment in an incorrect type")

	semver = SemverSelector{
		MajorVersion: SemverSelectorSegment{
			Type: SemverSegmentTypeUnspecified,
		},
	}
	assert.Panics(t, func() {
		_ = semver.String()
	}, "semver serialization should fail since major segment in an incorrect type")
}
