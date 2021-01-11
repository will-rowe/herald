// Package server is the gRPC server implementation used in the background of a Herald instance and manages the service requests.
package server

import (
	"fmt"
	"log"
	"net"
	"os"

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
// for connection and handles
// requests.
// It will write to a log at the
// provided filepath.
func Start(logFilePath string) {

	// open the log
	// TODO: look at rolling log: https://github.com/natefinch/lumberjack
	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("starting herald server")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", TCPport))
	if err != nil {
		log.Printf("server failed to listen: %v", err)
	}

	// create a server instance
	heraldServer := HeraldServer{}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the service
	services.RegisterHeraldServer(grpcServer, &heraldServer)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("server failed to serve: %s", err)
	}

}
