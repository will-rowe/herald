//go:generate protoc -I=protobuf/services --go_out=plugins=grpc:src/services protobuf/services/*.proto

// Package clients contains the functions that make requests to services and parses their responses
package clients

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/will-rowe/herald/src/services"
)

// submitBasecallingProcess will send a basecalling request
func submitBasecallingProcess() {
	var conn *grpc.ClientConn
	//conn, err := grpc.Dial(fmt.Sprintf(":%d", services.TCPport), grpc.WithInsecure())
	conn, err := grpc.Dial(fmt.Sprintf(":%d", 7777), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := services.NewBasecallClient(conn)

	response, err := c.RunBasecalling(context.Background(), &services.ProcessSubmission2{Val1: "basecalling request"})
	if err != nil {
		log.Fatalf("Error when calling RunBasecalling: %s", err)
	}
	log.Printf("Response from server: %s", response.Val2)
}
