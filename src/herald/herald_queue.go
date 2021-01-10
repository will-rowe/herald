package herald

import (
	"fmt"

	"github.com/will-rowe/herald/src/services"
)

// AnnounceSamples will processes the queues and submit service requests
func (herald *Herald) AnnounceSamples() error {
	herald.Lock()
	defer herald.Unlock()
	if herald.announcementQueue.Len() == 0 {
		return fmt.Errorf("announcement queue is empty")
	}

	// check service providers are available
	for _, service := range services.ServiceRegister {
		if err := service.CheckAccess(); err != nil {
			return err
		}
	}

	// iterate once over the queue and process all the runs first
	for request := herald.announcementQueue.Front(); request != nil; request = request.Next() {
		switch v := request.Value.(type) {
		default:
			return fmt.Errorf("unexpected type in queue: %T", v)
		case *services.Sample:
			continue
		case *services.Run:

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

			// set the status to announced
			v.Metadata.SetStatus(services.Status_announced)

			// dequeue the sample
			herald.announcementQueue.Remove(request)
		}
	}

	// process the remaining queue (should just be samples now)
	for request := herald.announcementQueue.Front(); request != nil; request = request.Next() {

		// grab the sample that is first in the queue
		sample := request.Value.(*services.Sample)

		// get the tags in order
		for _, tag := range sample.Metadata.GetRequestOrder() {
			_ = tag
		}

		// evalute the sample

		// update fields and propogate to linked data

		// decide if it should be dequeued

		// update the status of the sample and dequeue it
		sample.Metadata.AddComment("sample announced.")
		sample.Metadata.SetStatus(services.Status_announced)
		if err := herald.updateRecord(sample); err != nil {
			return err
		}
		herald.announcementQueue.Remove(request)
	}

	if herald.announcementQueue.Len() != 0 {
		return fmt.Errorf("announcements sent but queue still contains %d requests", herald.announcementQueue.Len())
	}
	return nil
}
