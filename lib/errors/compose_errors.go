package errors

import (
	"bytes"
	"errors"
	"strconv"
)

// ComposeErrors joins all the accumulated individual errors into one combined
// error.
func ComposeErrors(msg string, errs []error) error {
	buffer := bytes.Buffer{}

	buffer.WriteString(msg)
	buffer.WriteString(". Bumped into ")
	buffer.WriteString(strconv.Itoa(len(errs)))
	buffer.WriteString(" problems: [ ")

	for i, err := range errs {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(err.Error())
	}

	buffer.WriteString(" ].")

	return errors.New(buffer.String())
}
