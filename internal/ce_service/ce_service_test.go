package ce_service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCeClient(t *testing.T) {
	_, err := NewContainerEngineService("podman")
	assert.Error(t, err)
	_, err = NewContainerEngineService("docker")
	assert.NoError(t, err)
}
