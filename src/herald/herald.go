// Package herald acts as an interface between the storage and service packages
package herald

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/will-rowe/herald/src/services"
	"github.com/will-rowe/herald/src/storage"
)

// Herald is the struct for holding runtime data
type Herald struct {
	sync.Mutex                         // to make the UI binding thread safe
	store             *storage.Storage // the key-value store for the samples
	announcementQueue *list.List       // a FIFO queue for announcements

	// runtime info for JS:
	storeLocation        string     // where the store is located on disk
	experimentCount      int        // the number of experiments currently in the store
	sampleCount          int        // the number of samples currently in the store
	untaggedSampleCount  int        // the number of samples in the store that are untagged
	taggedSampleCount    int        // the number of samples in the store that are tagged with at least one process
	announcedSampleCount int        // the number of samples in the store that are currently being announced
	sampleDetails        [][]string // used to store all the sample labels, creation dates and corresponding experiment in memory (for JS to access)
	experimentLabels     []string   // used to store all the experiment names in memory (for JS to access)
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
		store:             store,
		announcementQueue: list.New(),
		storeLocation:     storeLocation,
		sampleDetails:     make([][]string, 3),
	}, nil
}

// GetRuntimeInfo makes a pass of the experiment and sample stores before populating the Herald instance with data:
// - how many samples are in the storage
// - notes any samples with tags
// - loads all sample labels into a slice (for JS to access)
func (herald *Herald) GetRuntimeInfo() error {
	herald.Lock()
	defer herald.Unlock()

	// reset the runtime data
	herald.experimentCount = 0
	herald.sampleCount = 0
	herald.untaggedSampleCount = 0
	herald.taggedSampleCount = 0
	herald.announcedSampleCount = 0

	// get the experiment and sample counts from the store
	herald.experimentCount = herald.store.GetNumExperiments()
	herald.sampleCount = herald.store.GetNumSamples()

	// create the holders
	herald.experimentLabels = make([]string, herald.experimentCount)
	herald.sampleDetails = make([][]string, 3)
	for i := 0; i < 3; i++ {
		herald.sampleDetails[i] = make([]string, herald.sampleCount)
	}

	// range over the sample labels via the key channel from the bit cask
	i := 0
	for label := range herald.store.GetSampleLabels() {

		// get the full sample
		sample, err := herald.store.GetSample(string(label))
		if err != nil {
			return err
		}

		// add the details to the store
		herald.sampleDetails[0][i] = sample.Metadata.GetLabel()
		herald.sampleDetails[1][i] = sample.Metadata.Created.String()
		herald.sampleDetails[2][i] = sample.GetParentExperiment()

		// check the status, update counts and refresh queues
		// todo: refresh queues
		status := sample.Metadata.GetStatus().String()
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

	// range over the experiment labels via the key channel from the bit cask
	i = 0
	for label := range herald.store.GetExperimentLabels() {
		herald.experimentLabels[i] = string(label)
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

// WipeStorage will clear all samples and experiments from storage
func (herald *Herald) WipeStorage() error {
	herald.Lock()
	defer herald.Unlock()
	if err := herald.store.Wipe(); err != nil {
		return err
	}
	herald.sampleCount = 0
	herald.experimentCount = 0
	return nil
}

// GetDbPath returns the location of the storage on disk
func (herald *Herald) GetDbPath() string {
	herald.Lock()
	defer herald.Unlock()
	return herald.storeLocation
}

// GetExperimentCount returns the current number of experiments in storage
func (herald *Herald) GetExperimentCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.experimentCount
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

// CreateExperiment creates an experiment record, updates the runtime info and adds the record to storage
// TODO: this might be bypassed later and instead get JS to encode the form to protobuf directly
func (herald *Herald) CreateExperiment(expLabel, outDir, fast5Dir, fastqDir, comment string, tags []string) error {
	herald.Lock()
	defer herald.Unlock()

	// create the experiment
	exp := services.InitExperiment(expLabel, outDir, fast5Dir, fastqDir)

	// add any comment
	if len(comment) != 0 {
		if err := exp.Metadata.AddComment(comment); err != nil {
			return err
		}
	}

	// tag the experiment, add to announcement queue and update it's status
	if len(tags) != 0 {
		if err := exp.Metadata.AddTags(tags); err != nil {
			return err
		}
		herald.announcementQueue.PushBack(exp)
	}

	// add the experiment to the store
	if err := herald.store.AddExperiment(exp); err != nil {
		return err
	}

	// update the runtime info (grow the label slice, tag slice etc.)
	herald.experimentLabels = append(herald.experimentLabels, expLabel)
	herald.experimentCount++
	return nil
}

// CreateSample creates a sample record, updates the runtime info and adds the record to storage
// TODO: this might be bypassed later and instead get JS to encode the form to protobuf directly
func (herald *Herald) CreateSample(label string, experimentName string, barcode int32, comment string, tags []string) error {
	herald.Lock()
	defer herald.Unlock()

	// get the experiment from storage
	exp, err := herald.store.GetExperiment(experimentName)
	if err != nil {
		return err
	}

	// copy the tags over to the samples (sequence and basecall)
	tags = append(exp.Metadata.GetRequestOrder(), tags...)

	// create the sample
	sample := services.InitSample(label, exp.Metadata.GetLabel(), barcode)
	if len(comment) != 0 {
		if err := sample.Metadata.AddComment(comment); err != nil {
			return err
		}
	}

	// tag the sample, add to the announcement queue and update status
	if len(tags) != 0 {
		if err := sample.Metadata.AddTags(tags); err != nil {
			return err
		}
		herald.announcementQueue.PushBack(sample)
	}

	// add the sample to the store
	if err := herald.store.AddSample(sample); err != nil {
		return err
	}

	// update the runtime info (grow the label slice, tag slice etc.)
	herald.sampleDetails[0] = append(herald.sampleDetails[0], label)
	herald.sampleDetails[1] = append(herald.sampleDetails[1], sample.Metadata.GetCreated().String())
	herald.sampleDetails[2] = append(herald.sampleDetails[2], experimentName)
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
	return sample.Metadata.GetStatus().String(), nil
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
	status := sample.Metadata.GetStatus().String()

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

// GetSampleLabel is used by JS to collect a sample label from the runtime slice of sample data
// NOTE: this assumes the caller has already run GetSampleCount (or similar) to find the iterator range
// TODO: add error on return too (will require re-write of JS function)
func (herald *Herald) GetSampleLabel(iterator int) string {
	herald.Lock()
	defer herald.Unlock()
	return herald.sampleDetails[0][iterator]
}

// GetSampleCreation is used by JS to collect a sample created timestamp from the runtime slice of sample data
// NOTE: this assumes the caller has already run GetSampleCount (or similar) to find the iterator range
// TODO: add error on return too (will require re-write of JS function)
func (herald *Herald) GetSampleCreation(iterator int) string {
	herald.Lock()
	defer herald.Unlock()
	return herald.sampleDetails[1][iterator]
}

// GetSampleExperiment is used by JS to collect a sample experiment name from the runtime slice of sample data
// NOTE: this assumes the caller has already run GetSampleCount (or similar) to find the iterator range
// TODO: add error on return too (will require re-write of JS function)
func (herald *Herald) GetSampleExperiment(iterator int) string {
	herald.Lock()
	defer herald.Unlock()
	return herald.sampleDetails[2][iterator]
}

// GetLabel is used by JS to collect an experiment name from the runtime slice of experiment names
// NOTE: this assumes the caller has already run GetExperimentCount (or similar) to find the iterator range
// TODO: add error on return too (will require re-write of JS function)
func (herald *Herald) GetLabel(iterator int) string {
	herald.Lock()
	defer herald.Unlock()
	return herald.experimentLabels[iterator]
}
