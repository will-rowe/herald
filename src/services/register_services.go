package services

import (
	"fmt"

	"github.com/will-rowe/herald/src/records"
)

// DefaultAddress
var DefaultAddress string = "127.0.0.1"

// DefaultTCPport
var DefaultTCPport int = 60742

// DefaultArcherVersion is the API version to use for Archer
var DefaultArcherVersion string = "1"

// init is where we create the services
// that we want to offer at runtime.
// TODO: add more docs on this...
func init() {

	// init the process checker
	ServiceRegister = make(map[string]Service)

	///
	// create the service definitions
	//
	archerService := NewArcherService("Archer upload", records.RecordType_run, nil, DefaultAddress, DefaultTCPport)

	// check and register the services
	checkAndRegister(archerService)
}

// Service is an interface that allows Herald to submit requests to a service.
type Service interface {
	GetServiceName() string               // returns the name of the service
	GetRecordType() records.RecordType    // returns the record type (sample/run) that the service accepts
	GetAddress() string                   // returns the address of the server offering the service
	CheckAccess() bool                    // returns true if the service is accessible
	GetDependencies() []string            // retuns service dependencies as a slice of service names
	SendRequest(record interface{}) error // function to establish a client and submit the service request
}

// ServiceRegister is used to register all the
// available services to the current Herald
// runtime.
var ServiceRegister map[string]Service

// checkAndRegister will perform a few sanity checks
// and then register a service.
func checkAndRegister(service Service) {

	// check service name isn't taken
	if name, ok := ServiceRegister[service.GetServiceName()]; ok {
		panic(fmt.Sprintf("service name already exists: %v", name))
	}

	// check the record type is either sample or run
	switch service.GetRecordType() {
	case records.RecordType_run:
		break
	case records.RecordType_sample:
		break
	default:
		panic(fmt.Sprintf("unsupported record type: %v", service.GetRecordType()))
	}

	// check the dependencies
	deps := service.GetDependencies()
	if len(deps) != 0 {
		for _, depName := range deps {

			// can't depend on itself
			if depName == service.GetServiceName() {
				panic("process can't depend on itself")
			}

			// dependency must already be registered
			dependency, ok := ServiceRegister[depName]
			if !ok {
				panic(fmt.Sprintf("service dependency not registered, make sure to register %v first", depName))
			}
			// TODO: resolve the dependencies
			_ = dependency
		}
	}

	// register the service
	ServiceRegister[service.GetServiceName()] = service
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
