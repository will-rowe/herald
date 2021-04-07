package services

import (
	"fmt"
	"net"
	"time"

	"github.com/will-rowe/herald/src/records"
)

// minknowService is an adapter to submit requests
// from Herald to a Minknow service.
type minknowService struct {
	name       string             // name of the service
	recordType records.RecordType // the type of Herald record this service operates on (sun or rample)
	dependsOn  []string           // the other services that should have completed prior to this one being contacted
	address    string             // the gRPC address of the service
	port       int                // the gRPC port the service is accepting requests on
}

// NewMinknowService will construct a
// new Minknow service adaptor and
// return the Service interface.
func NewMinknowService(name string, recordType records.RecordType, dependsOn []string, address string, port int) Service {
	as := &minknowService{
		name:       name,
		recordType: recordType,
		dependsOn:  dependsOn,
		address:    address,
		port:       port,
	}

	return as
}

// GetServiceName returns the name of the service.
func (m *minknowService) GetServiceName() string {
	return m.name
}

// GetRecordType returns the record type (sample/run)
// that the service accepts.
func (m *minknowService) GetRecordType() records.RecordType {
	return m.recordType
}

// GetAddress will return the address of
// the gRPC server running the service.
func (m *minknowService) GetAddress() string {
	return fmt.Sprintf("%v:%d", m.address, m.port)
}

// CheckAccess returns true if the service is
// accessible.
func (m *minknowService) CheckAccess() bool {
	conn, err := net.DialTimeout("tcp", m.GetAddress(), time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// GetDependencies will return a slice
// of the dependency names.
func (m *minknowService) GetDependencies() []string {
	return m.dependsOn
}

// SendRequest will establish a minknow client, formulate
// a request and submit it to the running service.
func (m *minknowService) SendRequest(record interface{}) error {

	// TODO: this is yet to be implemented
	// Minknow is included as a service template.
	// refer to the Archer service for a complete example.

	return nil
}
