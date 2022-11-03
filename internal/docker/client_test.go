package docker

import (
	"fmt"
	"testing"

	mock_docker "github.com/arcalot/arcaflow-plugin-image-builder/mocks/docker"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
)

func TestClient_BuildImage(t *testing.T) {
	ctrl := gomock.NewController(t)

	dockerClientMock := mock_docker.NewMockDockerClient(ctrl)
	dockerClientMock.EXPECT().
		ImageBuild(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(types.ImageBuildResponse{}, fmt.Errorf("I totally crashed"))

	client := CEClient{
		client: dockerClientMock,
	}

	client.Build("some", "path", []string{"tag1", "tag2"})
}
