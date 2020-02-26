// Package herald acts as an interface between the storage and server packages
package herald

import (
	"sync"

	"github.com/will-rowe/herald/src/sample"
	"github.com/will-rowe/herald/src/server"
	"github.com/will-rowe/herald/src/storage"
)

// Herald is the struct for holding runtime data
type Herald struct {
	sync.Mutex                     // to make the UI binding thread safe
	Service       *server.Server   // orchestrates the announcements
	store         *storage.Storage // the key-value store for the samples
	storeLocation string           // where the store is located on disk
	sampleCount   int              // the number of samples currently in the store
	sampleLabels  []string         // used to store all the label names in memory (for JS to access)
	taggedSamples []*sample.Sample // tagged samples
}

// InitHerald will initiate the Herald instance
func InitHerald(storeLocation string) (*Herald, error) {

	// load the store
	var store *storage.Storage
	var err error
	if store, err = storage.OpenStorage(storeLocation); err != nil {
		return nil, err
	}

	// create an instance
	herald := &Herald{
		store:         store,
		storeLocation: storeLocation,
	}

	// populate the instance with runtime data and return
	herald.CheckAllSamples()
	return herald, nil
}

// CheckAllSamples makes a pass of the sample store and populates the Herald instance with data:
// - how many samples are in the storage
// - notes any samples with tags
// - loads all sample labels into a slice (for JS to access)
func (herald *Herald) CheckAllSamples() {
	herald.Lock()
	defer herald.Unlock()

	// get the count
	herald.sampleCount = herald.store.GetNumEntries()

	// create the holders
	herald.sampleLabels = make([]string, herald.sampleCount)

	// range over the store keys
	i := 0
	for label := range herald.store.GetLabels() {
		herald.sampleLabels[i] = string(label)
		i++
	}
	return
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
	return nil
}

// TagSample will add a tag to a sample
func (herald *Herald) TagSample(sampleLabel string, tags []string) error {
	return nil
}

// DeleteSample removes a sample record from storage and updates the counts
func (herald *Herald) DeleteSample(label string) error {
	herald.Lock()
	defer herald.Unlock()

	// TODO: check the sample isn't being used (has been announced)

	// TODO: add delete logic (inverse of the CreateSample)

	if err := herald.store.DeleteSample(label); err != nil {
		return err
	}
	herald.sampleCount--
	return nil
}

// GetSampleLabel is used by JS to collect a sample label from the runtime list
func (herald *Herald) GetSampleLabel(iterator int) string {
	herald.Lock()
	defer herald.Unlock()
	return herald.sampleLabels[iterator]
}

// PrintSampleToString collects a sample from the database and returns a string of the sample protobuf data
func (herald *Herald) PrintSampleToString(label string) string {
	herald.Lock()
	defer herald.Unlock()

	// TODO: check the error from getSampleDump method
	sampleString, _ := herald.store.GetSampleDump(label)

	return sampleString
}
