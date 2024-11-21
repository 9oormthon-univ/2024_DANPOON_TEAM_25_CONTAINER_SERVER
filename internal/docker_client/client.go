package dockerclient

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type AuthConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DockerClient struct {
	Client           *client.Client
	DockerCredential string
	Username         string
}

func NewDockerClient() (*DockerClient, error) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	username := os.Getenv("DOCKER_USERNAME")
	authConfig := AuthConfig{
		Username: username,
		Password: os.Getenv("DOCKER_CREDENTIAL"),
	}
	authBytes, err := json.Marshal(authConfig)
	if err != nil {
		return nil, err
	}
	authEncoded := base64.StdEncoding.EncodeToString(authBytes)
	return &DockerClient{Client: client, DockerCredential: authEncoded, Username: username}, nil
}

func (d *DockerClient) CreateImage(ctx context.Context, imageTag string, specs []string, logCallback func(log string)) error {
	imageName := fmt.Sprintf("%s/ide:%s", d.Username, imageTag)
	log.Println(imageName)
	dockerfileSpec := ""
	for _, spec := range specs {
		dockerfileSpec += fmt.Sprintf("nixpkgs.%s ", spec)
	}
	dockerfile := fmt.Sprintf(CODESERVER_DOCKERFILE, dockerfileSpec)
	err := d.buildImage(ctx, imageName, dockerfile, logCallback)
	if err != nil {
		return err
	}
	err = d.pushImage(ctx, imageName, logCallback)
	if err != nil {
		return err
	}
	return nil
}

func (d *DockerClient) buildImage(ctx context.Context, imageName string, dockerfile string, logCallback func(log string)) error {
	buildContext, err := createDockerfileContext(dockerfile)
	if err != nil {
		return err
	}
	buildOptions := types.ImageBuildOptions{
		Tags:       []string{imageName},
		Dockerfile: "Dockerfile",
		Remove:     true,
	}
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
	log.Printf("Image %s build successfully", imageName)
	return nil
}

func (d *DockerClient) pushImage(ctx context.Context, imageName string, logCallback func(log string)) error {
	pushResponse, err := d.Client.ImagePush(ctx, imageName, image.PushOptions{
		RegistryAuth: d.DockerCredential,
	})
	if err != nil {
		return fmt.Errorf("failed to push image: %w", err)
	}
	defer pushResponse.Close()

	decoder := json.NewDecoder(pushResponse)
	for {
		var line map[string]interface{}
		if err := decoder.Decode(&line); err != nil {

			if err == io.EOF {
				break
			}

			fmt.Printf("error decoding push response: %v\n", err)
		}

		if errMsg, ok := line["error"].(string); ok {
			fmt.Printf("error decoding push response: %v\n", errMsg)
		}

		if status, ok := line["status"].(string); ok {
			logCallback(status)
		}
	}

	log.Printf("Image %s pushed successfully", imageName)
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
