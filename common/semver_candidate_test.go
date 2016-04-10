package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var inequalityTuples = []*semverInequalityTuple{
	// Absolute comparison greater-thans
	version("1.2.3").isGreaterThanVersion("1.2.2"),
	version("2.3.4").isGreaterThanVersion("1.3.4"),
	version("2.3.4").isGreaterThanVersion("2.2.4"),
	version("2.3.4").isGreaterThanVersion("2.3.4-alpha"),
	version("1.2.3").isGreaterThanVersion("1.2.2"),
	version("2.3.4").isGreaterThanVersion("1.3.4"),
	version("2.3.4").isGreaterThanVersion("2.2.4"),
	version("2.3.4").isGreaterThanVersion("2.3.4-alpha"),
	version("1.2.3").isGreaterThanVersion("1.2.2"),
	version("2.3.4").isGreaterThanVersion("2.3.3"),
	version("2.3.4-alpha.2").isGreaterThanVersion("2.3.4-alpha.1"),
	// Absolute comparison less-thans
	version("1.2.3").isLessThanVersion("2.2.3"),
	version("1.2.3").isLessThanVersion("1.3.3"),
	version("1.2.3").isLessThanVersion("1.2.4"),
	version("1.2.3-alpha").isLessThanVersion("1.2.3"),
	version("1.2.3-alpha").isLessThanVersion("1.2.3-beta"),
	version("1.2.3-beta.1").isLessThanVersion("1.2.3-beta.3"),
	version("1.2.3").isLessThanVersion("2.2.3"),
	version("1.2.3").isLessThanVersion("1.3.3"),
	version("1.2.3").isLessThanVersion("1.2.4"),
	version("1.2.3-alpha").isLessThanVersion("1.2.3"),
	version("1.2.3-alpha").isLessThanVersion("1.2.3-beta"),
	version("1.2.3-beta.1").isLessThanVersion("1.2.3-beta.3"),
	version("1.2.2").isLessThanVersion("1.2.3"),
	version("2.3.3").isLessThanVersion("2.3.4"),
	version("2.3.4-alpha.1").isLessThanVersion("2.3.4-alpha.2"),
	// Absolute comparison equal-tos
	version("1.2.2").isEqualToVersion("1.2.2"),
	version("1.2.3-beta").isEqualToVersion("1.2.3-beta"),
	version("1.2.3-beta.1").isEqualToVersion("1.2.3-beta.1"),
}

type semverInequalityTuple struct {
	expected                        int
	subject, object                 string
	compiledSubject, compiledObject SemverCandidate
}

func version(subject string) *semverInequalityTuple {
	return &semverInequalityTuple{subject: subject}
}

func (tuple *semverInequalityTuple) isGreaterThanVersion(object string) *semverInequalityTuple {
	tuple.object = object
	tuple.expected = 1
	return tuple
}

func (tuple *semverInequalityTuple) isEqualToVersion(object string) *semverInequalityTuple {
	tuple.object = object
	tuple.expected = 0
	return tuple
}

func (tuple *semverInequalityTuple) isLessThanVersion(object string) *semverInequalityTuple {
	tuple.object = object
	tuple.expected = -1
	return tuple
}

func (tuple *semverInequalityTuple) compile() {
	matches := candidateRegex.FindStringSubmatch(tuple.subject)
	if matches != nil {
		compiledSubject, err := NewSemverCandidate(
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
			tuple.compiledSubject = compiledSubject
		} else {
			panic(fmt.Sprint("A test subject could not be initialized:", tuple.subject, "(", err, ")"))
		}
	} else {
		panic(fmt.Sprint("A test subject was invalid:", tuple.subject))
	}

	matches = candidateRegex.FindStringSubmatch(tuple.object)
	if matches != nil {
		compiledObject, err := NewSemverCandidate(
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
			tuple.compiledObject = compiledObject
		} else {
			panic(fmt.Sprint("A test object could not be initialized:", tuple.object, "(", err, ")"))
		}
	} else {
		panic(fmt.Sprint("A test object was invalid:", tuple.object))
	}
}

func TestNewSemverCandidate(t *testing.T) {
	var (
		err       error
		candidate SemverCandidate
	)

	candidate, err = NewSemverCandidate("", "b", "", "", "", "", "", "")
	assert.NotNil(t, err, "should fail for invalid ref hash")

	candidate, err = NewSemverCandidate("a", "", "", "", "", "", "", "")
	assert.NotNil(t, err, "should fail for invalid ref name")

	candidate, err = NewSemverCandidate("a", "b", "", "1asdlkj", "", "", "", "")
	assert.NotNil(t, err, "should fail for invalid major")

	candidate, err = NewSemverCandidate("a", "b", "", "", "2", "", "", "")
	assert.NotNil(t, err, "should fail for invalid major")

	candidate, err = NewSemverCandidate("a", "b", "", "1", "uds", "", "", "")
	assert.NotNil(t, err, "should fail for invalid minor")

	candidate, err = NewSemverCandidate("a", "b", "", "1", "2", "asdjhak", "", "")
	assert.NotNil(t, err, "should fail for invalid patch")

	candidate, err = NewSemverCandidate("a", "b", "", "1", "2", "3", "whatever", "kjdahkdjsh")
	assert.NotNil(t, err, "should fail for invalid prerelease")

	candidate, err = NewSemverCandidate("a", "b", "", "1", "", "", "", "")
	assert.Nil(t, err, "should not return an error in simple cases")
	assert.Equal(t, "a", candidate.GitRefHash, "hash should match")
	assert.Equal(t, "b", candidate.GitRefName, "name should match")
	assert.Equal(t, 1, candidate.MajorVersion, "major should match")
	assert.Equal(t, 0, candidate.MinorVersion, "minor should match")
	assert.Equal(t, 0, candidate.PatchVersion, "patch should match")
	assert.Equal(t, "", candidate.PrereleaseLabel, "prerelease label should match")
	assert.Equal(t, 0, candidate.PrereleaseVersion, "prerelease version should match")
	assert.Equal(t, false, candidate.PrereleaseVersionExists, "prerelease version should not exist")

	candidate, err = NewSemverCandidate("a", "b", "", "1", "2", "3", "alpha", "4")
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

func TestSemverCandidateCompareTo(t *testing.T) {
	// Compile all the things first
	for _, tuple := range inequalityTuples {
		tuple.compile()
	}

	// Run all the tests fam
	for _, tuple := range inequalityTuples {
		assert.Equal(
			t,
			tuple.expected,
			tuple.compiledSubject.CompareTo(tuple.compiledObject),
			fmt.Sprintf(
				`"%s" compared to "%s" should be %d`,
				tuple.subject,
				tuple.object,
				tuple.expected,
			),
		)
	}
}

func TestSemverCandidateList(t *testing.T) {
	var (
		list                               SemverCandidateList
		c, expectedLowest, expectedHighest SemverCandidate
	)

	c, _ = NewSemverCandidate("a", "b", "c", "1", "2", "3", "", "")
	list = append(list, c)
	c, _ = NewSemverCandidate("a", "b", "c", "2", "3", "4", "", "")
	list = append(list, c)
	expectedLowest, _ = NewSemverCandidate("a", "b", "c", "1", "2", "3", "alpha", "2")
	list = append(list, expectedLowest)
	c, _ = NewSemverCandidate("a", "b", "c", "2", "1", "1", "", "")
	list = append(list, c)
	expectedHighest, _ = NewSemverCandidate("a", "b", "c", "3", "5", "6", "", "")
	list = append(list, expectedHighest)

	lowest := list.Lowest()
	highest := list.Highest()

	assert.Equal(t, 0, lowest.CompareTo(expectedLowest), "The lowest candidate should be correct")
	assert.Equal(t, 0, highest.CompareTo(expectedHighest), "The highest candidate should be correct")

	selector, err := NewSemverSelector("", "1", "x", "", "", "", "")
	matchedList := list.Match(selector)
	assert.Nil(t, err, "Match should always work for normal candidate lists")
	assert.Equal(t, 2, len(matchedList), "Match should return the correct number of elements")

	list = []SemverCandidate{}
	assert.Nil(t, list.Lowest(), "Lowset should return nil when there are no elements to sort")
	assert.Nil(t, list.Highest(), "Highest should return nil when there are no elements to sort")
}
