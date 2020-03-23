package httpw

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError_String(t *testing.T) {
	e := &E{
		Status: http.StatusNotFound,
	}
	assert.Equal(t, "error 404 (Not Found)", e.Error())

	e = &E{
		Status: http.StatusNotFound,
		Err:    errors.New("eww"),
	}
	assert.Equal(t, "eww", e.Error())
}

func TestError_Error(t *testing.T) {
	e := Error{
		Status: http.StatusNotFound,
		Err:    errors.New("eww"),
	}

	assert.Equal(t, e.String(), e.Error())
}
