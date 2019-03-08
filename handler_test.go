package httpw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerFunc_ServeHTTP(t *testing.T) {
	var (
		r = httptest.NewRequest(http.MethodGet, "/", nil)
		h = func(r *http.Request) (*Response, error) {
			return &R{Status: http.StatusTeapot}, nil
		}
	)

	resp, err := HandlerFunc(h).ServeHTTP(r)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTeapot, resp.Status)
}
