package docker_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"go.arcalot.io/arcaflow-container-toolkit/internal/docker"

	"github.com/docker/docker/api/types/build"
	mock_docker "go.arcalot.io/arcaflow-container-toolkit/mocks/docker"
	"go.arcalot.io/assert"
	"go.uber.org/mock/gomock"
)

func TestClient_BuildImage(t *testing.T) {
	ctrl := gomock.NewController(t)

	dockerClientMock := mock_docker.NewMockDockerClient(ctrl)
	dockerClientMock.EXPECT().
		ImageBuild(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(build.ImageBuildResponse{}, fmt.Errorf("I totally crashed"))

	client := docker.CEClient{
		Client: dockerClientMock,
	}

	assert.Error(t, client.Build("some", "path", []string{"tag1", "tag2"},
		"amd64", docker.DefaultBuildOptions()))
}

func TestClient_ImageTag(t *testing.T) {
	ctrl := gomock.NewController(t)

	dockerClientMock := mock_docker.NewMockDockerClient(ctrl)
	dockerClientMock.EXPECT().
		ImageTag(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1).
		Return(fmt.Errorf("I totally crashed"))

	client := docker.CEClient{
		Client: dockerClientMock,
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

	client := docker.CEClient{
		Client: dockerClientMock,
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
	err = docker.Show(rdr_jsons, buf)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equals(t, string(stream_txt), buf.String())

	bad_jsons, err := os.ReadFile("../../fixtures/jsons/bad.jsons")
	if err != nil {
		panic(err)
	}
	rdr_jsons = io.NopCloser(strings.NewReader(string(bad_jsons)))
	assert.Error(t, docker.Show(rdr_jsons, new(bytes.Buffer)))

	bad_jsons, err = os.ReadFile("../../fixtures/jsons/malformed_docker_error_details.jsons")
	if err != nil {
		panic(err)
	}
	rdr_jsons = io.NopCloser(strings.NewReader(string(bad_jsons)))
	assert.Error(t, docker.Show(rdr_jsons, new(bytes.Buffer)))

	bad_jsons, err = os.ReadFile("../../fixtures/jsons/error_details.jsons")
	if err != nil {
		panic(err)
	}
	rdr_jsons = io.NopCloser(strings.NewReader(string(bad_jsons)))
	assert.Error(t, docker.Show(rdr_jsons, new(bytes.Buffer)))
}

func TestNewCEClient(t *testing.T) {
	cec, _ := docker.NewCEClient()
	assert.NotNil(t, cec)
}
