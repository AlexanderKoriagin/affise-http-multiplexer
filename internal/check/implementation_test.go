package check

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_CheckUrl(t *testing.T) {
	s := &service{urlQuantity: 8}

	assert.False(t, s.CheckUrl([]string{}))
	assert.True(t, s.CheckUrl([]string{"1", "2", "3", "4", "5", "6", "7"}))
	assert.True(t, s.CheckUrl([]string{"1", "2", "3", "4", "5", "6", "7", "8"}))
	assert.False(t, s.CheckUrl([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}))
}
