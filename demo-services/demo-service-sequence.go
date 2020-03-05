//go:generate protoc -I=protobuf/services --go_out=plugins=grpc:src/services protobuf/services/*.proto

// Package main is used to run a demo service
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/will-rowe/herald/src/services"
	grpc "google.golang.org/grpc"
)

// TCPport is the port on which the service is running
const TCPport = 7777

// Server represents the gRPC server
type Server struct {
}

// RunSequencing receives a Sequence request, checks if it has a label and then returns a response
func (s *Server) RunSequencing(ctx context.Context, submission *services.ProcessSubmission) (*services.ProcessSummary, error) {
	log.Printf("received sequencing request")
	if submission.GetVal1() == "" {
		return nil, fmt.Errorf("no label sent")
	}
	log.Printf("experiment label: %s", submission.GetVal1())
	return &services.ProcessSummary{Val2: fmt.Sprintf("here is a response for the client regarding experiment label: %v", submission.GetVal1())}, nil
}

// StartServices will start a gRPC server, attach the service and then waits for connection
func main() {
	log.Println("starting the demo sequence service")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", TCPport))
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}

	// create a server instance
	s := Server{}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the service
	services.RegisterSequenceServer(grpcServer, &s)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("failed to serve: %s", err)
	}

}
