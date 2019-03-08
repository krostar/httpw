package httpw

import (
	"net/http"

	"github.com/pkg/errors"
)

// Wrap wraps a classic http.Handler to be used as a wrapped Handler.
func (w Wrapper) Wrap(h Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		status, responseData, responseHeader, err := w.unwrap(h, r)
		if err != nil {
			for _, onError := range w.o.onError {
				onError(r, err)
			}
		}

		rwH := rw.Header()
		for key, values := range responseHeader {
			for _, value := range values {
				rwH.Add(key, value)
			}
		}

		rw.WriteHeader(status)
		if responseData != nil {
			rw.Write(responseData) // nolint: errcheck, gosec
		}
	}
}

func (w Wrapper) unwrap(h Handler, r *http.Request) (int, []byte, http.Header, error) {
	var (
		status int
		data   interface{}
		header http.Header
		err    error
	)

	wrappedResponse, handlerErr := h.ServeHTTP(r)
	if handlerErr != nil {
		switch errT := handlerErr.(type) {
		case *Error:
			status, data, header, err = errT.Status, errT.Data, errT.Header, errT.Err
		case Error:
			status, data, header, err = errT.Status, errT.Data, errT.Header, errT.Err
		default:
			err = errT
		}
		if status == 0 {
			status = w.o.defaultErrorStatus
		}
	} else if wrappedResponse != nil {
		status, data, header = wrappedResponse.Status, wrappedResponse.Data, wrappedResponse.Header
	}

	resp, mErr := w.o.dataMarshaler(data)
	if mErr != nil {
		if err != nil {
			mErr = errors.Wrap(err, mErr.Error()+", original error was: ")
		}
		err = mErr
	}
	return status, resp, header, err
}
