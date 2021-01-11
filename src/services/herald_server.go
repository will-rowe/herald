// Package services is the gRPC server implementation used in the background of a Herald instance and manages the service requests.
package services

import (
	"fmt"
	"log"
	"net"
	"os"

	grpc "google.golang.org/grpc"
)

// DefaultTCPport is the port on
// which the Herald gRPC is running
// services.
const DefaultTCPport = 7779

// Server represents the gRPC server
// that is used by Herald to orchestrate
// services during runtime.
type Server struct {
	logFile *os.File
	tcpPort int
}

// Option is a wrapper struct used to pass functional
// options to the HeraldServer constructor.
type Option func(x *Server) error

// SetLog will set the log file for a server
// and open it for writing
func SetLog(logFile string) Option {
	return func(x *Server) error {
		return x.setLog(logFile)
	}
}

// SetPort will set the TCP port for a
// server.
func SetPort(port int) Option {
	return func(x *Server) error {
		return x.setPort(port)
	}
}

// Start the HeraldServer.
func (x *Server) Start() {
	log.Println("starting herald server")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", x.tcpPort))
	if err != nil {
		log.Printf("server failed to listen: %v", err)
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach the service
	RegisterHeraldServer(grpcServer, x)

	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("server failed to serve: %s", err)
	}
}

// Stop the HeraldServer.
func (x *Server) Stop() {
	log.Println("stopping herald server")
	x.logFile.Close()
}

// NewServer will return a HeraldServer.
func NewServer(options ...Option) (*Server, error) {

	// create the server
	x := &Server{
		tcpPort: DefaultTCPport,
	}

	// add the user provided options
	for _, option := range options {
		err := option(x)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

// setLog.
func (x *Server) setLog(logFile string) error {

	// TODO: look at rolling log: https://github.com/natefinch/lumberjack
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	log.SetOutput(f)
	x.logFile = f
	return nil
}

// setPort.
func (x *Server) setPort(port int) error {
	x.tcpPort = port
	return nil
}
