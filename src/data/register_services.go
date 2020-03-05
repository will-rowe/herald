// Package data adds some wrapper functions to the protobuf messages
package data

import (
	"fmt"

	"github.com/will-rowe/herald/src/clients"
)

// Service
type Service struct {
	name      string
	dependsOn []string
	port      int
	callback  func()
}

// ServiceRegister is used to register all the available processes
var ServiceRegister map[string]*Service

// init the process definitions at runtime
func init() {

	// init the process checker
	ServiceRegister = make(map[string]*Service)

	// create the process definitions
	createServiceDefinition("sequence", nil, 7777, clients.DummyProcess)
	createServiceDefinition("basecall", nil, 7778, clients.DummyProcess)
	//createServiceDefinition("rampart", nil)
	//createServiceDefinition("pipelineA", []string{"sequence", "basecall"})

}

// createServiceDefinition will init a process
func createServiceDefinition(sName string, sDependsOn []string, sPort int, sFunc func()) {

	// check the process does not already exist
	if _, exists := ServiceRegister[sName]; exists {
		panic(fmt.Sprintf("process already exists: %v", sName))
	}

	// init the process
	newService := &Service{
		name:      sName,
		dependsOn: []string{},
		port:      sPort,
		callback:  sFunc,
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
