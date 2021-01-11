package server

import (
	"context"
	"fmt"

	grpc "google.golang.org/grpc"

	"github.com/will-rowe/herald/src/services"
)

// SubmitMinionPipeline will send a pipeline request
func SubmitMinionPipeline(heraldRecord interface{}, service *Service) error {

	// assert we have a Sample, not a Run
	var sample *services.Sample
	switch heraldRecord.(type) {
	case *services.Run:
		return fmt.Errorf("can't submit Run in pipeline request, need a Sample")
	case *services.Sample:
		sample = heraldRecord.(*services.Sample)
	default:
		return fmt.Errorf("unsupported Herald record type")
	}

	// TODO: check if pipeline has already been submitted
	//if err := pipelinealreadysubmitted(); err == nil {

	// mark tag as complete
	//	run.Metadata.GetTags()[service.name] = true
	//	return nil
	//}

	// instantiate a client connection, on the TCP port the server is bound to
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf(":%d", TCPport), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect to Herald server: %s", err)
	}
	defer conn.Close()
	c := services.NewUploadClient(conn)

	// form the request
	request := &services.MinionPipelineRequest{
		Pipeline: sample.GetMetadata().Label,
		Param1:   "here you go",
	}

	// send the request and collect the response
	response, err := c.SubmitMinionPipeline(context.Background(), request)
	if err != nil {
		return fmt.Errorf("could not run upload: %s", err)
	}

	// process the response
	//fmt.Printf("Response from server: %s", response.Val)

	return fmt.Errorf("Response from server: %s", response.Val)
}
