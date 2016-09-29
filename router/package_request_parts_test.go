package main

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/gophr-pm/gophr/common/semver"
	"github.com/stretchr/testify/assert"
)

func TestHasSHASelector(t *testing.T) {
	parts := packageRequestParts{
		shaSelector: "thisisashaselectorthisisashaselecto",
	}
	assert.True(t, parts.hasSHASelector())

	parts = packageRequestParts{
		shaSelector: "",
	}
	assert.False(t, parts.hasSHASelector())
}

func TestHasSemverSelector(t *testing.T) {
	req := &http.Request{URL: &url.URL{Path: "/abc/def"}}
	parts, _ := readPackageRequestParts(req)
	assert.False(t, parts.hasSemverSelector())

	req = &http.Request{URL: &url.URL{Path: "/abc/def@2.3.4"}}
	parts, _ = readPackageRequestParts(req)
	assert.True(t, parts.hasSemverSelector())

	req = &http.Request{URL: &url.URL{Path: "/abc/def@1234567890123456789012345678901234567890"}}
	parts, _ = readPackageRequestParts(req)
	assert.False(t, parts.hasSemverSelector())
}

func TestReadPackageRequestParts(t *testing.T) {
	req := &http.Request{URL: &url.URL{Path: "/abc/def"}}
	expectedParts := &packageRequestParts{
		url:                   "/abc/def",
		repo:                  "def",
		author:                "abc",
		subpath:               "",
		selector:              "",
		shaSelector:           "",
		semverSelector:        semver.SemverSelector{},
		semverSelectorDefined: false,
	}
	actualParts, err := readPackageRequestParts(req)
	assert.Nil(t, err)
	assert.True(
		t,
		reflect.DeepEqual(expectedParts, actualParts),
		fmt.Sprintf("%s should equal %s", actualParts.String(), expectedParts.String()))

	req = &http.Request{URL: &url.URL{Path: "/abc/def@1.2.5+"}}
	semsel, _ := semver.NewSemverSelector("", "1", "2", "5", "", "", "+")
	expectedParts = &packageRequestParts{
		url:                   "/abc/def@1.2.5+",
		repo:                  "def",
		author:                "abc",
		subpath:               "",
		selector:              "1.2.5+",
		shaSelector:           "",
		semverSelector:        semsel,
		semverSelectorDefined: true,
	}
	actualParts, err = readPackageRequestParts(req)
	assert.Nil(t, err)
	assert.True(
		t,
		reflect.DeepEqual(expectedParts, actualParts),
		fmt.Sprintf("%s should equal %s", actualParts.String(), expectedParts.String()))

	req = &http.Request{URL: &url.URL{Path: "/abc/def@123456abcd123456abcd123456abcd123456abcd"}}
	expectedParts = &packageRequestParts{
		url:                   "/abc/def@123456abcd123456abcd123456abcd123456abcd",
		repo:                  "def",
		author:                "abc",
		subpath:               "",
		selector:              "123456abcd123456abcd123456abcd123456abcd",
		shaSelector:           "123456abcd123456abcd123456abcd123456abcd",
		semverSelector:        semver.SemverSelector{},
		semverSelectorDefined: false,
	}
	actualParts, err = readPackageRequestParts(req)
	assert.Nil(t, err)
	assert.True(
		t,
		reflect.DeepEqual(expectedParts, actualParts),
		fmt.Sprintf("%s should equal %s", actualParts.String(), expectedParts.String()))

	req = &http.Request{URL: &url.URL{Path: "/abc/def/ghi"}}
	expectedParts = &packageRequestParts{
		url:                   "/abc/def/ghi",
		repo:                  "def",
		author:                "abc",
		subpath:               "/ghi",
		selector:              "",
		shaSelector:           "",
		semverSelector:        semver.SemverSelector{},
		semverSelectorDefined: false,
	}
	actualParts, err = readPackageRequestParts(req)
	assert.Nil(t, err)
	assert.True(
		t,
		reflect.DeepEqual(expectedParts, actualParts),
		fmt.Sprintf("%s should equal %s", actualParts.String(), expectedParts.String()))

	req = &http.Request{URL: &url.URL{Path: "/abc/def@1.2.5+/ghi"}}
	semsel, _ = semver.NewSemverSelector("", "1", "2", "5", "", "", "+")
	expectedParts = &packageRequestParts{
		url:                   "/abc/def@1.2.5+/ghi",
		repo:                  "def",
		author:                "abc",
		subpath:               "/ghi",
		selector:              "1.2.5+",
		shaSelector:           "",
		semverSelector:        semsel,
		semverSelectorDefined: true,
	}
	actualParts, err = readPackageRequestParts(req)
	assert.Nil(t, err)
	assert.True(
		t,
		reflect.DeepEqual(expectedParts, actualParts),
		fmt.Sprintf("%s should equal %s", actualParts.String(), expectedParts.String()))

	req = &http.Request{URL: &url.URL{Path: "/abc/def@123456abcd123456abcd123456abcd123456abcd/ghi"}}
	expectedParts = &packageRequestParts{
		url:                   "/abc/def@123456abcd123456abcd123456abcd123456abcd/ghi",
		repo:                  "def",
		author:                "abc",
		subpath:               "/ghi",
		selector:              "123456abcd123456abcd123456abcd123456abcd",
		shaSelector:           "123456abcd123456abcd123456abcd123456abcd",
		semverSelector:        semver.SemverSelector{},
		semverSelectorDefined: false,
	}
	actualParts, err = readPackageRequestParts(req)
	assert.Nil(t, err)
	assert.True(
		t,
		reflect.DeepEqual(expectedParts, actualParts),
		fmt.Sprintf("%s should equal %s", actualParts.String(), expectedParts.String()))

	req = &http.Request{URL: &url.URL{Path: "/abc/def@sddm/ghi"}}
	_, err = readPackageRequestParts(req)
	assert.NotNil(t, err)

	req = &http.Request{URL: &url.URL{Path: "/abc/def@1.x.x+/ghi"}}
	_, err = readPackageRequestParts(req)
	assert.NotNil(t, err)

	req = &http.Request{URL: &url.URL{Path: "//"}}
	_, err = readPackageRequestParts(req)
	assert.NotNil(t, err)

	req = &http.Request{URL: &url.URL{Path: "/a//b"}}
	_, err = readPackageRequestParts(req)
	assert.NotNil(t, err)

	req = &http.Request{URL: &url.URL{Path: "/a/b@/c"}}
	_, err = readPackageRequestParts(req)
	assert.NotNil(t, err)

	req = &http.Request{URL: &url.URL{Path: "/a/b@1.x/"}}
	_, err = readPackageRequestParts(req)
	assert.NotNil(t, err)
}
