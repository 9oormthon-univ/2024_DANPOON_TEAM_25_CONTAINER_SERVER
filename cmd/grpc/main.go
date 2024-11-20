package main

import (
	"context"
	"log"
	"net"

	pb "github.com/9oormthon-univ/2024_DANPOON_TEAM_25_CONTAINER_SERVER/proto/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedCourseIDEServiceServer
}

func (s *server) Create(ctx context.Context, req *pb.CourseIDECreateRequest) (*pb.CourseIDECreateResponse, error) {
	//TODO: Create Docker Image With Nix env and flake.
	return &pb.CourseIDECreateResponse{Ok: true}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Fail to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCourseIDEServiceServer(grpcServer, &server{})

	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Fail to serve: %v", err)
	}
}
