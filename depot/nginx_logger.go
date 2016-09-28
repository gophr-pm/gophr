package main

import (
	"log"

	"github.com/hpcloud/tail"
)

const (
	nginxErrorLogPath  = "/var/log/nginx/error.log"
	nginxAccessLogPath = "/var/log/nginx/access.log"
)

// tailNginxLogs tails all the nginx logs in prints them with the default go
// logger.
func tailNginxLogs() error {
	// Create both of the tails.
	errorLogTail, err := tail.TailFile(
		nginxErrorLogPath,
		tail.Config{Follow: true})
	if err != nil {
		return err
	}
	accessLogTail, err := tail.TailFile(
		nginxAccessLogPath,
		tail.Config{Follow: true})
	if err != nil {
		return err
	}

	// Print logs in the background.
	go tailLinePrinter(errorLogTail, "[nginx:error]")
	go tailLinePrinter(accessLogTail, "[nginx:access]")

	// Exit successfully.
	return nil
}

// tailLinePrinter prints all lines from a tail.
func tailLinePrinter(t *tail.Tail, prefix string) {
	for line := range t.Lines {
		log.Println(prefix, line.Text)
	}
}
