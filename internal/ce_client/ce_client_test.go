package ce_client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func IntMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestIntMinBasic(t *testing.T) {
	assert.Equal(t, -2, IntMin(2, -2))
}
