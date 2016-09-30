package nr

import (
	"log"

	newrelic "github.com/newrelic/go-agent"
)

// ReportNewRelicError report
func ReportNewRelicError(txn newrelic.Transaction, err error, isDev bool) {
	if !isDev {
		log.Println("Reporting error to newrelic: ", err)
		txn.NoticeError(err)
	}
}
