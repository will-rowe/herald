package clients

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/will-rowe/herald/src/services"
)

// SubmitSequencingProcess will send a sequencing request
func SubmitSequencingProcess() {
	// instantiate a client connection, on the TCP port the server is bound to
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf(":%d", 7777), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := services.NewSequenceClient(conn)

	// send the request and collect the response
	response, err := c.RunSequencing(context.Background(), &services.ProcessSubmission{Val1: "sequencing request"})
	if err != nil {
		log.Fatalf("Error when calling RunSequencing: %s", err)
	}

	// process the response
	log.Printf("Response from server: %s", response.Val2)
}
