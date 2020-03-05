//go:generate protoc -I=protobuf --go_out=plugins=grpc:src/services protobuf/*.proto

// Package services contains the wrappers for the protobuf messages and the client request callbacks
package services

import (
	"fmt"

	"google.golang.org/grpc"
)

// Service is holds the information needed by Herald to send messages to a service provider
type Service struct {
	name            string                                               // name of the service
	dependsOn       []string                                             // the other services that should have completed prior to this one being contacted
	port            int                                                  // the port the service is accepting requests on
	requestCallback func(experiment *Experiment, service *Service) error // the function to run when contacting the service
}

// ServiceRegister is used to register all the available processes
var ServiceRegister map[string]*Service

// init the process definitions at runtime
func init() {

	// init the process checker
	ServiceRegister = make(map[string]*Service)

	///
	// create the process definitions
	//
	registerService("sequence", nil, 7777, SubmitSequencingProcess)
	registerService("basecall", nil, 7778, SubmitSequencingProcess)
	registerService("pipelineA", []string{"sequence", "basecall"}, 7779, SubmitSequencingProcess)
	//
	//
	//
}

// registerService will init a process
func registerService(sName string, sDependsOn []string, sPort int, sFunc func(experiment *Experiment, service *Service) error) {

	// check the process does not already exist
	if _, exists := ServiceRegister[sName]; exists {
		panic(fmt.Sprintf("process already exists: %v", sName))
	}

	// init the process
	newService := &Service{
		name:            sName,
		dependsOn:       []string{},
		port:            sPort,
		requestCallback: sFunc,
	}

	// check the dependencies
	if len(sDependsOn) != 0 {
		for _, depName := range sDependsOn {

			// can't depend on itself
			if depName == sName {
				panic("process can't depend on itself")
			}

			// dependency must already be registered
			dependency, ok := ServiceRegister[depName]
			if !ok {
				panic(fmt.Sprintf("process dependency not registered: %v", depName))
			}

			// TODO: this needs proper checking and probably lends itself to some sort of DAG scenario
			// then we can check for infinite loops etc.

			// add the dependency name
			newService.dependsOn = append(newService.dependsOn, dependency.name)
		}
	}

	// register the process
	ServiceRegister[sName] = newService
	return
}

// CheckAccess will check to see if the service port can be accessed
func (service *Service) CheckAccess() error {
	// instantiate a client connection, on the TCP port the server is bound to
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(fmt.Sprintf(":%d", service.port), grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("did not connect to port %d: %s", service.port, err)
	}
	conn.Close()
	return nil
}

// GetDeps will return a slice of the dependency names
func (service *Service) GetDeps() []string {
	return service.dependsOn
}

// SendRequest will run the service callback function
func (service *Service) SendRequest(experiment *Experiment) error {
	return service.requestCallback(experiment, service)
}
