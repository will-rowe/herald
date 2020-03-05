package services

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/will-rowe/herald/src/helpers"
)

// SubmitSequencingProcess will send a sequencing request
func SubmitSequencingProcess(experiment *Experiment, service *Service) {

	// TODO: check if sequencing has already been done
	// this is just a temp hack to test out the call back
	if err := helpers.CheckDirExists(experiment.GetFast5OutputDirectory()); err == nil {

		// mark tag as complete
		experiment.Metadata.GetTags()[service.name] = true
		return
	}

	// instantiate a client connection, on the TCP port the server is bound to
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf(":%d", service.port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := NewSequenceClient(conn)

	// send the request and collect the response
	response, err := c.RunSequencing(context.Background(), &ProcessSubmission{Val1: "sequencing request"})
	if err != nil {
		log.Fatalf("Error when calling RunSequencing: %s", err)
	}

	// process the response
	log.Printf("Response from server: %s", response.Val2)
}
