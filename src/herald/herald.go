// Package herald acts as an interface between the storage and server packages
package herald

import (
	"fmt"
	"sync"

	"github.com/will-rowe/herald/src/sample"
	"github.com/will-rowe/herald/src/server"
	"github.com/will-rowe/herald/src/storage"
)

// Herald is the struct for holding runtime data
type Herald struct {
	sync.Mutex                  // to make the UI binding thread safe
	service    *server.Server   // orchestrates the announcements
	store      *storage.Storage // the key-value store for the samples

	// runtime info for JS:
	storeLocation        string   // where the store is located on disk
	sampleCount          int      // the number of samples currently in the store
	untaggedSampleCount  int      // the number of samples in the store that are untagged
	taggedSampleCount    int      // the number of samples in the store that are tagged with at least one process
	announcedSampleCount int      // the number of samples in the store that are currently being announced
	sampleLabels         []string // used to store all the label names in memory (for JS to access)
}

// InitHerald will initiate the Herald instance
func InitHerald(storeLocation string) (*Herald, error) {

	// load the store
	var store *storage.Storage
	var err error
	if store, err = storage.OpenStorage(storeLocation); err != nil {
		return nil, err
	}

	// return an instance
	return &Herald{
		store:         store,
		storeLocation: storeLocation,
	}, nil
}

// CheckAllSamples makes a pass of the sample store and populates the Herald instance with data:
// - how many samples are in the storage
// - notes any samples with tags
// - loads all sample labels into a slice (for JS to access)
func (herald *Herald) CheckAllSamples() error {
	herald.Lock()
	defer herald.Unlock()

	// reset the runtime data
	herald.sampleCount = 0
	herald.untaggedSampleCount = 0
	herald.taggedSampleCount = 0
	herald.announcedSampleCount = 0

	// get the sample count from the store
	herald.sampleCount = herald.store.GetNumEntries()

	// create the holders
	herald.sampleLabels = make([]string, herald.sampleCount)

	// range over the store key channel (sample labels)
	i := 0
	for label := range herald.store.GetLabels() {

		// add the sample label to the holder
		herald.sampleLabels[i] = string(label)

		// get the full sample
		sample, err := herald.store.GetSample(herald.sampleLabels[i])
		if err != nil {
			return err
		}

		// check the status, update counts and refresh queues
		// todo: refresh queues
		status := sample.GetStatus().String()
		switch status {
		case "untagged":
			herald.untaggedSampleCount++
		case "tagged":
			herald.taggedSampleCount++
		case "announced":
			herald.announcedSampleCount++
		default:
			return fmt.Errorf("encountered sample with unknown status (%v)", status)
		}

		i++
	}
	return nil
}

// Destroy will properly close down the Herald instance and sync the store to disk
func (herald *Herald) Destroy() error {
	herald.Lock()
	defer herald.Unlock()
	return herald.store.CloseStorage()
}

// WipeStorage will clear all samples from storage
func (herald *Herald) WipeStorage() error {
	herald.Lock()
	defer herald.Unlock()
	if err := herald.store.Wipe(); err != nil {
		return err
	}
	herald.sampleCount = 0
	return nil
}

// GetDbPath returns the location of the storage on disk
func (herald *Herald) GetDbPath() string {
	herald.Lock()
	defer herald.Unlock()
	return herald.storeLocation
}

// GetSampleCount returns the current number of samples in storage
func (herald *Herald) GetSampleCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.sampleCount
}

// GetUntaggedSampleCount returns the current number of samples in storage that are untagged
func (herald *Herald) GetUntaggedSampleCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.untaggedSampleCount
}

// GetTaggedSampleCount returns the current number of samples in storage that are tagged with at least one process
func (herald *Herald) GetTaggedSampleCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.taggedSampleCount
}

// GetAnnouncedSampleCount returns the current number of samples that have been announced
func (herald *Herald) GetAnnouncedSampleCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.announcedSampleCount
}

// CreateSample creates a sample record, updates the runtime info and adds the record to storage
// TODO: this might be bypassed later and instead get JS to encode the form to protobuf directly
func (herald *Herald) CreateSample(label string, barcode int32, comment string, tags []string) error {
	herald.Lock()
	defer herald.Unlock()

	// create the sample
	sample := sample.InitSample(label, barcode, comment)

	// tag the sample
	if len(tags) != 0 {
		if err := sample.AddTags(tags); err != nil {
			return err
		}
	}

	// add the sample to the store
	if err := herald.store.AddSample(sample); err != nil {
		return err
	}

	// update the runtime info (grow the label slice, tag slice etc.)
	herald.sampleLabels = append(herald.sampleLabels, label)
	herald.sampleCount++
	if len(tags) != 0 {
		herald.taggedSampleCount++
	}
	return nil
}

// GetSampleStatus will check the status of a sample
func (herald *Herald) GetSampleStatus(sampleLabel string) (string, error) {
	herald.Lock()
	defer herald.Unlock()

	// get the sample from storage
	sample, err := herald.store.GetSample(sampleLabel)
	if err != nil {
		return "", err
	}

	// return status
	return sample.GetStatus().String(), nil
}

// DeleteSample removes a sample record from storage and updates the counts
func (herald *Herald) DeleteSample(sampleLabel string) error {
	herald.Lock()
	defer herald.Unlock()

	// get the sample from storage
	sample, err := herald.store.GetSample(sampleLabel)
	if err != nil {
		return err
	}

	// get the status
	status := sample.GetStatus().String()

	// prevent announced samples from being deleted
	if status == "announced" {
		return fmt.Errorf("can't delete sample during announcement: %v", sampleLabel)
	}

	// delete the sample from the store
	if err := herald.store.DeleteSample(sampleLabel); err != nil {
		return err
	}

	// update the counts
	herald.sampleCount--
	if status == "tagged" {
		herald.taggedSampleCount--
	}
	return nil
}

// PrintSampleToJSONstring collects a sample from the database and returns a string of the sample protobuf data in JSON
func (herald *Herald) PrintSampleToJSONstring(label string) string {
	herald.Lock()
	defer herald.Unlock()

	// TODO: check the error from GetSampleJSONDump method
	sampleString, _ := herald.store.GetSampleJSONDump(label)

	return sampleString
}

// GetSampleLabel is used by JS to collect a sample label from the runtime slice of sample labels
// NOTE: this assumes the caller has already run GetSampleCount (or similar) to find the iterator range
// TODO: add error on return too (will require re-write of JS function)
func (herald *Herald) GetSampleLabel(iterator int) string {
	herald.Lock()
	defer herald.Unlock()
	return herald.sampleLabels[iterator]
}
