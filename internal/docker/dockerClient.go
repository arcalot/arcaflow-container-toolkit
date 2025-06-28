package docker

import (
	"context"
	"github.com/docker/docker/api/types/build"
	"io"

	"github.com/docker/docker/api/types/image"
)

type DockerClient interface {
	ImageBuild(context.Context, io.Reader, build.ImageBuildOptions) (build.ImageBuildResponse, error)
	ImageTag(context.Context, string, string) error
	ImagePush(context.Context, string, image.PushOptions) (io.ReadCloser, error)
}
