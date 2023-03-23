package ce_service_test

import (
	"testing"

	"go.arcalot.io/arcaflow-container-toolkit/internal/ce_service"

	"go.arcalot.io/assert"
)

func TestNewCeClient(t *testing.T) {
	_, err := ce_service.NewContainerEngineService("podman")
	assert.Error(t, err)
	_, err = ce_service.NewContainerEngineService("docker")
	assert.NoError(t, err)
}
