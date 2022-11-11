package ce_service_test

import (
	"github.com/arcalot/arcaflow-plugin-image-builder/internal/ce_service"
	"testing"

	"go.arcalot.io/assert"
)

func TestNewCeClient(t *testing.T) {
	_, err := ce_service.NewContainerEngineService("podman")
	assert.Error(t, err)
	_, err = ce_service.NewContainerEngineService("docker")
	assert.NoError(t, err)
}
