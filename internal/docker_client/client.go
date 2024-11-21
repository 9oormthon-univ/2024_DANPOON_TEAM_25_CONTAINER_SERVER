package dockerclient

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/docker/docker/api/types"
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

func (d *DockerClient) dockerLogin(ctx context.Context) error {
	// echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
	log.Println(d.Username)
	log.Println(d.DockerCredential)
	cmd := exec.CommandContext(ctx, "echo", d.DockerCredential, "docker", "login", "-u", d.Username, "--password-stdin")
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}
func (d *DockerClient) pushImage(ctx context.Context, imageName string, logCallback func(log string)) error {
	err := d.dockerLogin(ctx)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "docker", "push", imageName)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Start the docker push command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start docker push command: %w", err)
	}

	// Channels to receive stdout and stderr lines
	outputChan := make(chan string)
	errorChan := make(chan string)
	doneChan := make(chan struct{})

	// Function to read from a pipe and send lines to a channel
	readPipe := func(pipe io.ReadCloser, outChan chan<- string) {
		defer close(outChan)
		scanner := bufio.NewScanner(pipe)
		for scanner.Scan() {
			outChan <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			log.Printf("error reading pipe: %v", err)
		}
	}

	go readPipe(stdoutPipe, outputChan)
	go readPipe(stderrPipe, errorChan)

	go func() {
		cmd.Wait()
		close(doneChan)
	}()

	for {
		select {
		case line, ok := <-outputChan:
			if ok {
				processPushOutput(line, logCallback)
			}
		case line, ok := <-errorChan:
			if ok {
				processPushOutput(line, logCallback)
			}
		case <-doneChan:
			if ctx.Err() == context.DeadlineExceeded {
				return fmt.Errorf("docker push timed out")
			}
			if err := cmd.ProcessState.ExitCode(); err != 0 {
				return fmt.Errorf("docker push failed with exit code %d", cmd.ProcessState.ExitCode())
			}
			log.Printf("Image %s pushed successfully", imageName)
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func processPushOutput(line string, logCallback func(log string)) {
	var jsonLine map[string]interface{}
	if err := json.Unmarshal([]byte(line), &jsonLine); err == nil {
		if errMsg, ok := jsonLine["error"].(string); ok {
			logCallback(fmt.Sprintf("Error: %s", errMsg))
			return
		}
		if status, ok := jsonLine["status"].(string); ok {
			var logMsg string
			if progress, ok := jsonLine["progress"].(string); ok {
				logMsg = fmt.Sprintf("%s %s", status, progress)
			} else {
				logMsg = status
			}
			logCallback(logMsg)
			return
		}
	} else {
		logCallback(line)
	}
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
