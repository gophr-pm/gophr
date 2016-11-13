package datadog

import (
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

// Constants used to references against statsd's unexported
// AlertType.
const (
	Error   = "error"
	Info    = "info"
	Success = "success"
)

// EventCreator is a interface for passing a statsd DataDog NewEvent
// and responsible - Google Search for creating new metric events.
type EventCreator func(title, text string) *statsd.Event

// TrackTranscationArgs is the args structs for tracking transactions
// to DataDog.
type TrackTranscationArgs struct {
	Tags            []string
	MetricName      string
	Client          *statsd.Client
	StartTime       time.Time
	AlertType       string
	EventInfo       []string
	CreateEvent     EventCreator
	CustomEventName string
}
