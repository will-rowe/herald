// Package storage wraps two bit casks as the disk-backed key-value store for sample information and experiment information
package storage

import (
	"bytes"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/prologic/bitcask"

	"github.com/will-rowe/herald/src/data"
)

// dbMaxEntries is used to cap the number of db elements that can be added
const dbMaxEntries = 10000

// useSync will run sync on every bit cask transaction, improving stability at the expense of time
const useSync = true

// Storage holds a bitcask db and some extra stuff
type Storage struct {
	sampleDB     *bitcask.Bitcask // the key-value store for samples
	experimentDB *bitcask.Bitcask // the key-value store for experiments
	dbLocation   string           // where the store is stored
}

// OpenStorage will create/open up the databases and return a storage struct or an error
func OpenStorage(dbLocation string) (*Storage, error) {

	// get the names for both databases
	sampleDBname := fmt.Sprintf("%s/sampleCask", dbLocation)
	experimentDBname := fmt.Sprintf("%s/experimentCask", dbLocation)

	// open the databases
	sdb, err := bitcask.Open(sampleDBname, bitcask.WithSync(useSync))
	if err != nil {
		return nil, err
	}
	edb, err := bitcask.Open(experimentDBname, bitcask.WithSync(useSync))
	if err != nil {
		return nil, err
	}

	// create the storage struct
	store := &Storage{
		sampleDB:     sdb,
		experimentDB: edb,
		dbLocation:   dbLocation,
	}
	return store, nil
}

// CloseStorage will flush and close the storage databases
func (storage *Storage) CloseStorage() error {
	if err := storage.sampleDB.Sync(); err != nil {
		return err
	}
	if err := storage.sampleDB.Close(); err != nil {
		return err
	}
	if err := storage.experimentDB.Sync(); err != nil {
		return err
	}
	if err := storage.experimentDB.Close(); err != nil {
		return err
	}
	return nil
}

// Wipe clears all entries from the samples and experiments databases
func (storage *Storage) Wipe() error {
	if err := storage.experimentDB.DeleteAll(); err != nil {
		return err
	}
	return storage.sampleDB.DeleteAll()
}

// GetNumSamples returns the current number of samples in storage
func (storage *Storage) GetNumSamples() int {
	return storage.sampleDB.Len()
}

// GetNumExperiments returns the current number of experiments in storage
func (storage *Storage) GetNumExperiments() int {
	return storage.experimentDB.Len()
}

// GetSampleLabels returns a channel of sample labels (keys) held in storage
func (storage *Storage) GetSampleLabels() chan []byte {
	return storage.sampleDB.Keys()
}

// GetExperimentNames returns a channel of experiment names (keys) held in storage
func (storage *Storage) GetExperimentNames() chan []byte {
	return storage.experimentDB.Keys()
}

// DeleteSample is a method to remove a sample from storage
func (storage *Storage) DeleteSample(sampleLabel string) error {
	return storage.sampleDB.Delete([]byte(sampleLabel))
}

// DeleteExperiment is a method to remove an experiment from storage
func (storage *Storage) DeleteExperiment(experimentName string) error {
	return storage.experimentDB.Delete([]byte(experimentName))
}

// AddSample is a method to marshal a sample and store it
func (storage *Storage) AddSample(sample *data.Sample) error {

	// check the DB limit hasn't been reached
	if storage.sampleDB.Len() == dbMaxEntries {
		return fmt.Errorf("database entry limit reached (%d)", dbMaxEntries)
	}

	// check the sample is not in the database already
	if storage.sampleDB.Has([]byte(sample.Metadata.Label)) {
		return fmt.Errorf("duplicate label can't be added to the database (%s)", sample.Metadata.Label)
	}

	// marshal the sample
	data, err := proto.Marshal(sample)
	if err != nil {
		return err
	}

	// add the sample
	if err := storage.sampleDB.Put([]byte(sample.Metadata.Label), data); err != nil {
		return err
	}
	return nil
}

// AddExperiment is a method to marshal an experiment and store it
func (storage *Storage) AddExperiment(experiment *data.Experiment) error {

	// check the DB limit hasn't been reached
	if storage.experimentDB.Len() == dbMaxEntries {
		return fmt.Errorf("database entry limit reached (%d)", dbMaxEntries)
	}

	// check the sample is not in the database already
	if storage.experimentDB.Has([]byte(experiment.Metadata.Label)) {
		return fmt.Errorf("duplicate experiment name can't be added to the database (%s)", experiment.Metadata.Label)
	}

	// marshal the sample
	data, err := proto.Marshal(experiment)
	if err != nil {
		return err
	}

	// add the sample
	if err := storage.experimentDB.Put([]byte(experiment.Metadata.Label), data); err != nil {
		return err
	}
	return nil
}

// GetSample is a method to retrieve a sample from storage and unmarshal it to a struct
func (storage *Storage) GetSample(sampleLabel string) (*data.Sample, error) {

	// get the sample from the bit cask
	dbData, err := storage.sampleDB.Get([]byte(sampleLabel))
	if err != nil {
		return nil, err
	}

	// unmarshal the sample
	sample := &data.Sample{}
	if err := proto.Unmarshal(dbData, sample); err != nil {
		return nil, err
	}
	return sample, nil
}

// GetExperiment is a method to retrieve an experiment from storage and unmarshal it to a struct
func (storage *Storage) GetExperiment(experimentName string) (*data.Experiment, error) {

	// get the experiment from the bit cask
	dbData, err := storage.experimentDB.Get([]byte(experimentName))
	if err != nil {
		return nil, err
	}

	// unmarshal the sample
	exp := &data.Experiment{}
	if err := proto.Unmarshal(dbData, exp); err != nil {
		return nil, err
	}
	return exp, nil
}

// GetSampleProtoDump is a method to retrieve a sample from storage and return a string dump of the protobuf message
func (storage *Storage) GetSampleProtoDump(sampleLabel string) (string, error) {

	// get the sample from the bit cask
	dbData, err := storage.sampleDB.Get([]byte(sampleLabel))
	if err != nil {
		return "", err
	}

	// unmarshal the sample
	sample := &data.Sample{}
	if err := proto.Unmarshal(dbData, sample); err != nil {
		return "", err
	}

	return proto.MarshalTextString(sample), nil
}

// GetSampleJSONDump is a method to retrieve a sample from storage and return a string dump of the protobuf message in JSON
func (storage *Storage) GetSampleJSONDump(sampleLabel string) (string, error) {

	// get the sample from the bit cask
	dbData, err := storage.sampleDB.Get([]byte(sampleLabel))
	if err != nil {
		return "", err
	}

	// unmarshal the sample
	sample := &data.Sample{}
	if err := proto.Unmarshal(dbData, sample); err != nil {
		return "", err
	}

	// convert to JSON
	buf := &bytes.Buffer{}
	jsonMarshaller := jsonpb.Marshaler{
		EnumsAsInts:  false, // Whether to render enum values as integers, as opposed to string values.
		EmitDefaults: false, // Whether to render fields with zero values
		Indent:       "\t",  // A string to indent each level by
		OrigName:     false, // Whether to use the original (.proto) name for fields
	}
	jsonMarshaller.Marshal(buf, sample)

	return string(buf.Bytes()), nil
}
