package httpw

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapper_Wrap(t *testing.T) {
	headers := make(http.Header)
	headers.Set("yolo", "yili")

	tests := map[string]struct {
		r               *R
		e               error
		expectedStatus  int
		expectedBody    []byte
		expectedHeaders http.Header
		expectedErr     error
	}{
		"nil": {
			expectedHeaders: make(http.Header),
		},
		"no errors": {
			r: &R{
				Status: http.StatusUnavailableForLegalReasons,
				Data:   "hello world",
				Header: headers,
			},
			expectedStatus:  http.StatusUnavailableForLegalReasons,
			expectedBody:    []byte(`"hello world"`),
			expectedHeaders: headers,
		},
		"errors": {
			e: &E{
				Status: http.StatusUnavailableForLegalReasons,
				Data:   "hello world",
				Header: headers,
				Err:    errors.New("eww"),
			},
			expectedStatus:  http.StatusUnavailableForLegalReasons,
			expectedBody:    []byte(`"hello world"`),
			expectedHeaders: headers,
			expectedErr:     errors.New("eww"),
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			wrapper := New(WithOnErrorCallback(func(r *http.Request, err error) {
				require.Equal(t, test.expectedErr, err)
			}))

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			rw := httptest.NewRecorder()

			wrapper.Wrap(HandlerFunc(
				func(r *http.Request) (*R, error) {
					return test.r, test.e
				},
			)).ServeHTTP(rw, r)

			require.Equal(t, test.expectedStatus, rw.Code)
			assert.Equal(t, string(test.expectedBody), rw.Body.String())
			assert.Equal(t, test.expectedHeaders, rw.Header())
		})
	}
}

func TestWrapper_unwrap(t *testing.T) {
	wrapper := New()
	tests := map[string]struct {
		handler          HandlerFunc
		mockDMBytes      []byte
		mockDMErr        error
		expectedStatus   int
		expectedBody     []byte
		expectedErrorStr string
	}{
		"nil": {
			handler: func(r *http.Request) (*R, error) { return nil, nil },
		},
		"nil response, std error": {
			handler:          func(r *http.Request) (*R, error) { return nil, errors.New("eww") },
			expectedStatus:   http.StatusInternalServerError,
			expectedErrorStr: "eww",
		},
		"nil response, Error": {
			handler: func(r *http.Request) (*R, error) {
				return nil, E{
					Status: http.StatusBadGateway, Data: "42", Err: errors.New("eww"),
				}
			},
			expectedStatus:   http.StatusBadGateway,
			expectedBody:     []byte(`"42"`),
			expectedErrorStr: "eww",
		},
		"nil response, *Error": {
			handler: func(r *http.Request) (*R, error) {
				return nil, &E{
					Status: http.StatusBadGateway, Data: "42", Err: errors.New("eww"),
				}
			},
			expectedStatus:   http.StatusBadGateway,
			expectedBody:     []byte(`"42"`),
			expectedErrorStr: "eww",
		},
		"nil response, Error without status": {
			handler: func(r *http.Request) (*R, error) {
				return nil, E{Err: errors.New("eww")}
			},
			expectedStatus:   wrapper.o.defaultErrorStatus,
			expectedErrorStr: "eww",
		},
		"nil response, Error without error": {
			handler: func(r *http.Request) (*R, error) {
				return nil, E{}
			},
			expectedStatus: wrapper.o.defaultErrorStatus,
		},
		"Response, nil error": {
			handler: func(r *http.Request) (*R, error) {
				return &R{
					Status: http.StatusAccepted,
					Data:   "42",
				}, nil
			},
			expectedStatus: http.StatusAccepted,
			expectedBody:   []byte(`"42"`),
		},
		"Response, E": {
			handler: func(r *http.Request) (*R, error) {
				return &R{
						Status: http.StatusAccepted,
						Data:   "42",
					}, &E{
						Status: http.StatusBadRequest,
						Err:    errors.New("eww"),
					}
			},
			expectedStatus:   http.StatusBadRequest,
			expectedErrorStr: "eww",
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r := httptest.NewRequest(http.MethodGet, "/", nil)

			status, body, _, err := wrapper.unwrap(test.handler, r)
			assert.Equal(t, test.expectedStatus, status)
			assert.Equal(t, string(test.expectedBody), string(body))
			if test.expectedErrorStr != "" {
				require.NotNil(t, err)
				assert.Equal(t, test.expectedErrorStr, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
