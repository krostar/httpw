package httpw

import (
	"fmt"
	"net/http"
)

// Error implements the error interface and defines a way
// to return a complete error to the final user.
type Error struct {
	// Status is the HTTP status code to return.
	Status int
	// Header stores the headers to set in the response.
	Header http.Header
	// Data is an interface which will be json marshaled
	// and sent to the final user.
	Data interface{}
	// Err stores the internal error to give to the OnErrorFunc callback.
	Err error
}

// E is an alias to easily use Error.
type E = Error

// String implements the stringer interface.
func (e Error) String() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("error %d (%s)", e.Status, http.StatusText(e.Status))
}

// Error implements the error interface.
func (e Error) Error() string {
	return e.String()
}
