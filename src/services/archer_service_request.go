package services

import (
	"context"
	"fmt"

	archer "github.com/will-rowe/archer/pkg/api/v1"
	"google.golang.org/grpc"

	"github.com/will-rowe/herald/src/records"
)

// SubmitArcherUpload will request a data upload
// using the archer service.
func SubmitArcherUpload(heraldRecord interface{}, service *Service) error {

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

	// form an archer request
	var request *archer.ProcessRequest
	_ = run

	// connect to the gRPC server
	conn, err := grpc.Dial(service.GetAddress(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	// establish the client
	client := archer.NewArcherClient(conn)

	// send the request
	resp, err := client.Process(context.Background(), request)
	if err != nil {
		return err
	}

	return fmt.Errorf("upload request response: %s", resp)
}
