package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

type MalformedErrorDetails struct {
	messge string
}

func (err MalformedErrorDetails) Error() string {
	return err.messge
}

func NewMalformedErrorDetails(msg string) *MalformedErrorDetails {
	return &MalformedErrorDetails{
		messge: msg,
	}
}

type ErrorDetails struct {
	messge string
}

func (err ErrorDetails) Error() string {
	return err.messge
}

func NewErrorDetails(msg string) *ErrorDetails {
	return &ErrorDetails{
		messge: msg,
	}
}

type CEClient struct {
	client DockerClient
}

func NewCEClient() (*CEClient, error) {
	container_cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("error while creating Docker client (%w)", err)
	}
	return &CEClient{
		client: container_cli,
	}, nil
}

func (ce CEClient) Build(filepath string, name string, tags []string) error {
	image_tag := name + ":" + tags[0]
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()
	tar, err := archive.TarWithOptions(filepath, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("error archiving %s (%w)", filepath, err)
	}
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{image_tag},
	}
	res, err := ce.client.ImageBuild(ctx, tar, opts)
	if err != nil {
		return fmt.Errorf("error building %s (%w)", name, err)
	}
	if res.Body != nil {
		err := Show(res.Body, os.Stdout)
		if err != nil {
			return fmt.Errorf("error for %s found by container engine during build (%w)", name, err)
		}
		err = res.Body.Close()
		if err != nil {
			return fmt.Errorf("error closing image build response (%w)", err)
		}
	}
	return nil
}

func Show(rd io.Reader, writer io.Writer) error {
	var lastLine string
	var nextLine string
	scanner := bufio.NewScanner(rd)
	line := &StreamLine{}

	for scanner.Scan() {
		lastLine = scanner.Text()
		nextLine = scanner.Text()
		err := json.Unmarshal([]byte(nextLine), &line)
		if err != nil {
			return fmt.Errorf("error unmarshalling jsons stream %s (%w)", lastLine, err)
		}
		if _, err := writer.Write([]byte(line.Stream)); err != nil {
			return fmt.Errorf("error writing json stream (%w)", err)
		}
		line = &StreamLine{}
	}

	errLine := &ErrorLine{}
	err := json.Unmarshal([]byte(lastLine), errLine)
	if err != nil {
		return NewMalformedErrorDetails(
			fmt.Sprintf(
				"error unmarshalling error details from jsons stream producer  %s (%v)",
				lastLine, err))
	}

	if errLine.Error != "" {
		return NewErrorDetails(
			fmt.Sprintf("error details from jsons stream producer (%s)", errLine.Error))
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning jsons stream (%w)", err)
	}

	return nil
}

func (ce CEClient) Tag(image_tag string, destination string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	err := ce.client.ImageTag(ctx, image_tag, destination)

	if err != nil {
		return fmt.Errorf("error tagging %s (%w)", destination, err)
	}
	return nil
}

func (ce CEClient) Push(destination string, username string, password string, registry_address string) error {
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
		return fmt.Errorf("error pushing %s (%w)", destination, err)
	}
	if rdr != nil {
		err := Show(rdr, os.Stdout)
		if err != nil {
			return fmt.Errorf("error for %s found by container engine during push (%w)", destination, err)
		}
		err = rdr.Close()
		if err != nil {
			return fmt.Errorf("error closing image push reader (%w)", err)
		}
	}
	return nil
}
