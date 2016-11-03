package awesome

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func TestFetchAwesomeGoList(t *testing.T) {
	Convey("Given a url", t, func() {
		var (
			badGodocMarkdown  = "*[chalk](https://github.com/ttacon/chalk) - Intuitive package for prettifying terminal/console output."
			goodGodocMarkdown = "### Contents\n*[chalk](https://github.com/ttacon/chalk) - Intuitive package for prettifying terminal/console output."
		)

		Convey("Should return correct packageTuple struct", func() {
			actualOutput, err := fetchAwesomeGoList(
				fetchAwesomeGoListArgs{
					doHTTPGet: func(url string) (*http.Response, error) {
						return &http.Response{
							Body: nopCloser{strings.NewReader(goodGodocMarkdown)},
						}, nil
					},
				},
			)

			So(err, ShouldBeNil)
			So(actualOutput, ShouldResemble, []packageTuple{
				packageTuple{author: "ttacon", repo: "chalk"},
			})
		})

		Convey("Should error, since the markdown does not contain '### Content'", func() {
			actualOutput, err := fetchAwesomeGoList(
				fetchAwesomeGoListArgs{
					doHTTPGet: func(url string) (*http.Response, error) {
						return &http.Response{
							Body: nopCloser{strings.NewReader(badGodocMarkdown)},
						}, nil
					},
				},
			)

			So(err, ShouldNotBeNil)
			So(actualOutput, ShouldBeNil)
		})

		Convey("Should error, since the url failed to return", func() {
			actualOutput, err := fetchAwesomeGoList(
				fetchAwesomeGoListArgs{
					doHTTPGet: func(url string) (*http.Response, error) {
						return nil, errors.New("failed")
					},
				},
			)

			So(err, ShouldNotBeNil)
			So(actualOutput, ShouldBeNil)
		})
	})
}
