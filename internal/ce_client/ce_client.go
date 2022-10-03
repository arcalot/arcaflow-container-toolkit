package ce_client

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type ContainerEngineClient interface {
	Build(filepath string, name string, tags []string) error
	Tag(image_tag string, destination string) error
	Push(destination string, username string, password string, registry_address string) error
}

type docker struct {
	client *client.Client
}

func NewCeClient(choice string) (ContainerEngineClient, error) {
	choice = strings.ToLower(choice)
	switch choice {
	case "podman":
		return nil, fmt.Errorf("Podman is not supported yet.")
	case "docker-cli":
		return nil, fmt.Errorf("Docker CLI is not supported yet.")
	default: // docker
		container_cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return nil, fmt.Errorf("error while creating Docker client (%w)", err)
		}
		return docker{client: container_cli}, nil
	}
}

func (ce docker) Build(filepath string, name string, tags []string) error {
	image_tag := name + ":" + tags[0]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	tar, err := archive.TarWithOptions(filepath, &archive.TarOptions{})
	if err != nil {
		return err
	}
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{image_tag},
	}
	res, err := ce.client.ImageBuild(ctx, tar, opts)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	defer res.Body.Close()
	err = Show(res.Body)
	if err != nil {
		return err
	}
	return nil
}

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type StreamLine struct {
	Stream string `json:"stream"`
}

func Show(rd io.Reader) error {
	var lastLine string
	var nextLine string
	scanner := bufio.NewScanner(rd)
	line := &StreamLine{}

	for scanner.Scan() {
		lastLine = scanner.Text()
		nextLine = scanner.Text()
		json.Unmarshal([]byte(nextLine), &line)
		if _, err := os.Stdout.Write([]byte(line.Stream)); err != nil {
			return err
		}
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (ce docker) Tag(image_tag string, destination string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	err := ce.client.ImageTag(ctx, image_tag, destination)
	if err != nil {
		return err
	}
	return nil
}

func (ce docker) Push(destination string, username string, password string, registry_address string) error {
	authConfig := types.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: registry_address,
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)
	opts := types.ImagePushOptions{RegistryAuth: authConfigEncoded}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	rdr, err := ce.client.ImagePush(ctx, destination, opts)
	if err != nil {
		return err
	}
	defer rdr.Close()
	err = Show(rdr)
	if err != nil {
		return err
	}
	return nil
}
