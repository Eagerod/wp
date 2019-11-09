package wp

import (
	"errors"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestCreation(t *testing.T) {
	e1 := errors.New("This thing didn't work")
	e2 := errors.New("This other thing failed")

	me := MultiErrorFromErrors([]error{e1, e2})

	assert.Equal(
		t,
		"This thing didn't work\nThis other thing failed",
		me.Error(),
	)
	assert.Equal(t, true, me.Exists())
}

func TestIgnoresEmptyErrors(t *testing.T) {
	e1 := errors.New("This thing didn't work")
	e2 := errors.New("This other thing failed")

	me := MultiErrorFromErrors([]error{nil, e1, nil, e2, nil})

	assert.Equal(
		t,
		"This thing didn't work\nThis other thing failed",
		me.Error(),
	)
	assert.Equal(t, true, me.Exists())
}

func TestEmpty(t *testing.T) {
	me := MultiErrorFromErrors([]error{})

	assert.Equal(t, false, me.Exists())
}
