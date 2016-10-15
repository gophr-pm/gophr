package common

import "bytes"

// MockHTTPResponseBody is a utility to make mocking http.Responses easier.
type MockHTTPResponseBody struct {
	bytes.Buffer
	closed bool
}

// NewMockHTTPResponseBody creates a new MockHTTPResponseBody with some initial
// data to start off with.
func NewMockHTTPResponseBody(data []byte) *MockHTTPResponseBody {
	mhrb := &MockHTTPResponseBody{}
	mhrb.Write(data)
	return mhrb
}

// Close marks this response body closed.
func (mhrb *MockHTTPResponseBody) Close() error {
	mhrb.closed = true
	return nil
}

// WasClosed returns true if the Close() method was invoked on this instance.
func (mhrb *MockHTTPResponseBody) WasClosed() bool {
	return mhrb.closed
}
