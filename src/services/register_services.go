package services

import (
	"fmt"

	"google.golang.org/grpc"

	"github.com/will-rowe/herald/src/records"
)

// DefaultAddress
var DefaultAddress string = ""

// DefaultTCPport
var DefaultTCPport int = 50051

// init is where we create the services
// that we want to offer at runtime.
// TODO: add more docs on this...
func init() {

	// init the process checker
	ServiceRegister = make(map[string]*Service)

	///
	// create the service definitions
	//  - define the request callback function for your service (e.g. SubmitArcherUpload)
	//  - call the registerService function to tell herald about the service and what record type to use
	// Run services
	registerService(records.RecordType_run, "Upload (archer)", nil, DefaultAddress, DefaultTCPport, SubmitArcherUpload)
	//
	// Sample services
	//registerService(records.RecordType_sample, "ARTIC pipeline (medaka)", nil, DefaultTCPport, SubmitMinionPipeline)
	//registerService(RecordType_sample, "pipelineA", []string{"sequence", "basecall", "upload"}, 7780, SubmitSequencingProcess)
}

// Service is holds the information needed by Herald to send messages to a service provider
type Service struct {
	recordType      records.RecordType                                     // the type of Herald record this service operates on (sun or rample)
	name            string                                                 // name of the service
	dependsOn       []string                                               // the other services that should have completed prior to this one being contacted
	address         string                                                 // the gRPC address of the service
	port            int                                                    // the gRPC port the service is accepting requests on
	requestCallback func(heraldRecord interface{}, service *Service) error // the function to run when contacting the service
}

// ServiceRegister is used to register
// all the available services to the
// current Herald runtime.
var ServiceRegister map[string]*Service

// GetRecordType will return what type of
// Herald record the service will run on.
func (service *Service) GetRecordType() string {
	return service.recordType.String()
}

// GetAddress will return the address of
// the gRPC server running the service.
func (service *Service) GetAddress() string {
	return fmt.Sprintf("%v:%d", service.address, service.port)
}

// CheckAccess will return true if the
// gRPC service can be accessed.
func (service *Service) CheckAccess() bool {

	// instantiate a client connection, on the TCP port the server is bound to
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(service.GetAddress(), grpc.WithInsecure())
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// GetDeps will return a slice
// of the dependency names.
func (service *Service) GetDeps() []string {
	return service.dependsOn
}

// SendRequest will run the
// service callback function.
func (service *Service) SendRequest(heraldRecord interface{}) error {

	// check the access
	if !service.CheckAccess() {
		return fmt.Errorf("can't access %s service on %s", service.name, service.GetAddress())
	}

	// run the defined request
	return service.requestCallback(heraldRecord, service)
}

// registerService will register a service to the
// Herald runtime.
func registerService(recordType records.RecordType, sName string, sDependsOn []string, sAddress string, sPort int, sFunc func(heraldRecord interface{}, service *Service) error) {

	// check the record type is either sample or run
	switch recordType {
	case records.RecordType_run:
		break
	case records.RecordType_sample:
		break
	default:
		panic(fmt.Sprintf("unsupported record type: %v", recordType))
	}

	// check the process does not already exist
	if _, exists := ServiceRegister[sName]; exists {
		panic(fmt.Sprintf("process already exists: %v", sName))
	}

	// init the process
	newService := &Service{
		recordType:      recordType,
		name:            sName,
		dependsOn:       []string{},
		address:         sAddress,
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

/*

TODO: moved this here for now
we need to arrange services by a priority DAG
this was done during tagging but should be done here instead to remove cyclic imports throughout codebase

// createServiceDAG creates a linear ordering of services that accounts for service dependencies
func (heraldData *HeraldData) createServiceDAG() error {

	// reset requestOrder
	heraldData.RequestOrder = []string{}

	// transfer service names from map to slice
	serviceList := make([]string, len(heraldData.Tags))
	numServices := 0
	for serviceName := range heraldData.GetTags() {
		serviceList[numServices] = serviceName
		numServices++
	}

	// create a dag
	dag := toposort.NewGraph(numServices)

	// create the nodes
	dag.AddNodes(serviceList...)

	// loop through input list and create edges
	for _, serviceName := range serviceList {

		// ignore services with no dependencies
		service := ServiceRegister[serviceName]
		if len(service.GetDeps()) == 0 {
			continue
		}

		// loop over the depencies and draw edges
		for _, dependencyName := range service.GetDeps() {
			dag.AddEdge(dependencyName, serviceName)
		}
	}

	// sort the graph and check for cycles
	result, ok := dag.Toposort()
	if !ok {
		return fmt.Errorf("service dependency cycle detected")
	}
	heraldData.RequestOrder = result
	return nil
}



*/
