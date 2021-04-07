package services

import (
	"context"
	"fmt"
	"net"
	"time"

	archer "github.com/will-rowe/archer/pkg/api/v1"
	"google.golang.org/grpc"

	"github.com/will-rowe/herald/src/records"
)

// archerService is an adapter to submit requests
// from Herald to an archer service.
type archerService struct {
	name       string             // name of the service
	recordType records.RecordType // the type of Herald record this service operates on (sun or rample)
	dependsOn  []string           // the other services that should have completed prior to this one being contacted
	address    string             // the gRPC address of the service
	port       int                // the gRPC port the service is accepting requests on
}

// NewArcherService will construct a
// new Archer service adaptor and
// return the Service interface.
func NewArcherService(name string, recordType records.RecordType, dependsOn []string, address string, port int) Service {
	as := &archerService{
		name:       name,
		recordType: recordType,
		dependsOn:  dependsOn,
		address:    address,
		port:       port,
	}

	return as
}

// GetServiceName returns the name of the service.
func (a *archerService) GetServiceName() string {
	return a.name
}

// GetRecordType returns the record type (sample/run)
// that the service accepts.
func (a *archerService) GetRecordType() records.RecordType {
	return a.recordType
}

// GetAddress will return the address of
// the gRPC server running the service.
func (a *archerService) GetAddress() string {
	return fmt.Sprintf("%v:%d", a.address, a.port)
}

// CheckAccess returns true if the service is
// accessible.
func (a *archerService) CheckAccess() bool {

	conn, err := net.DialTimeout("tcp", a.GetAddress(), time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// GetDependencies will return a slice
// of the dependency names.
func (a *archerService) GetDependencies() []string {
	return a.dependsOn
}

// SendRequest will establish an archer client, formulat
// a request and submit it to the running service.
func (a *archerService) SendRequest(record interface{}) error {

	// assert we have a Sample, not a Run
	var run *records.Run
	switch record.(type) {
	case *records.Sample:
		return fmt.Errorf("can't submit Sample in data upload request, need a Run")
	case *records.Run:
		run = record.(*records.Run)
	default:
		return fmt.Errorf("unsupported Herald record type")
	}

	// form an archer request
	var request *archer.ProcessRequest
	_ = run

	// connect to the gRPC server
	conn, err := grpc.Dial(a.GetAddress(), grpc.WithInsecure())
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
