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

// Client is used to interface with the DataDog API.
type Client interface {
	Incr(name string, tags []string, rate float64) error
	Event(e *statsd.Event) error
	Gauge(name string, value float64, tags []string, rate float64) error
}

// NewClient returns a new statsd DataDog client for sending metrics to the
// DataDog agent.
func NewClient(conf *config.Config, nameSpace string) (Client, error) {
	if conf.IsDev {
		return NewFakeDataDogClient(), nil
	}

	c, err := statsd.New(
		fmt.Sprintf("%s:%s", datadogServiceName, datadogServicePort),
	)
	if err != nil {
		return nil, err
	}

	c.Namespace = nameSpace

	return c, err
}
