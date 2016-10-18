package awesome

import (
	"bytes"
	"errors"
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

// FetchAwesomeGoList parses all the packages from the README on github.com/avelino/awesome-go.
func FetchAwesomeGoList(args FetchAwesomeGoListArgs) ([]PackageTuple, error) {
	body, err := args.doHTTPGet(awesomeGoURL)
	if err != nil {
		return nil, errors.New("Failed to make request to awesome-go.")
	}

	i := bytes.Index(body, []byte(awesomeGoReadmeContentsHeading))
	if i == -1 {
		return nil, errors.New("Failed to find contents heading in awesome go readme.")
	}

	matches := linkRegex.FindAllStringSubmatch(string(body[i:]), -1)
	var PackageTuples []PackageTuple
	for _, match := range matches {
		packageParts := strings.Split(match[1], "/")
		if len(packageParts) >= 2 {
			PackageTuples = append(PackageTuples, PackageTuple{
				author: packageParts[0],
				repo:   packageParts[1],
			})
		}
	}

	return PackageTuples, nil
}
