package datadog

import (
	"fmt"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/config"
)

const (
	datadogServicePort = "8125"
	datadogServiceName = "dd-agent-svc"
)

// NewClient returns a new statsd DataDog client for sending metrics to the
// DataDog agent.
func NewClient(conf *config.Config, nameSpace string) (*statsd.Client, error) {
	c, err := statsd.New(
		fmt.Sprintf("%s:%s", datadogServiceName, datadogServicePort),
	)
	if err != nil {
		return nil, err
	}

	c.Namespace = nameSpace

	return c, err
}
