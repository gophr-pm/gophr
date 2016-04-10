package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoGetTemplateDataSource(t *testing.T) {
	selector, _ := NewSemverSelector("", "1", "x", "", "", "", "")

	getConfig().dev = true
	ds := NewGoGetTemplateDataSource("skeswa", "gophr", false, SemverSelector{}, "", false, SemverCandidate{})
	assert.Equal(t, "github.com/skeswa/gophr", ds.GithubRoot, "Github root should be correct")
	assert.Equal(t, "master", ds.GithubTree, "Github tree should be correct")
	assert.Equal(t, "gophr.dev/skeswa/gophr", ds.GophrPath, "Gophr patch should be correct")
	assert.Equal(t, "gophr.dev/skeswa/gophr", ds.GophrRoot, "Gophr root should be correct")
	assert.Equal(t, "http", ds.Protocol, "Protocol should be http for dev")

	getConfig().dev = false
	ds = NewGoGetTemplateDataSource("skeswa", "gophr", true, selector, "/test", false, SemverCandidate{})
	assert.Equal(t, "github.com/skeswa/gophr", ds.GithubRoot, "Github root should be correct")
	assert.Equal(t, "master", ds.GithubTree, "Github tree should be correct")
	assert.Equal(t, "gophr.dev/skeswa/gophr@1.x/test", ds.GophrPath, "Gophr path should be correct")
	assert.Equal(t, "gophr.dev/skeswa/gophr@1.x", ds.GophrRoot, "Gophr root should be correct")
	assert.Equal(t, "https", ds.Protocol, "Protocol should be https for not dev")

	getConfig().dev = true
	candidate, _ := NewSemverCandidate("a", "b", "somelabel", "1", "1", "1", "", "")
	ds = NewGoGetTemplateDataSource("skeswa", "gophr", true, selector, "", true, candidate)
	assert.Equal(t, "github.com/skeswa/gophr", ds.GithubRoot, "Github root should be correct")
	assert.Equal(t, "somelabel", ds.GithubTree, "Github tree should be correct")
	assert.Equal(t, "gophr.dev/skeswa/gophr@1.x", ds.GophrPath, "Gophr patch should be correct")
	assert.Equal(t, "gophr.dev/skeswa/gophr@1.x", ds.GophrRoot, "Gophr root should be correct")
}
