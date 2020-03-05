package herald

import (
	"fmt"

	"github.com/will-rowe/herald/src/services"
)

// AnnounceSamples will processes the queues and submit service requests
func (herald *Herald) AnnounceSamples() error {

	if herald.announcementQueue.Len() == 0 {
		return fmt.Errorf("announcement queue is empty")
	}

	// check service providers are available
	for _, service := range services.ServiceRegister {
		if err := service.CheckAccess(); err != nil {
			return err
		}
	}

	// process the queue
	for herald.announcementQueue.Len() > 0 {

		// grab the sample that is first in the queue
		s := herald.announcementQueue.Front()
		sample := s.Value.(*services.Sample)

		// get the tags in order
		for _, tag := range sample.Metadata.GetRequestOrder() {

			// get the service details
			service := services.ServiceRegister[tag]

			if tag == "sequencing" {
				exp, err := herald.store.GetExperiment(sample.GetParentExperiment())
				if err != nil {
					return err
				}

				// run the service request
				service.SendRequest(exp)

			}

			// run the service request
			//service.SendRequest(sample)

			// TODO: if the service is blocking, wait for it
		}

		fmt.Print(sample)

		// evalute the sample

		// update fields and propogate to linked data

		// decide if it should be dequeued

		// dequeue the sample
		herald.announcementQueue.Remove(s)
	}

	return nil
}
