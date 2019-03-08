package httpw

import "net/http"

// Handler defines what a wrapped handler should look like.
type Handler interface {
	// ServeHTTP defines the wrapped handler signature.
	ServeHTTP(*http.Request) (*Response, error)
}

// HandlerFunc has the same signature as the Handler.ServeHTTP
// method and helps with usability.
type HandlerFunc func(*http.Request) (*Response, error)

// ServeHTTP implements the Handler interface.
func (f HandlerFunc) ServeHTTP(r *http.Request) (*Response, error) { return f(r) }
