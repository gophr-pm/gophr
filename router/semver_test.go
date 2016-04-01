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

	semver, err = NewSemver("6", "", "", "", "", "", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("", "c", "", "", "", "", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("", "", "", "", "", "", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("", "1", "", "1", "", "", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("", "1", "z", "", "", "", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("", "1", "x", "1", "", "", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("", "1", "x", "x", "", "", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("~", "1", "x", "", "", "", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("~", "1", "1", "x", "", "", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("~", "1", "1", "z", "", "", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("~", "1", "1", "", "alpha", "", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("~", "1", "1", "1", "alpha", "x", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("~", "1", "1", "", "", "x", "")
	assert.NotNil(t, err)

	semver, err = NewSemver("~", "1", "1", "1", "alpha", "z", "")
	assert.NotNil(t, err)

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
}
