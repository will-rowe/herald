package services

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	c := NewBasecallClient(conn)

	response, err := c.RunBasecalling(context.Background(), &ProcessSubmission2{Val1: "basecalling request"})
	if err != nil {
		log.Fatalf("Error when calling RunBasecalling: %s", err)
	}
	log.Printf("Response from server: %s", response.Val2)
}
