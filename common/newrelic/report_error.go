package nr

import (
	"log"

	newrelic "github.com/newrelic/go-agent"
)

// ReportNewRelicError reports errors to newrelic.
func ReportNewRelicError(txn newrelic.Transaction, err error, isDev bool) {
	if !isDev {
		log.Println("Reporting error to newrelic: ", err)
		txn.NoticeError(err)
	}
}
