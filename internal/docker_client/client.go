package dockerclient

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	Client *client.Client
}

func NewDockerClient() (*DockerClient, error) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerClient{Client: client}, nil
}

func (d *DockerClient) CreateImage(imageName string, specs []string, logCallback func(log string)) error {
	dockerfileSpec := ""
	for _, spec := range specs {
		dockerfileSpec += fmt.Sprintf("nixpkgs.%s ", spec)
	}
	dockerfile := fmt.Sprintf(CODESERVER_DOCKERFILE, dockerfileSpec)
	err := d.buildImage(imageName, dockerfile, logCallback)
	if err != nil {
		return err
	}
	return nil
}

func (d *DockerClient) buildImage(name string, dockerfile string, logCallback func(log string)) error {
	buildContext, err := createDockerfileContext(dockerfile)
	if err != nil {
		return err
	}
	buildOptions := types.ImageBuildOptions{
		Tags:       []string{name},
		Dockerfile: "Dockerfile",
		Remove:     true,
	}
	ctx := context.Background()
	buildResponse, err := d.Client.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return err
	}
	defer buildResponse.Body.Close()
	decoder := json.NewDecoder(buildResponse.Body)
	for {
		var line map[string]interface{}
		if err := decoder.Decode(&line); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode build log: %w", err)
		}

		if stream, ok := line["stream"].(string); ok {
			logCallback(stream)
		}
	}
	log.Printf("Image %s build successfully", name)
	return nil
}

func createDockerfileContext(dockerfile string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	if err := addFileToTar(tw, "Dockerfile", dockerfile); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}

func addFileToTar(tw *tar.Writer, name, dockerfile string) error {
	header := &tar.Header{
		Name: name,
		Mode: 0644,
		Size: int64(len(dockerfile)),
	}
	if err := tw.WriteHeader(header); err != nil {
		return err
	}
	if _, err := tw.Write([]byte(dockerfile)); err != nil {
		return err
	}
	return nil
}
