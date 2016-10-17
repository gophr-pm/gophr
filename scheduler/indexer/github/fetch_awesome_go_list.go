package github

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

const (
	awesomeGoURL                   = "https://raw.githubusercontent.com/avelino/awesome-go/master/README.md"
	awesomeGoReadmeContentsHeading = "### Contents"
)

var (
	linkRegex = regexp.MustCompile(`\[[^\]]+\]\(https://github\.com/([^\)]+)`)
)

type awesomePackage struct {
	author string
	repo   string
}

// fetchAwesomeGoList returns a map of all awesome go packages.
func fetchAwesomeGoList() ([]awesomePackage, error) {
	resp, err := http.Get(awesomeGoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	i := bytes.IndexAny(body, awesomeGoReadmeContentsHeading)
	if i == -1 {
		return nil, errors.New("Failed to find contents heading in awesome go readme.")
	}

	matches := linkRegex.FindAllStringSubmatch(string(body[i:]), -1)
	var awesomePackages []awesomePackage
	for _, match := range matches {
		packageParts := strings.Split(match[1], "/")
		if len(packageParts) >= 2 {
			awesomePackages = append(awesomePackages, awesomePackage{
				author: packageParts[0],
				repo:   packageParts[1],
			})
		}
	}

	return awesomePackages, nil
}
