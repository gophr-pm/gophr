package datadog

import (
	"fmt"
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
func (fc *FakeClient) Gauge(name string, value float64, tags []string, rate float64) error {
	log.Print(fmt.Sprintf(
		"Sending Gauge Metric to DataDog \n Name: %s,\n Value: %b,\n Tags: %s,\n Rate: %b \n",
		name, value, tags, rate,
	))
	return nil
}

// Event mocks FakeClient.Event.
func (fc *FakeClient) Event(e *statsd.Event) error {
	log.Print(fmt.Sprintf(
		"Sending Custom Event to DataDog\n Event: %+v \n", e,
	))
	return nil
}

// Incr mocks FakeClient.Incr.
func (fc *FakeClient) Incr(name string, tags []string, rate float64) error {
	log.Print(fmt.Sprintf(
		"Sending Increment Metric to DataDog \n Name: %s,\n Tags: %s,\n Rate: %b \n",
		name, tags, rate,
	))
	return nil
}
