// Package server is the gRPC server implementation used in the background of a Herald instance and manages the service requests.
package server

import (
	"fmt"
	"log"
	"net"

	"github.com/will-rowe/herald/src/services"
	grpc "google.golang.org/grpc"
)

// TCPport is the port on which the Herald
// service is running.
const TCPport = 7779

// HeraldServer represents the gRPC server
type HeraldServer struct {
}

// Start will start a gRPC server,
// attach the service and then waits
// for connection and handles requests/
func Start() {

	log.Println("starting herald server")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", TCPport))
	if err != nil {
		log.Printf("failed to listen: %v", err)
	}

	// create a server instance
	heraldServer := HeraldServer{}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the service
	services.RegisterUploadServer(grpcServer, &heraldServer)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("failed to serve: %s", err)
	}

}
