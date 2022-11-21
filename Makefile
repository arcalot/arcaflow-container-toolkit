generate-mocks:
	mockgen -source=internal/docker/dockerClient.go -destination=mocks/docker/dockerClient.go
	mockgen -source=internal/ce_client/ce_client.go -destination=mocks/ce_client/ce_client.go