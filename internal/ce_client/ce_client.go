package ce_client

import (
	"fmt"
	"strings"

	"github.com/arcalot/arcaflow-plugin-image-builder/internal/docker"
)

type ContainerEngineClient interface {
	Build(filepath string, name string, tags []string) error
	Tag(image_tag string, destination string) error
	Push(destination string, username string, password string, registry_address string) error
}

func NewCeClient(choice string) (ContainerEngineClient, error) {
	choice = strings.ToLower(choice)
	switch choice {
	case "podman":
		return nil, fmt.Errorf("podman is not supported yet")
	case "docker-cli":
		return nil, fmt.Errorf("docker CLI is not supported yet")
	default: // docker
		return docker.NewCEClient()
	}
}
