package httpw

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Wrapper handles wrapped function.
type Wrapper struct {
	o *options
}

// New returns a new wrapper.
func New(opts ...Option) *Wrapper {
	var o = options{
		dataMarshaler:      defaultDataMarshaler,
		defaultErrorStatus: http.StatusInternalServerError,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &Wrapper{o: &o}
}

// Wrap wraps a HTTP handler.
func Wrap(handler Handler, opts ...Option) http.Handler {
	return New(opts...).Wrap(handler)
}

// WrapF wraps a HTTP handler func.
func WrapF(handler HandlerFunc, opts ...Option) http.HandlerFunc {
	return New(opts...).WrapF(handler)
}

// WrapF wraps an HTTP handler func.
// nolint: interfacer
// The nolint above is due to the fact that interfacer is complaining about `handler HandlerFunc`
// that could be `handler Handler` as `HandlerFunc` implements `Handler`. The thing is that
// it avoids people to directly pass a function that has the same signature as HandlerFunc but which
// is not of type HandlerFunc so it's gonna stays that way.
func (w Wrapper) WrapF(handler HandlerFunc) http.HandlerFunc {
	return w.Wrap(handler)
}

func defaultDataMarshaler(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	resp, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal data")
	}

	return resp, nil
}
