package awesome

import (
	"bytes"
	"errors"
	"fmt"
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

// httpGetter executes an HTTP get to the specified URL and returns the
// corresponding response.
type httpGetter func(url string) (*http.Response, error)

// fetchAwesomeGoListArgs is the args struct for fetching awesome go packages
// from godoc.
type fetchAwesomeGoListArgs struct {
	doHTTPGet httpGetter
}

// fetchAwesomeGoList parses all the packages from the README on
// github.com/avelino/awesome-go.
func fetchAwesomeGoList(args fetchAwesomeGoListArgs) ([]packageTuple, error) {
	resp, err := args.doHTTPGet(awesomeGoURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to awesome-go: %v.", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read awesome-go response: %v.", err)
	}

	i := bytes.Index(body, []byte(awesomeGoReadmeContentsHeading))
	if i == -1 {
		return nil, errors.New(
			"Failed to find contents heading in awesome go readme.")
	}

	matches := linkRegex.FindAllStringSubmatch(string(body[i:]), -1)
	var packageTuples []packageTuple
	for _, match := range matches {
		packageParts := strings.Split(match[1], "/")
		if len(packageParts) >= 2 {
			packageTuples = append(packageTuples, packageTuple{
				author: packageParts[0],
				repo:   packageParts[1],
			})
		}
	}

	return packageTuples, nil
}
