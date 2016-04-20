package common

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type refsTest struct {
	summary           string
	original          string
	version           string
	changed           string
	versionCandidates []SemverCandidate
}

var refsTests = []refsTest{{
	"Version v0 works even without any references",
	reflines(
		"hash1 HEAD",
	),
	"v0",
	reflines(
		"hash1 HEAD",
	),
	nil,
}, {
	"Preserve original capabilities",
	reflines(
		"hash1 HEAD\x00caps",
	),
	"v0",
	reflines(
		"hash1 HEAD\x00caps",
	),
	nil,
}, {
	"Matching major version branch",
	reflines(
		"00000000000000000000000000000000000hash1 HEAD",
		"00000000000000000000000000000000000hash2 refs/heads/v0",
		"00000000000000000000000000000000000hash3 refs/heads/v1",
		"00000000000000000000000000000000000hash4 refs/heads/v2",
	),
	"v1",
	reflines(
		"00000000000000000000000000000000000hash4 HEAD\x00symref=HEAD:refs/heads/v2",
		"00000000000000000000000000000000000hash4 refs/heads/master",
		"00000000000000000000000000000000000hash2 refs/heads/v0",
		"00000000000000000000000000000000000hash3 refs/heads/v1",
		"00000000000000000000000000000000000hash4 refs/heads/v2",
	),
	[]SemverCandidate{
		{
			"00000000000000000000000000000000000hash2", // hash
			"refs/heads/v0",                            // name
			"v0",                                       // label
			0,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		}, {
			"00000000000000000000000000000000000hash3", // hash
			"refs/heads/v1",                            // name
			"v1",                                       // label
			1,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		}, {
			"00000000000000000000000000000000000hash4", // hash
			"refs/heads/v2",                            // name
			"v2",                                       // label
			2,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		},
	},
}, {
	"Matching minor version branch",
	reflines(
		"00000000000000000000000000000000000hash1 HEAD",
		"00000000000000000000000000000000000hash2 refs/heads/v1.1",
		"00000000000000000000000000000000000hash3 refs/heads/v1.3",
		"00000000000000000000000000000000000hash4 refs/heads/v1.2",
	),
	"v1",
	reflines(
		"00000000000000000000000000000000000hash4 HEAD\x00symref=HEAD:refs/heads/v1.2",
		"00000000000000000000000000000000000hash4 refs/heads/master",
		"00000000000000000000000000000000000hash2 refs/heads/v1.1",
		"00000000000000000000000000000000000hash3 refs/heads/v1.3",
		"00000000000000000000000000000000000hash4 refs/heads/v1.2",
	),
	[]SemverCandidate{
		{
			"00000000000000000000000000000000000hash2", // hash
			"refs/heads/v1.1",                          // name
			"v1.1",                                     // label
			1,                                          // major
			1,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		}, {
			"00000000000000000000000000000000000hash3", // hash
			"refs/heads/v1.3",                          // name
			"v1.3",                                     // label
			1,                                          // major
			3,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		}, {
			"00000000000000000000000000000000000hash4", // hash
			"refs/heads/v1.2",                          // name
			"v1.2",                                     // label
			1,                                          // major
			2,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		},
	},
}, {
	"Disable original symref capability",
	reflines(
		"00000000000000000000000000000000000hash1 HEAD\x00foo symref=bar baz",
		"00000000000000000000000000000000000hash2 refs/heads/v1",
	),
	"v1",
	reflines(
		"00000000000000000000000000000000000hash2 HEAD\x00symref=HEAD:refs/heads/v1 foo oldref=bar baz",
		"00000000000000000000000000000000000hash2 refs/heads/master",
		"00000000000000000000000000000000000hash2 refs/heads/v1",
	),
	[]SemverCandidate{
		{
			"00000000000000000000000000000000000hash2", // hash
			"refs/heads/v1",                            // name
			"v1",                                       // label
			1,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		},
	},
}, {
	"Disable original symref capability with tags",
	reflines(
		"00000000000000000000000000000000000hash1 HEAD\x00foo symref=bar baz",
		"00000000000000000000000000000000000hash2 refs/tags/v1",
	),
	"v1",
	reflines(
		"00000000000000000000000000000000000hash2 HEAD\x00foo oldref=bar baz",
		"00000000000000000000000000000000000hash2 refs/heads/master",
		"00000000000000000000000000000000000hash2 refs/tags/v1",
	),
	[]SemverCandidate{
		{
			"00000000000000000000000000000000000hash2", // hash
			"refs/tags/v1",                             // name
			"v1",                                       // label
			1,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		},
	},
}, {
	"Replace original master branch",
	reflines(
		"00000000000000000000000000000000000hash1 HEAD",
		"00000000000000000000000000000000000hash1 refs/heads/master",
		"00000000000000000000000000000000000hash2 refs/heads/v1",
	),
	"v1",
	reflines(
		"00000000000000000000000000000000000hash2 HEAD\x00symref=HEAD:refs/heads/v1",
		"00000000000000000000000000000000000hash2 refs/heads/master",
		"00000000000000000000000000000000000hash2 refs/heads/v1",
	),
	[]SemverCandidate{
		{
			"00000000000000000000000000000000000hash2", // hash
			"refs/heads/v1",                            // name
			"v1",                                       // label
			1,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		},
	},
}, {
	"Matching tag",
	reflines(
		"00000000000000000000000000000000000hash1 HEAD",
		"00000000000000000000000000000000000hash2 refs/tags/v0",
		"00000000000000000000000000000000000hash3 refs/tags/v1",
		"00000000000000000000000000000000000hash4 refs/tags/v2",
	),
	"v1",
	reflines(
		"00000000000000000000000000000000000hash4 HEAD",
		"00000000000000000000000000000000000hash4 refs/heads/master",
		"00000000000000000000000000000000000hash2 refs/tags/v0",
		"00000000000000000000000000000000000hash3 refs/tags/v1",
		"00000000000000000000000000000000000hash4 refs/tags/v2",
	),
	[]SemverCandidate{
		{
			"00000000000000000000000000000000000hash2", // hash
			"refs/tags/v0",                             // name
			"v0",                                       // label
			0,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		}, {
			"00000000000000000000000000000000000hash3", // hash
			"refs/tags/v1",                             // name
			"v1",                                       // label
			1,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		}, {
			"00000000000000000000000000000000000hash4", // hash
			"refs/tags/v2",                             // name
			"v2",                                       // label
			2,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		},
	},
}, {
	"Tag peeling",
	reflines(
		"00000000000000000000000000000000000hash1 HEAD",
		"00000000000000000000000000000000000hash2 refs/heads/master",
		"00000000000000000000000000000000000hash3 refs/tags/v1",
		"00000000000000000000000000000000000hash4 refs/tags/v1^{}",
		"00000000000000000000000000000000000hash5 refs/tags/v2",
	),
	"v1",
	reflines(
		"00000000000000000000000000000000000hash5 HEAD",
		"00000000000000000000000000000000000hash5 refs/heads/master",
		"00000000000000000000000000000000000hash3 refs/tags/v1",
		"00000000000000000000000000000000000hash4 refs/tags/v1^{}",
		"00000000000000000000000000000000000hash5 refs/tags/v2",
	),
	[]SemverCandidate{
		{
			"00000000000000000000000000000000000hash3", // hash
			"refs/tags/v1",                             // name
			"v1",                                       // label
			1,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		}, {
			"00000000000000000000000000000000000hash4", // hash
			"refs/tags/v1",                             // name
			"v1",                                       // label
			1,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		}, {
			"00000000000000000000000000000000000hash5", // hash
			"refs/tags/v2",                             // name
			"v2",                                       // label
			2,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		},
	},
}, {
	"Matching unstable versions",
	reflines(
		"00000000000000000000000000000000000hash1 HEAD",
		"00000000000000000000000000000000000hash2 refs/heads/master",
		"00000000000000000000000000000000000hash3 refs/heads/v1",
		"00000000000000000000000000000000000hash4 refs/heads/v1.1-unstable",
		"00000000000000000000000000000000000hash5 refs/heads/v1.3-unstable",
		"00000000000000000000000000000000000hash6 refs/heads/v1.2-unstable",
		"00000000000000000000000000000000000hash7 refs/heads/v2",
	),
	"v1-unstable",
	reflines(
		"00000000000000000000000000000000000hash7 HEAD\x00symref=HEAD:refs/heads/v2",
		"00000000000000000000000000000000000hash7 refs/heads/master",
		"00000000000000000000000000000000000hash3 refs/heads/v1",
		"00000000000000000000000000000000000hash4 refs/heads/v1.1-unstable",
		"00000000000000000000000000000000000hash5 refs/heads/v1.3-unstable",
		"00000000000000000000000000000000000hash6 refs/heads/v1.2-unstable",
		"00000000000000000000000000000000000hash7 refs/heads/v2",
	),
	[]SemverCandidate{
		{
			"00000000000000000000000000000000000hash3", // hash
			"refs/heads/v1",                            // name
			"v1",                                       // label
			1,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		}, {
			"00000000000000000000000000000000000hash4", // hash
			"refs/heads/v1.1-unstable",                 // name
			"v1.1-unstable",                            // label
			1,                                          // major
			1,                                          // minor
			0,                                          // patch
			"unstable",                                 // pre-release label
			0,                                          // pre-release version
			true,                                       // pre-release exists
		}, {
			"00000000000000000000000000000000000hash5", // hash
			"refs/heads/v1.3-unstable",                 // name
			"v1.3-unstable",                            // label
			1,                                          // major
			3,                                          // minor
			0,                                          // patch
			"unstable",                                 // pre-release label
			0,                                          // pre-release version
			true,                                       // pre-release exists
		}, {
			"00000000000000000000000000000000000hash6", // hash
			"refs/heads/v1.2-unstable",                 // name
			"v1.2-unstable",                            // label
			1,                                          // major
			2,                                          // minor
			0,                                          // patch
			"unstable",                                 // pre-release label
			0,                                          // pre-release version
			true,                                       // pre-release exists
		}, {
			"00000000000000000000000000000000000hash7", // hash
			"refs/heads/v2",                            // name
			"v2",                                       // label
			2,                                          // major
			0,                                          // minor
			0,                                          // patch
			"",                                         // pre-release label
			0,                                          // pre-release version
			false,                                      // pre-release exists
		},
	},
}}

func reflines(lines ...string) string {
	var buf bytes.Buffer
	buf.WriteString("001e# service=git-upload-pack\n0000")
	for _, l := range lines {
		buf.WriteString(fmt.Sprintf("%04x%s\n", len(l)+5, l))
	}
	buf.WriteString("0000")
	return buf.String()
}

func invalidSizeStringReflines() string {
	var buf bytes.Buffer
	buf.WriteString("001z# service=git-upload-pack\n0000")
	return buf.String()
}

func sizeTooBigReflines() string {
	var buf bytes.Buffer
	buf.WriteString("9999# service=git-upload-pack\n0000")
	return buf.String()
}

func candidatesMatch(t *testing.T, candidates1 []SemverCandidate, candidates2 []SemverCandidate) bool {
	if candidates1 == nil && candidates2 != nil {
		t.Logf("Candidates 1 was nil and Candidates 2 was not")
		return false
	} else if candidates1 != nil && candidates2 == nil {
		t.Logf("Candidates 2 was nil and Candidates 1 was not")
		return false
	} else if len(candidates1) != len(candidates2) {
		t.Logf("Candidates 1 (length %d) was a different length than Candidates 2 (length %d)", len(candidates1), len(candidates2))
		return false
	}

	for i, c1 := range candidates1 {
		c2 := candidates2[i]
		if !(c1.GitRefHash == c1.GitRefHash &&
			c1.GitRefName == c2.GitRefName &&
			c1.MajorVersion == c2.MajorVersion &&
			c1.MinorVersion == c2.MinorVersion &&
			c1.PatchVersion == c2.PatchVersion &&
			c1.PrereleaseLabel == c2.PrereleaseLabel &&
			c1.PrereleaseVersion == c2.PrereleaseVersion &&
			c1.PrereleaseVersionExists == c2.PrereleaseVersionExists) {
			t.Logf("Candidate 1 (%v) was different than Candidate 2 (%v)", candidates1, candidates2)
			return false
		}
	}

	return true
}

func TestGetRefs(t *testing.T) {
	_, err := FetchRefs("thisisnotarealthing", "notevenalittle")
	assert.NotNil(t, err, "fetch should fail since the github root is invalid")

	_, err = FetchRefs("skeswa", "gophr")
	assert.Nil(t, err, "fetch should work for valid repos")
}

func TestUseRefs(t *testing.T) {
	refs, err := NewRefs([]byte(invalidSizeStringReflines()))
	assert.NotNil(t, err, "refs parsing should have failed because the size wasn't a number")

	refs, err = NewRefs([]byte(sizeTooBigReflines()))
	assert.NotNil(t, err, "refs parsing should have failed because the size too big")

	for _, test := range refsTests {
		t.Log(test.summary)

		refs, err = NewRefs([]byte(test.original))
		assert.Nil(t, err, "refs should have been parsed correctly")
		assert.True(t, candidatesMatch(t, test.versionCandidates, refs.Candidates), "parsed version candidates should be correct")

		if len(refs.Candidates) > 0 {
			lastCandidate := refs.Candidates[len(refs.Candidates)-1]
			reserializedRefs := refs.Reserialize(lastCandidate.GitRefName, lastCandidate.GitRefHash)
			assert.Equal(t, test.changed, string(reserializedRefs[:]), "refs should have been serialized correctly")
		}
	}
}
