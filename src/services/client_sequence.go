package services

import (
	"fmt"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/will-rowe/herald/src/helpers"
)

// SubmitSequencingProcess will send a sequencing request
func SubmitSequencingProcess(experiment *Experiment, service *Service) error {

	// TODO: check if sequencing has already been done
	// this is just a temp hack to test out the call back
	if err := helpers.CheckDirExists(experiment.GetFast5OutputDirectory()); err == nil {

		// mark tag as complete
		experiment.Metadata.GetTags()[service.name] = true
		return nil
	}

	// instantiate a client connection, on the TCP port the server is bound to
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf(":%d", service.port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect: %s", err)
	}
	defer conn.Close()
	c := NewSequenceClient(conn)

	// send the request and collect the response
	response, err := c.RunSequencing(context.Background(), &ProcessSubmission{Val1: "sequencing request"})
	if err != nil {
		return fmt.Errorf("Error when calling RunSequencing: %s", err)
	}

	// process the response
	//log.Printf("Response from server: %s", response.Val2)
	_ = response
	return nil
}
