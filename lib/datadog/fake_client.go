package datadog

import (
	"log"

	"github.com/DataDog/datadog-go/statsd"
)

// FakeClient is a mock for FakeClient.
type FakeClient struct {
	Tags      []string
	Namespace string
}

// NewFakeDataDogClient creates a new NewFakeDataDogClient.
func NewFakeDataDogClient() Client {
	return &FakeClient{}
}

// Gauge mocks FakeClient.Gauge.
func (fc *FakeClient) Gauge(
	name string,
	value float64,
	tags []string,
	rate float64,
) error {
	log.Printf(
		"Sending Gauge Metric to DataDog "+
			"(name: %s, value: %b, tags: %s, rate: %b).\n",
		name,
		value,
		tags,
		rate)

	return nil
}

// Event mocks FakeClient.Event.
func (fc *FakeClient) Event(e *statsd.Event) error {
	log.Printf(
		"Sending Custom Event to DataDog (event: %+v).\n",
		e)

	return nil
}

// Incr mocks FakeClient.Incr.
func (fc *FakeClient) Incr(name string, tags []string, rate float64) error {
	log.Printf(
		"Sending Increment Metric to DataDog "+
			"(name: %s, tags: %s, rate: %b).\n",
		name,
		tags,
		rate)

	return nil
}
