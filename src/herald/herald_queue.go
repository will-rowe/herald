package herald

import (
	"fmt"

	"github.com/will-rowe/herald/src/services"
)

// AnnounceSamples will processes the queues and submit service requests
func (herald *Herald) AnnounceSamples() error {

	// check service providers are available
	for _, service := range services.ServiceRegister {
		if err := service.CheckAccess(); err != nil {
			return err
		}
	}

	// process the experiments first
	if err := herald.processExperimentQueue(); err != nil {
		return err
	}

	return nil
}

// processExperimentQueue will run the service requests in the experiment queue
func (herald *Herald) processExperimentQueue() error {
	if herald.experimentQueue.Len() == 0 {
		return fmt.Errorf("experiment queue is empty")
	}

	// process the queue
	for herald.experimentQueue.Len() > 0 {

		// grab the experiment that is first in the queue
		e := herald.experimentQueue.Front()
		exp := e.Value.(*services.Experiment)

		// get the tags in order
		for _, tag := range exp.Metadata.GetRequestOrder() {

			// get the service details
			service := services.ServiceRegister[tag]

			// run the service request
			service.SendRequest(exp)

			// TODO: if the service is blocking, wait for it
		}

		fmt.Print(exp)

		// dequeue the experiment
		herald.experimentQueue.Remove(e)
	}

	return nil
}
