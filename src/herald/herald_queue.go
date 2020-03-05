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

	// iterate once over the queue and process all the experiments that need sequencing and basecalling
	for e := herald.announcementQueue.Front(); e != nil; e = e.Next() {

		switch v := e.Value.(type) {
		default:
			return fmt.Errorf("unexpected type in queue: %T", v)
		case *services.Sample:
			continue
		case *services.Experiment:

			// get the tags in order
			for _, tag := range v.Metadata.GetRequestOrder() {

				// check it's not been completed already
				if complete := v.Metadata.GetTags()[tag]; complete {
					continue
				}

				// get the service details
				service := services.ServiceRegister[tag]

				// run the service request
				if err := service.SendRequest(v); err != nil {
					return err
				}
			}

			// dequeue the sample
			herald.announcementQueue.Remove(e)
		}
	}

	// process the remaining queue (should just be samples now)
	for herald.announcementQueue.Len() > 0 {

		// grab the sample that is first in the queue
		s := herald.announcementQueue.Front()
		sample := s.Value.(*services.Sample)

		// get the tags in order
		for _, tag := range sample.Metadata.GetRequestOrder() {
			_ = tag
		}

		// evalute the sample

		// update fields and propogate to linked data

		// decide if it should be dequeued

		// dequeue the sample
		herald.announcementQueue.Remove(s)
	}

	return nil
}
