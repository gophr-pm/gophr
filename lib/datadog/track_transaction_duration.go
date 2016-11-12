package datadog

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

// TrackTransaction is responsible for tracking a specific request TrackTransaction.
// This can be internal requests, or external user-facing requests.
// It will always report the transaction duration and optionally an
// associated event.
func TrackTransaction(args TrackTranscationArgs) {
	// Calculate the total request duration.
	reqDuration := float64(time.Since(args.StartTime))
	log.Println("Transaction duration = ", reqDuration)

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
	if args.CreateEvent != nil {
		// Create the base event.
		event := args.CreateEvent(
			args.CustomEventName,
			strings.Join(args.EventInfo[:], ","),
		)

		// Identify the event alert type.
		switch args.AlertType {
		case "success":
			event.AlertType = statsd.Success
			//  count for error
		case "error":
			event.AlertType = statsd.Error
			// send cound for error
		default:
			event.AlertType = statsd.Info
		}

		// Report an increment count for the given event type.
		if err := args.Client.Incr(
			fmt.Sprintf(
				"%s.%s",
				args.CustomEventName,
				args.AlertType,
			), args.Tags, 1); err != nil {
			log.Println(err)
		}

		// Report the event to Datadog.
		if err := args.Client.Event(event); err != nil {
			log.Println(err)
		}
	}
}
