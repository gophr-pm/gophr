package awesome

import (
	"errors"
	"testing"

	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/lib/db/query"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPersistAwesomeGoPackage(t *testing.T) {
	Convey("Given a set of awesome packages", t, func() {

		Convey("if batch executor fails, it should return an error", func() {
			err := PersistAwesomePackages(
				PersistAwesomePackagesArgs{
					Session: &gocql.Session{},
					NewBatchCreator: func(batchType gocql.BatchType) *gocql.Batch {
						return &gocql.Batch{}
					},
					BatchExecutor: func(
						session query.BatchingQueryable,
						currentBatch *gocql.Batch,
						resultsError chan error,
					) {
						resultsError <- errors.New("Executing batch failed")
					},
					PackageTuples: generateRandomAwesomePackages(100),
				},
			)

			So(err, ShouldNotBeNil)
		})

		Convey("if batch executor completes with a package number divisable by 50, since the url failed to return", func() {
			err := PersistAwesomePackages(
				PersistAwesomePackagesArgs{
					Session: &gocql.Session{},
					NewBatchCreator: func(batchType gocql.BatchType) *gocql.Batch {
						return &gocql.Batch{}
					},
					BatchExecutor: func(
						session query.BatchingQueryable,
						currentBatch *gocql.Batch,
						resultsError chan error,
					) {
						resultsError <- nil
					},
					PackageTuples: generateRandomAwesomePackages(100),
				},
			)

			So(err, ShouldBeNil)
		})

		Convey("if batch executors pass with a package number not divisable by 50, since the url failed to return", func() {
			err := PersistAwesomePackages(
				PersistAwesomePackagesArgs{
					Session: &gocql.Session{},
					NewBatchCreator: func(batchType gocql.BatchType) *gocql.Batch {
						return &gocql.Batch{}
					},
					BatchExecutor: func(
						session query.BatchingQueryable,
						currentBatch *gocql.Batch,
						resultsError chan error,
					) {
						resultsError <- nil
					},
					PackageTuples: generateRandomAwesomePackages(201),
				},
			)

			So(err, ShouldBeNil)
		})
	})
}
