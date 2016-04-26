package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common"
)

// RecordInstallHandler creates an HTTP request handler that responds to install
// recording requests.
func RecordInstallHandler(session *gocql.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			common.RespondWithError(w, NewInvalidRequestBodyError(
				"a valid PackageInstallDTO",
				"",
				err,
			))
			return
		}

		dto := &common.PackageInstallDTO{}
		err = dto.UnmarshalJSON(body)
		if err != nil {
			common.RespondWithError(w, NewInvalidRequestBodyError(
				"a valid PackageInstallDTO",
				string(body[:]),
				err,
			))
			return
		}

		if len(dto.Author) < 1 && len(dto.Repo) < 1 {
			common.RespondWithError(w, NewInvalidRequestBodyError(
				"a valid PackageInstallDTO",
				string(body[:]),
			))
			return
		}

		err = recordPackageInstall(session, dto.Author, dto.Repo)
		if err != nil {
			common.RespondWithError(w, err)
			return
		}

		w.WriteHeader(200)
	}
}
