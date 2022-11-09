package docker

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	mock_docker "github.com/arcalot/arcaflow-plugin-image-builder/mocks/docker"
	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

	assert.Error(t, client.Build("some", "path", []string{"tag1", "tag2"}))
}

func TestClient_ImageTag(t *testing.T) {
	ctrl := gomock.NewController(t)

	dockerClientMock := mock_docker.NewMockDockerClient(ctrl)
	dockerClientMock.EXPECT().
		ImageTag(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(fmt.Errorf("I totally crashed"))

	client := CEClient{
		client: dockerClientMock,
	}

	assert.Error(t, client.Tag("some:path", "sky.io"))
}

func TestClient_ImagePush(t *testing.T) {
	ctrl := gomock.NewController(t)

	dockerClientMock := mock_docker.NewMockDockerClient(ctrl)
	dockerClientMock.EXPECT().
		ImagePush(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(io.NopCloser(strings.NewReader("stream")), fmt.Errorf("I totally crashed"))

	client := CEClient{
		client: dockerClientMock,
	}

	assert.Error(t, client.Push("some:path", "user", "pass", "sky.io"))
}

func TestClient_Show(t *testing.T) {
	stream_jsons, err := os.ReadFile("../../fixtures/jsons/stream_no_errors.jsons")
	if err != nil {
		panic(err)
	}
	stream_txt, err := os.ReadFile("../../fixtures/jsons/stream_no_errors.txt")
	if err != nil {
		panic(err)
	}
	rdr_jsons := io.NopCloser(strings.NewReader(string(stream_jsons)))
	buf := new(bytes.Buffer)
	err = Show(rdr_jsons, buf)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, string(stream_txt), buf.String())

	bad_jsons, err := os.ReadFile("../../fixtures/jsons/bad.jsons")
	if err != nil {
		panic(err)
	}
	rdr_jsons = io.NopCloser(strings.NewReader(string(bad_jsons)))
	assert.Error(t, Show(rdr_jsons, new(bytes.Buffer)))

	bad_jsons, err = os.ReadFile("../../fixtures/jsons/malformed_docker_error_details.jsons")
	if err != nil {
		panic(err)
	}
	rdr_jsons = io.NopCloser(strings.NewReader(string(bad_jsons)))
	assert.Error(t, Show(rdr_jsons, new(bytes.Buffer)))

	bad_jsons, err = os.ReadFile("../../fixtures/jsons/error_details.jsons")
	if err != nil {
		panic(err)
	}
	rdr_jsons = io.NopCloser(strings.NewReader(string(bad_jsons)))
	assert.Error(t, Show(rdr_jsons, new(bytes.Buffer)))
}

func TestNewCEClient(t *testing.T) {
	cec, err := NewCEClient()
	var cec_t *CEClient
	if err != nil {
		assert.Nil(t, cec)
	} else {
		assert.IsType(t, cec_t, cec)
	}
}
