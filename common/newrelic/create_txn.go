package nr

import (
	"log"
	"net/http"

	"github.com/gophr-pm/gophr/common/config"
	newrelic "github.com/newrelic/go-agent"
)

// CreateNewRelicTxn starts a transaction that will log activity
// to new relic.
func CreateNewRelicTxn(
	newRelicApp newrelic.Application,
	conf *config.Config,
	w *http.ResponseWriter,
	r *http.Request,
) newrelic.Transaction {
	var txn newrelic.Transaction
	if !conf.IsDev {
		log.Printf("Logging request for %s \n", r.URL.String())
		// Create a new relic transaction.
		txn = newRelicApp.StartTransaction(r.URL.String(), *w, r)
		defer txn.End()
	}

	return txn
}
