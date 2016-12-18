package main

import (
	"strings"
	"time"

	gosxnotifier "github.com/gosx-notifier"
)

const (
	gophrTitle                     = "gophrctl finished running"
	gophrAlertSound                = gosxnotifier.Default
	gophrAppIconPath               = "./gophr.png"
	notificationThresholdInSeconds = 10
)

func showNotification(startTime time.Time, args []string) {
	// Calculate our time difference since initially running a command.
	timeDiff := time.Since(startTime).Seconds()
	// Determine if we should push an alert.
	shouldAlert := timeDiff > notificationThresholdInSeconds

	// Only alert if the time thershold has been passed.
	if shouldAlert {
		// Create a notification.
		note := &gosxnotifier.Notification{
			Title:   gophrTitle,
			Sound:   gophrAlertSound,
			Message: strings.Join(args, " "),
			AppIcon: gophrAppIconPath,
		}

		// Push notification to the terminal.
		note.Push()
	}
}
