package main

import (
	"log"
	"net"

	dockerclient "github.com/9oormthon-univ/2024_DANPOON_TEAM_25_CONTAINER_SERVER/internal/docker_client"
	pb "github.com/9oormthon-univ/2024_DANPOON_TEAM_25_CONTAINER_SERVER/proto/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedCourseIDEServiceServer
	Client *dockerclient.DockerClient
}

func NewServer() (*server, error) {
	dockerClient, err := dockerclient.NewDockerClient()
	if err != nil {
		return nil, err
	}
	return &server{Client: dockerClient}, nil
}

func (s *server) Create(req *pb.CourseIDECreateRequest, stream pb.CourseIDEService_CreateServer) error {
	err := s.Client.CreateImage("test", req.Spec, func(logMessage string) {
		if err := stream.Send(&pb.CourseIDECreateResponse{Message: logMessage, Ok: false}); err != nil {
			log.Printf("Fail to stream log: %v", err)
		}
	})
	if err != nil {
		return err
	}
	return nil
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

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Fail to serve: %v", err)
	}
}
