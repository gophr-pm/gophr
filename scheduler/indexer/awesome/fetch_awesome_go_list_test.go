package awesome

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFetchAwesomeGoList(t *testing.T) {
	Convey("Given a url", t, func() {
		var (
			badGodocMarkdown  = []byte("*[chalk](https://github.com/ttacon/chalk) - Intuitive package for prettifying terminal/console output.")
			goodGodocMarkdown = []byte("### Contents\n*[chalk](https://github.com/ttacon/chalk) - Intuitive package for prettifying terminal/console output.")
		)

		Convey("Should return correct PackageTuple struct", func() {
			actualOutput, err := FetchAwesomeGoList(
				FetchAwesomeGoListArgs{
					doHTTPGet: func(url string) ([]byte, error) {
						return goodGodocMarkdown, nil
					},
				},
			)

			So(err, ShouldBeNil)
			So(actualOutput, ShouldResemble, []PackageTuple{
				PackageTuple{author: "ttacon", repo: "chalk"},
			})
		})

		Convey("Should error, since the markdown does not contain '### Content'", func() {
			actualOutput, err := FetchAwesomeGoList(
				FetchAwesomeGoListArgs{
					doHTTPGet: func(url string) ([]byte, error) {
						return badGodocMarkdown, nil
					},
				},
			)

			So(err, ShouldNotBeNil)
			So(actualOutput, ShouldBeNil)
		})

		Convey("Should error, since the url failed to return", func() {
			actualOutput, err := FetchAwesomeGoList(
				FetchAwesomeGoListArgs{
					doHTTPGet: func(url string) ([]byte, error) {
						return nil, errors.New("failed")
					},
				},
			)

			So(err, ShouldNotBeNil)
			So(actualOutput, ShouldBeNil)
		})
	})
}
