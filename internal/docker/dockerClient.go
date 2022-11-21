package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
)

type DockerClient interface {
	ImageBuild(context.Context, io.Reader, types.ImageBuildOptions) (types.ImageBuildResponse, error)
	ImageTag(context.Context, string, string) error
	ImagePush(context.Context, string, types.ImagePushOptions) (io.ReadCloser, error)
}
