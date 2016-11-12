package datadog

import (
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

const (
	Error   = "error"
	Info    = "info"
	Success = "success"
)

// EventCreator lol
type EventCreator func(title, text string) *statsd.Event

// TrackTranscationArgs lol
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
