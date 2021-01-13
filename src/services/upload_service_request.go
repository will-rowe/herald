package services

import (
	"context"
	"fmt"

	grpc "google.golang.org/grpc"

	"github.com/will-rowe/herald/src/records"
)

// SubmitUpload will request a data upload
// using the
func SubmitUpload(heraldRecord interface{}, service *Service) error {

	// assert we have a Sample, not a Run
	var run *records.Run
	switch heraldRecord.(type) {
	case *records.Sample:
		return fmt.Errorf("can't submit Sample in data upload request, need a Run")
	case *records.Run:
		run = heraldRecord.(*records.Run)
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
	conn, err := grpc.Dial(fmt.Sprintf(":%d", service.port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect to Herald server: %s", err)
	}
	defer conn.Close()
	c := NewHeraldClient(conn)

	// send the request and collect the response
	response, err := c.SubmitUpload(context.Background(), &UploadRequest{Val1: run.GetFastqOutputDirectory()})
	if err != nil {
		return fmt.Errorf("could not run upload: %s", err)
	}

	// process the response
	fmt.Printf("Response from server: %s", response.Val2)

	return fmt.Errorf("Response from server: %s", response.Val2)
}
