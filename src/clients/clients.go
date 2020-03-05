//go:generate protoc -I=protobuf/services --go_out=plugins=grpc:src/services protobuf/services/*.proto

// Package clients contains the functions that make requests to services and parses their responses
package clients

import (
	"fmt"

	"google.golang.org/grpc"
)

// TestServiceConnection will test if a gRPC connection can be made on a given port
func TestServiceConnection(portNum int) error {

	// instantiate a client connection, on the TCP port the server is bound to
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf(":%d", portNum), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect to port %d: %s", portNum, err)
	}
	conn.Close()
	return nil
}

//
func DummyProcess() {
	fmt.Println("dummy called")
}
