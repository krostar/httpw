package httpw

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkWithoutWrapper(b *testing.B) {
	var handler = func(rw http.ResponseWriter, _ *http.Request) {
		b, err := defaultDataMarshaler("42")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusAccepted)
		rw.Write(b) // nolint: errcheck, gosec
	}
	bench(b, http.HandlerFunc(handler))
}

func BenchmarkWithWrapper(b *testing.B) {
	var (
		w       = New()
		handler = func(r *http.Request) (*R, error) {
			return &R{
				Status: http.StatusAccepted,
				Data:   "42",
			}, nil
		}
	)
	bench(b, w.WrapF(handler))
}

func bench(b *testing.B, h http.Handler) {
	for n := 0; n < b.N; n++ {
		var (
			rw = httptest.NewRecorder()
			r  = httptest.NewRequest(http.MethodGet, "/", nil)
		)
		h.ServeHTTP(rw, r)
	}
}
