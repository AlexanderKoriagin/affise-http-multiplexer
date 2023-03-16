package limit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Acquire(t *testing.T) {
	l := &service{max: 8}

	for i := 0; i < 8; i++ {
		assert.True(t, l.Acquire())
	}
	assert.False(t, l.Acquire())
}

func TestService_Release(t *testing.T) {
	l := &service{max: 8}

	for i := 0; i < 8; i++ {
		assert.True(t, l.Acquire())
	}
	assert.False(t, l.Acquire())

	l.Release()
	assert.True(t, l.Acquire())

	for i := 0; i < 8; i++ {
		l.Release()
	}

	for i := 0; i < 8; i++ {
		assert.True(t, l.Acquire())
	}
}
