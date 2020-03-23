package httpw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerFunc_ServeHTTP(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	resp, err := HandlerFunc(
		func(r *http.Request) (*Response, error) {
			return &R{Status: http.StatusTeapot}, nil
		},
	).ServeHTTP(r)

	require.NoError(t, err)
	assert.Equal(t, http.StatusTeapot, resp.Status)
}
