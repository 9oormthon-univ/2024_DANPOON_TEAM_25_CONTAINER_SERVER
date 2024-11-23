package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"time"

	dockerclient "github.com/9oormthon-univ/2024_DANPOON_TEAM_25_CONTAINER_SERVER/internal/docker_client"
	gitclient "github.com/9oormthon-univ/2024_DANPOON_TEAM_25_CONTAINER_SERVER/internal/git_client"
	pb "github.com/9oormthon-univ/2024_DANPOON_TEAM_25_CONTAINER_SERVER/proto/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedCourseIDEServiceServer
	DockerClient *dockerclient.DockerClient
	GitClient    *gitclient.Gitclient
}

func NewServer() (*server, error) {
	dockerClient, err := dockerclient.NewDockerClient()
	if err != nil {
		return nil, err
	}
	gitClient := gitclient.NewGitClient()
	return &server{DockerClient: dockerClient, GitClient: gitClient}, nil
}

func (s *server) CreateImage(req *pb.CourseIDECreateRequest, stream pb.CourseIDEService_CreateImageServer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	imageTag := fmt.Sprintf("course%s", req.CourseId)
	encodedTag := base64.StdEncoding.EncodeToString([]byte(imageTag))
	err := s.DockerClient.CreateImage(ctx, encodedTag, req.Spec, func(logMessage string) {
		if err := stream.Send(&pb.CourseIDECreateResponse{Message: logMessage, Ok: false}); err != nil {
		}
	})
	if err != nil {
		log.Printf("Build Image error: %v", err)
		return err
	}
	return nil
}

func (s *server) CreatePod(ctx context.Context, req *pb.PodCreateRequest) (*pb.PodCreateResponse, error) {
	err := s.GitClient.ModifyRepository(req.CourseId, req.StudentId)
	if err != nil {
		return nil, err
	}
	return &pb.PodCreateResponse{Ok: true}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fail to listen: %v", err)
	}
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Fail to create server: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCourseIDEServiceServer(grpcServer, server)

	reflection.Register(grpcServer)
	log.Println("Server running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Fail to serve: %v", err)
	}
}
