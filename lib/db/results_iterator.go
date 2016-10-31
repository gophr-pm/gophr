package db

// ResultsIterator iterates through the results of a query.
type ResultsIterator interface {
	// Scan consumes the next row of the iterator and copies the columns of the
	// current row into the values pointed at by dest. Use nil as a dest value to
	// skip the corresponding column. Scan might send additional queries to the
	// database to retrieve the next set of rows if paging was enabled.
	//
	// Scan returns true if the row was successfully unmarshaled or false if the
	// end of the result set was reached or if an error occurred. Close should be
	// called afterwards to retrieve any potential errors.
	Scan(dest ...interface{}) bool
	// Close closes the iterator and returns any errors that happened during the
	// query or the iteration.
	Close() error
}
