package datadog

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

// EventCreator is a interface for passing a statsd DataDog NewEvent
// and responsible - Google Search for creating new metric events.
type EventCreator func(title, text string) *statsd.Event

// TrackTransactionArgs is the args structs for tracking transactions
// to DataDog.
type TrackTransactionArgs struct {
	Tags            []string
	Client          Client
	StartTime       time.Time
	AlertType       string
	EventInfo       []string
	MetricName      string
	CreateEvent     EventCreator
	CustomEventName string
}

// TrackTransaction is responsible for tracking a specific request
// TrackTransaction. This can be internal requests, or external user-facing
// requests. It will always report the transaction duration and optionally an
// associated event.
func TrackTransaction(args TrackTransactionArgs) {
	// Calculate the total request duration.
	reqDuration := float64(time.Since(args.StartTime) / time.Millisecond)

	// Then send a guage for the request metric.
	if err := args.Client.Gauge(
		args.MetricName,
		reqDuration,
		args.Tags,
		1,
	); err != nil {
		log.Println(err)
	}

	// Check if we need to send a custom event as well.
	if args.CreateEvent == nil {
		return
	}

	// Create the base event.
	event := args.CreateEvent(
		args.CustomEventName,
		strings.Join(args.EventInfo[:], ","))

	// Pass the metric tags into the event tags.
	event.Tags = args.Tags

	// Identify the event alert type.
	switch args.AlertType {
	case "success":
		event.AlertType = statsd.Success
	case "error":
		event.AlertType = statsd.Error
	default:
		event.AlertType = statsd.Info
	}

	// Report an increment count for the given event type.
	if err := args.Client.Incr(
		fmt.Sprintf("%s.%s", args.CustomEventName, args.AlertType),
		args.Tags,
		1,
	); err != nil {
		log.Println(err)
	}

	// Report the event to Datadog.
	if err := args.Client.Event(event); err != nil {
		log.Println(err)
	}
}
