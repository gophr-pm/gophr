package awesome

import (
	"errors"
	"testing"

	"github.com/gophr-pm/gophr/lib/db"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestPersistAwesomeGoPackage(t *testing.T) {
	Convey("Given a set of awesome packages", t, func() {

		Convey("if batch executor fails, it should return an error", func() {
			batch := db.NewMockBatch()
			client := db.NewMockClient()

			batch.On(
				"Query",
				"insert into gophr.awesome_packages (author,repo) values (?,?)",
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"))
			client.On("NewLoggedBatch").Return(batch)

			err := persistAwesomePackages(
				persistAwesomePackagesArgs{
					q: client,
					batchExecutor: func(
						currentBatch db.Batch,
						resultsError chan error,
					) {
						resultsError <- errors.New("Executing batch failed")
					},
					packageTuples: generateRandomAwesomePackages(100),
				},
			)

			So(err, ShouldNotBeNil)
		})

		Convey("if batch executor completes with a package number divisable by 50, since the url failed to return", func() {
			batch := db.NewMockBatch()
			client := db.NewMockClient()

			batch.On(
				"Query",
				"insert into gophr.awesome_packages (author,repo) values (?,?)",
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"))
			client.On("NewLoggedBatch").Return(batch)

			err := persistAwesomePackages(
				persistAwesomePackagesArgs{
					q: client,
					batchExecutor: func(
						currentBatch db.Batch,
						resultsError chan error,
					) {
						resultsError <- nil
					},
					packageTuples: generateRandomAwesomePackages(100),
				},
			)

			So(err, ShouldBeNil)
		})

		Convey("if batch executors pass with a package number not divisable by 50, since the url failed to return", func() {
			batch := db.NewMockBatch()
			client := db.NewMockClient()

			batch.On(
				"Query",
				"insert into gophr.awesome_packages (author,repo) values (?,?)",
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"))
			client.On("NewLoggedBatch").Return(batch)

			err := persistAwesomePackages(
				persistAwesomePackagesArgs{
					q: client,
					batchExecutor: func(
						currentBatch db.Batch,
						resultsError chan error,
					) {
						resultsError <- nil
					},
					packageTuples: generateRandomAwesomePackages(201),
				},
			)

			So(err, ShouldBeNil)
		})
	})
}
