package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSemver(t *testing.T) {
	var (
		err    error
		semver Semver
	)

	semver, err = NewSemver("6", "1", "", "", "", "", "")
	assert.NotNil(t, err, "should fail on illegal prefixes")

	semver, err = NewSemver("", "1", "", "", "", "", "?")
	assert.NotNil(t, err, "should fail on illegal suffixes")

	semver, err = NewSemver("", "c", "", "", "", "", "")
	assert.NotNil(t, err, "should fail on illegal major segment")

	semver, err = NewSemver("", "", "", "", "", "", "")
	assert.NotNil(t, err, "should fail on no major segment provided")

	semver, err = NewSemver("", "1", "", "1", "", "", "")
	assert.NotNil(t, err, "should fail on gap between version segments")

	semver, err = NewSemver("", "1", "z", "", "", "", "")
	assert.NotNil(t, err, "should fail on illegal minor segment")

	semver, err = NewSemver("", "1", "x", "1", "", "", "")
	assert.NotNil(t, err, "should fail on an segment trailing a wildcard")

	semver, err = NewSemver("", "1", "x", "x", "", "", "")
	assert.NotNil(t, err, "should fail on an segment trailing a wildcard")

	semver, err = NewSemver("~", "1", "x", "", "", "", "")
	assert.NotNil(t, err, "should fail when prefix is mixed with minor wildcard")

	semver, err = NewSemver("~", "1", "1", "x", "", "", "")
	assert.NotNil(t, err, "should fail when prefix is mixed with patch wildcard")

	semver, err = NewSemver("", "1", "1", "x", "alpha", "", "")
	assert.NotNil(t, err, "should fail on an segment trailing a wildcard")

	semver, err = NewSemver("~", "1", "1", "z", "", "", "")
	assert.NotNil(t, err, "should fail on illegal patch segment")

	semver, err = NewSemver("~", "1", "1", "", "alpha", "", "")
	assert.NotNil(t, err, "should fail on gap between version segments")

	semver, err = NewSemver("~", "1", "1", "1", "alpha", "x", "")
	assert.NotNil(t, err, "should fail when prefix is mixed with prerelease wildcard")

	semver, err = NewSemver("~", "1", "1", "", "", "x", "")
	assert.NotNil(t, err, "should fail on gap between version segments")

	semver, err = NewSemver("~", "1", "1", "1", "alpha", "z", "")
	assert.NotNil(t, err, "should fail on illegal prelease segment")

	semver, err = NewSemver("~", "1", "2", "", "", "", "+")
	assert.NotNil(t, err, "should fail when prefix is mixed with suffix")

	semver, err = NewSemver("", "1", "2", "x", "", "", "+")
	assert.NotNil(t, err, "should fail when wildcard is mixed with suffix")

	semver, err = NewSemver("", "1", "2", "x", "", "", "x")
	assert.NotNil(t, err, "should fail when wildcard is mixed with suffix")

	// semver, err = NewSemver("~", "1", "", "", "", "", "")
	// assert.NotNil(t, err)

	semver, err = NewSemver("", "1", "", "", "", "", "")
	assert.Nil(t, err)
	assert.Equal(t, semverPrefixNone, semver.Prefix, "prefix should be unspecified")
	assert.Equal(t, semverSegmentTypeNumber, semver.MajorVersion.Type, "major should be type number")
	assert.Equal(t, 1, semver.MajorVersion.Number, "major should be the correct number")
	assert.Equal(t, semverSegmentTypeUnspecified, semver.MinorVersion.Type, "minor should be type number")
	assert.Equal(t, semverSegmentTypeUnspecified, semver.PatchVersion.Type, "patch should be type unspecified")
	assert.Equal(t, "", semver.PrereleaseLabel, "prerelease label should be empty")
	assert.Equal(t, semverSegmentTypeUnspecified, semver.PrereleaseVersion.Type, "prerelease should be type unspecified")
	assert.Equal(t, semverSuffixNone, semver.Suffix, "suffix should be unspecified")

	semver, err = NewSemver("", "2", "", "", "", "", "-")
	assert.Nil(t, err)
	assert.Equal(t, semverPrefixNone, semver.Prefix, "prefix should be unspecified")
	assert.Equal(t, semverSegmentTypeNumber, semver.MajorVersion.Type, "major should be type number")
	assert.Equal(t, 2, semver.MajorVersion.Number, "major should be the correct number")
	assert.Equal(t, semverSegmentTypeUnspecified, semver.MinorVersion.Type, "minor should be type number")
	assert.Equal(t, semverSegmentTypeUnspecified, semver.PatchVersion.Type, "patch should be type unspecified")
	assert.Equal(t, "", semver.PrereleaseLabel, "prerelease label should be empty")
	assert.Equal(t, semverSegmentTypeUnspecified, semver.PrereleaseVersion.Type, "prerelease should be type unspecified")
	assert.Equal(t, semverSuffixLessThan, semver.Suffix, "suffix should be less than")

	semver, err = NewSemver("~", "1", "2", "", "", "", "")
	assert.Nil(t, err)
	assert.Equal(t, semverPrefixTilde, semver.Prefix, "prefix should be a tilde")
	assert.Equal(t, semverSegmentTypeNumber, semver.MajorVersion.Type, "major should be type number")
	assert.Equal(t, 1, semver.MajorVersion.Number, "major should be the correct number")
	assert.Equal(t, semverSegmentTypeNumber, semver.MinorVersion.Type, "minor should be type number")
	assert.Equal(t, 2, semver.MinorVersion.Number, "minor should be the correct number")
	assert.Equal(t, semverSegmentTypeUnspecified, semver.PatchVersion.Type, "patch should be type unspecified")
	assert.Equal(t, "", semver.PrereleaseLabel, "prerelease label should be empty")
	assert.Equal(t, semverSegmentTypeUnspecified, semver.PrereleaseVersion.Type, "prerelease should be type unspecified")
	assert.Equal(t, semverSuffixNone, semver.Suffix, "suffix should be unspecified")

	semver, err = NewSemver("^", "1", "2", "3", "", "", "")
	assert.Nil(t, err)
	assert.Equal(t, semverPrefixCarat, semver.Prefix, "prefix should be a carat")
	assert.Equal(t, semverSegmentTypeNumber, semver.MajorVersion.Type, "major should be type number")
	assert.Equal(t, 1, semver.MajorVersion.Number, "major should be the correct number")
	assert.Equal(t, semverSegmentTypeNumber, semver.MinorVersion.Type, "minor should be type number")
	assert.Equal(t, 2, semver.MinorVersion.Number, "minor should be the correct number")
	assert.Equal(t, semverSegmentTypeNumber, semver.PatchVersion.Type, "patch should be type number")
	assert.Equal(t, 3, semver.PatchVersion.Number, "patch should be the correct number")
	assert.Equal(t, "", semver.PrereleaseLabel, "prerelease label should be empty")
	assert.Equal(t, semverSegmentTypeUnspecified, semver.PrereleaseVersion.Type, "prerelease should be type unspecified")
	assert.Equal(t, semverSuffixNone, semver.Suffix, "suffix should be unspecified")

	semver, err = NewSemver("", "1", "2", "3", "alpha", "x", "")
	assert.Nil(t, err)
	assert.Equal(t, semverPrefixNone, semver.Prefix, "prefix should be unspecified")
	assert.Equal(t, semverSegmentTypeNumber, semver.MajorVersion.Type, "major should be type number")
	assert.Equal(t, 1, semver.MajorVersion.Number, "major should be the correct number")
	assert.Equal(t, semverSegmentTypeNumber, semver.MinorVersion.Type, "minor should be type number")
	assert.Equal(t, 2, semver.MinorVersion.Number, "minor should be the correct number")
	assert.Equal(t, semverSegmentTypeNumber, semver.PatchVersion.Type, "patch should be type number")
	assert.Equal(t, 3, semver.PatchVersion.Number, "patch should be the correct number")
	assert.Equal(t, "alpha", semver.PrereleaseLabel, "prerelease label should be alpha")
	assert.Equal(t, semverSegmentTypeWildcard, semver.PrereleaseVersion.Type, "prerelease should be type wildcard")
	assert.Equal(t, semverSuffixNone, semver.Suffix, "suffix should be unspecified")

	semver, err = NewSemver("", "1", "2", "3", "beta", "43", "+")
	assert.Nil(t, err)
	assert.Equal(t, semverPrefixNone, semver.Prefix, "prefix should be unspecified")
	assert.Equal(t, semverSegmentTypeNumber, semver.MajorVersion.Type, "major should be type number")
	assert.Equal(t, 1, semver.MajorVersion.Number, "major should be the correct number")
	assert.Equal(t, semverSegmentTypeNumber, semver.MinorVersion.Type, "minor should be type number")
	assert.Equal(t, 2, semver.MinorVersion.Number, "minor should be the correct number")
	assert.Equal(t, semverSegmentTypeNumber, semver.PatchVersion.Type, "patch should be type number")
	assert.Equal(t, 3, semver.PatchVersion.Number, "patch should be the correct number")
	assert.Equal(t, "beta", semver.PrereleaseLabel, "prerelease label should be alpha")
	assert.Equal(t, semverSegmentTypeNumber, semver.PrereleaseVersion.Type, "prerelease should be type number")
	assert.Equal(t, 43, semver.PrereleaseVersion.Number, "prerelease should the correct number")
	assert.Equal(t, semverSuffixGreaterThan, semver.Suffix, "suffix should be greater than")
}
