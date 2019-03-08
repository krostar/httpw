package httpw

import "net/http"

// Response defines which response to give to the caller.
type Response struct {
	// Status is the HTTP status code to return.
	Status int
	// Header stores the headers to set in the response.
	Header http.Header
	// Data is an interface which will be json marshaled
	// to return to the final user.
	Data interface{}
}

// R is an alias to easily use Response.
type R = Response
