package httpw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	w := New()
	assert.Equal(t, http.StatusInternalServerError, w.o.defaultErrorStatus)
	assert.NotNil(t, w.o.dataMarshaler)

	w = New(WithDefaultErrorStatus(http.StatusBadRequest))
	assert.Equal(t, http.StatusBadRequest, w.o.defaultErrorStatus)
}

func TestWrap(t *testing.T) {
	var (
		h = HandlerFunc(func(r *http.Request) (*Response, error) {
			return &R{Status: http.StatusTeapot}, nil
		})
		recorder = httptest.NewRecorder()
		request  = httptest.NewRequest(http.MethodGet, "/", nil)
	)

	Wrap(h).ServeHTTP(recorder, request)

	require.Equal(t, http.StatusTeapot, recorder.Code)
}

func TestWrapF(t *testing.T) {
	var (
		h = func(r *http.Request) (*Response, error) {
			return &R{Status: http.StatusTeapot}, nil
		}
		recorder = httptest.NewRecorder()
		request  = httptest.NewRequest(http.MethodGet, "/", nil)
	)

	WrapF(h).ServeHTTP(recorder, request)

	require.Equal(t, http.StatusTeapot, recorder.Code)
}

func TestWrapper_WrapF(t *testing.T) {
	var (
		w = New()
		h = func(r *http.Request) (*Response, error) {
			return &R{Status: http.StatusTeapot}, nil
		}
		recorder = httptest.NewRecorder()
		request  = httptest.NewRequest(http.MethodGet, "/", nil)
	)

	w.WrapF(h).ServeHTTP(recorder, request)

	require.Equal(t, http.StatusTeapot, recorder.Code)
}

func TestDefaultDataMarshaler(t *testing.T) {
	var tests = map[string]struct {
		data            interface{}
		expectedRepr    []byte
		expectedFailure bool
		expectedErrRepr string
	}{
		"success": {
			data:            "42",
			expectedRepr:    []byte(`"42"`),
			expectedFailure: false,
		},
		"empty data": {
			data:            nil,
			expectedFailure: false,
		},
		"marshal failure": {
			data:            marshalFailer(42),
			expectedFailure: true,
			expectedErrRepr: "unable to marshal data: " +
				"json: error calling MarshalJSON for type httpw.marshalFailer: " +
				"fail",
		},
	}

	for name, test := range tests {
		var test = test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			repr, err := defaultDataMarshaler(test.data)
			if test.expectedFailure {
				require.Error(t, err)
				assert.Equal(t, test.expectedErrRepr, err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedRepr, repr)
			}
		})
	}
}
