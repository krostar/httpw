package httpw_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/krostar/httpw"
)

func Example() {
	var (
		wrapper = httpw.New(
			httpw.WithOnErrorCallback(func(r *http.Request, err error) {
				fmt.Println("err", err)
			}),
		)
		recorder = httptest.NewRecorder()
		request  = httptest.NewRequest(http.MethodGet, "/", nil)
	)

	wrapper.WrapF(func(r *http.Request) (*httpw.R, error) {
		return nil, errors.New("boum")
	}).ServeHTTP(recorder, request)

	fmt.Println("status", recorder.Code)

	// Output:
	// err boum
	// status 500
}
