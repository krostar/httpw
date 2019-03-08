package httpw

import (
	"net/http"
)

type options struct {
	onError            []OnErrorFunc
	dataMarshaler      DataMarshalerFunc
	defaultErrorStatus int
}

// OnErrorFunc defines the prototype on the function called when an error occur.
type OnErrorFunc func(r *http.Request, err error)

// DataMarshalerFunc defines the prototype on the function called to marshal data.
type DataMarshalerFunc func(data interface{}) ([]byte, error)

// Option defines the prototype of the option applier.
type Option func(*options)

// WithOnErrorCallback sets the callback functions called when an error occur.
func WithOnErrorCallback(fct OnErrorFunc) Option {
	return func(o *options) { o.onError = append(o.onError, fct) }
}

// WithDataMarshaler sets the function called to marshal data.
func WithDataMarshaler(fct DataMarshalerFunc) Option {
	return func(o *options) { o.dataMarshaler = fct }
}

// WithDefaultErrorStatus defines the default status set for error.
func WithDefaultErrorStatus(status int) Option {
	return func(o *options) { o.defaultErrorStatus = status }
}
