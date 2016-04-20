package main

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/skeswa/gophr/common"
)

const (
	urlVarAuthor = "author"
	urlVarRepo   = "repo"
)

func extractPackageRequestMetadata(r *http.Request) (string, string, error) {
	vars := mux.Vars(r)

	author := vars[urlVarAuthor]
	if len(author) < 0 {
		return "", "", NewInvalidURLParameterError(urlVarAuthor, author)
	}

	repo := vars[urlVarRepo]
	if len(repo) < 0 {
		return "", "", NewInvalidURLParameterError(urlVarRepo, repo)
	}

	return author, repo, nil
}

// VersionsHandler creates an HTTP request handler that responds to requests for
// all the versions of a package.
func VersionsHandler(session *gocql.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		author, repo, err := extractPackageRequestMetadata(r)
		if err != nil {
			common.RespondWithError(w, err)
			return
		}

		// Look in the database first.
		versions, err := common.FindPackageVersions(session, author, repo)
		if err != nil {
			common.RespondWithError(w, err)
			return
		} else if versions != nil && len(versions) > 0 {
			versionListDTO := common.NewVersionListDTOFromVersionStrings(versions)
			json, err := versionListDTO.MarshalJSON()
			if err != nil {
				common.RespondWithError(w, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(json)
			return
		} else {
			// We didn't find anything in the database, so ask github.
			refs, err := common.FetchRefs(author, repo)
			if err != nil {
				common.RespondWithError(w, err)
				return
			}

			if refs.Candidates != nil && len(refs.Candidates) > 0 {
				// TODO(skeswa): this means we found versions that we didn't know about
				// so this needs to  be put into the db for efficieny's sake.
				versionListDTO := common.NewVersionListDTOFromSemverCandidateList(refs.Candidates)
				json, err := versionListDTO.MarshalJSON()
				if err != nil {
					common.RespondWithError(w, err)
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write(json)
				return
			}

			// No versions could be found anywhere.
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(common.EmptyListJSON))
			return
		}
	}
}

// LatestVersionHandler creates an HTTP request handler that responds to
// requests for the latest version of a package.
func LatestVersionHandler(session *gocql.Session) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		author, repo, err := extractPackageRequestMetadata(r)
		if err != nil {
			common.RespondWithError(w, err)
			return
		}

		// Look in the database first.
		versions, err := common.FindPackageVersions(session, author, repo)
		if err != nil {
			common.RespondWithError(w, err)
			return
		} else if versions != nil && len(versions) > 0 {
			// In the database, the list of versions are sorted ascendingly.
			lastVersion := versions[len(versions)-1]
			versionDTO := common.NewVersionDTO(lastVersion)
			json, err := versionDTO.MarshalJSON()
			if err != nil {
				common.RespondWithError(w, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(json)
			return
		} else {
			// We didn't find anything in the database, so ask github.
			refs, err := common.FetchRefs(author, repo)
			if err != nil {
				common.RespondWithError(w, err)
				return
			}

			if refs.Candidates != nil && len(refs.Candidates) > 0 {
				lastCandidate := refs.Candidates.Highest()
				versionDTO := common.NewVersionDTO(lastCandidate.String())
				json, jsonErr := versionDTO.MarshalJSON()
				if jsonErr != nil {
					common.RespondWithError(w, jsonErr)
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write(json)
				return
			}

			versionDTO := common.NewVersionDTO(refs.MasterRefHash)
			json, err := versionDTO.MarshalJSON()
			if err != nil {
				common.RespondWithError(w, err)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(json)
			return
		}
	}
}
