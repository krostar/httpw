package httpw

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithOnErrorCallback(t *testing.T) {
	var w = New(WithOnErrorCallback(nil))
	require.Len(t, w.o.onError, 1)
	assert.Nil(t, w.o.onError[0])
}

func TestWithDataMarshaler(t *testing.T) {
	var w = New(WithDataMarshaler(nil))
	assert.Nil(t, w.o.dataMarshaler)
}

type marshalFailer int

// nolint: unparam
func (marshalFailer) MarshalJSON() ([]byte, error) { return nil, errors.New("fail") }

func TestWithDefaultErrorStatus(t *testing.T) {
	var w = New(WithDefaultErrorStatus(http.StatusTeapot))
	assert.Equal(t, http.StatusTeapot, w.o.defaultErrorStatus)
}
