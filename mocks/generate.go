package mocks

// this generates a mock for the DockerClient interface
//go:generate mockgen -destination=mocks/docker/dockerClient.go -package=mocks github.com/arcalot/arcaflow-plugin-image-builder/internal/docker DockerClient

// generate a mock for the ContainerEngineService interface
//go:generate mockgen -destination=mocks/ce_service/ce_service.go -package=mocks github.com/arcalot/arcaflow-plugin-image-builder/internal/ce_service ContainerEngineService
