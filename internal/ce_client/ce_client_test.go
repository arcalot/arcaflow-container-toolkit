package ce_client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCeClient(t *testing.T) {
	_, err := NewCeClient("podman")
	assert.Error(t, err)
	_, err = NewCeClient("docker")
	assert.NoError(t, err)
}
