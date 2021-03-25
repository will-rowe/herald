package herald

import (
	"fmt"

	"github.com/will-rowe/herald/src/records"
	"github.com/will-rowe/herald/src/services"
)

// AnnounceSamples will processes the queues and submit service requests
func (herald *Herald) AnnounceSamples() error {
	herald.Lock()
	defer herald.Unlock()
	if herald.announcementQueue.Len() == 0 {
		return fmt.Errorf("announcement queue is empty")
	}

	// iterate once over the queue and process all the runs first
	for request := herald.announcementQueue.Front(); request != nil; request = request.Next() {
		switch v := request.Value.(type) {
		default:
			return fmt.Errorf("unexpected type in queue: %T", v)
		case *records.Sample:
			continue
		case *records.Run:
			// make the service requests
			for tag, complete := range v.Metadata.GetTags() {

				// check it's not been completed already
				if complete {
					continue
				}

				// TODO: double dipping here - change the recieving method to do this
				// get the service details
				service := services.ServiceRegister[tag]

				// run the service request
				if err := service.SendRequest(v); err != nil {
					return err
				}
			}

			// set the status to announced
			v.Metadata.SetStatus(records.Status_announced)

			// dequeue the sample
			v.Metadata.AddComment("run announced.")
			v.Metadata.SetStatus(records.Status_announced)
			if err := herald.updateRecord(v); err != nil {
				return err
			}
			herald.announcementQueue.Remove(request)
		}
	}

	// process the remaining queue (should just be samples now)
	for request := herald.announcementQueue.Front(); request != nil; request = request.Next() {

		// grab the sample that is first in the queue
		sample := request.Value.(*records.Sample)

		// TODO:
		// evalute the sample
		// update fields and propogate to linked data
		// decide if it should be dequeued

		// make the service requests
		for tag, complete := range sample.Metadata.GetTags() {

			// check it's not been completed already
			if complete {
				continue
			}

			// TODO: double dipping here - change the recieving method to do this
			// get the service details
			service := services.ServiceRegister[tag]

			// run the service request
			if err := service.SendRequest(sample); err != nil {
				return err
			}
		}

		// update the status of the sample and dequeue it
		sample.Metadata.AddComment("sample announced.")
		sample.Metadata.SetStatus(records.Status_announced)
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
