package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSemverCandidate(t *testing.T) {
	var (
		err       error
		candidate SemverCandidate
	)

	candidate, err = NewSemverCandidate("", "b", "", "", "", "", "")
	assert.NotNil(t, err, "should fail for invalid ref hash")

	candidate, err = NewSemverCandidate("a", "", "", "", "", "", "")
	assert.NotNil(t, err, "should fail for invalid ref name")

	candidate, err = NewSemverCandidate("a", "b", "1asdlkj", "", "", "", "")
	assert.NotNil(t, err, "should fail for invalid major")

	candidate, err = NewSemverCandidate("a", "b", "", "2", "", "", "")
	assert.NotNil(t, err, "should fail for invalid major")

	candidate, err = NewSemverCandidate("a", "b", "1", "uds", "", "", "")
	assert.NotNil(t, err, "should fail for invalid minor")

	candidate, err = NewSemverCandidate("a", "b", "1", "2", "asdjhak", "", "")
	assert.NotNil(t, err, "should fail for invalid patch")

	candidate, err = NewSemverCandidate("a", "b", "1", "2", "3", "whatever", "kjdahkdjsh")
	assert.NotNil(t, err, "should fail for invalid prerelease")

	candidate, err = NewSemverCandidate("a", "b", "1", "", "", "", "")
	assert.Nil(t, err, "should not return an error in simple cases")
	assert.Equal(t, "a", candidate.GitRefHash, "hash should match")
	assert.Equal(t, "b", candidate.GitRefName, "name should match")
	assert.Equal(t, 1, candidate.MajorVersion, "major should match")
	assert.Equal(t, 0, candidate.MinorVersion, "minor should match")
	assert.Equal(t, 0, candidate.PatchVersion, "patch should match")
	assert.Equal(t, "", candidate.PrereleaseLabel, "prerelease label should match")
	assert.Equal(t, 0, candidate.PrereleaseVersion, "prerelease version should match")
	assert.Equal(t, false, candidate.PrereleaseVersionExists, "prerelease version should not exist")

	candidate, err = NewSemverCandidate("a", "b", "1", "2", "3", "alpha", "4")
	assert.Nil(t, err, "should not return an error in simple cases")
	assert.Equal(t, "a", candidate.GitRefHash, "hash should match")
	assert.Equal(t, "b", candidate.GitRefName, "name should match")
	assert.Equal(t, 1, candidate.MajorVersion, "major should match")
	assert.Equal(t, 2, candidate.MinorVersion, "minor should match")
	assert.Equal(t, 3, candidate.PatchVersion, "patch should match")
	assert.Equal(t, "alpha", candidate.PrereleaseLabel, "prerelease label should match")
	assert.Equal(t, 4, candidate.PrereleaseVersion, "prerelease version should match")
	assert.Equal(t, true, candidate.PrereleaseVersionExists, "prerelease version should exist")
}
