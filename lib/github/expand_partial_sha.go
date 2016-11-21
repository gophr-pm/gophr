package github

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/datadog"
)

const (
	etagHeader              = "Etag"
	baseGithubArchiveURL    = "https://github.com/%s/%s/archive/%s.zip"
	minSHALengthWithQuotes  = 42
	ddEventExpandPartialSHA = "github.expand-partial-sha"
)

// HTTPHeadReq executes an HTTP `HEAD` to the specified URL and returns the
// corresponding response.
type HTTPHeadReq func(url string) (*http.Header, error)

// ExpandPartialSHAArgs is the arguments struct for ExpandPartialSHA.
type ExpandPartialSHAArgs struct {
	Author     string
	Repo       string
	ShortSHA   string
	DoHTTPHead HTTPHeadReq
}

// ExpandPartialSHA is responsible for fetching a full commit SHA from a short
// SHA. This works by sending a HEAD request to the git archive endpoint with a
// short SHA. The request returns a full SHA of the archive in the `Etag`
// of the request header that is sent back.
func (svc *requestServiceImpl) ExpandPartialSHA(
	args ExpandPartialSHAArgs,
) (string, error) {
	// Specify monitoring parameters.
	trackingArgs := datadog.TrackTransactionArgs{
		Tags:      []string{"github", datadog.TagInternal},
		Client:    svc.ddClient,
		AlertType: datadog.Success,
		StartTime: time.Now(),
		EventInfo: []string{fmt.Sprintf(
			`{ author: "%s", repo: "%s", sha: "%s" }`,
			args.Author,
			args.Repo,
			args.ShortSHA,
		)},
		MetricName:      datadog.MetricJobDuration,
		CreateEvent:     statsd.NewEvent,
		CustomEventName: ddEventExpandPartialSHA,
	}

	// Ensure that the transaction is tracked after the job finishes.
	defer datadog.TrackTransaction(&trackingArgs)

	log.Printf(`Expanding partial SHA "%s" of "%s/%s".
`, args.ShortSHA, args.Author, args.Repo)

	archiveURL := fmt.Sprintf(
		baseGithubArchiveURL,
		args.Author,
		args.Repo,
		args.ShortSHA)

	gitHubRespHeader, err := args.DoHTTPHead(archiveURL)
	if err != nil {
		// Make sure that the error is recorded in the datadog transaction.
		trackingArgs.AlertType = datadog.Error
		trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())

		return "", err
	}

	eTagHeader := gitHubRespHeader.Get(etagHeader)
	if len(eTagHeader) != minSHALengthWithQuotes {
		// Manually create an error for this kind of failure.
		err = errors.New(
			"Unable to retrieve full commit SHA, " +
				"Etag header was incomplete or empty.")

		// Make sure that the error is recorded in the datadog transaction.
		trackingArgs.AlertType = datadog.Error
		trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())

		return "", err
	}

	// TODO(skeswa): the quotes can be a messy assumption. Should probably just be
	// using strings.Trim(). @Shikkic.
	// The Etag in the header contains the full SHA wrapped in quotes.
	// We need to remove the quotes.
	fullSHA := eTagHeader[1 : len(eTagHeader)-1]

	return fullSHA, nil
}
