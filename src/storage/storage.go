// Package storage wraps two bit casks as the disk-backed key-value store for sample information and run information
package storage

import (
	"bytes"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"git.mills.io/prologic/bitcask"

	"github.com/will-rowe/herald/src/records"
)

// dbMaxEntries is used to cap the number of db elements that can be added
const dbMaxEntries = 10000

// useSync will run sync on every bit cask transaction, improving stability at the expense of time
const useSync = true

// Storage holds a bitcask db and some extra stuff
type Storage struct {
	sampleDB   *bitcask.Bitcask // the key-value store for samples
	runDB      *bitcask.Bitcask // the key-value store for runs
	dbLocation string           // where the store is stored
}

// OpenStorage will create/open up the databases and return a storage struct or an error
func OpenStorage(dbLocation string) (*Storage, error) {

	// get the names for both databases
	sampleDBname := fmt.Sprintf("%s/sampleCask", dbLocation)
	runDBname := fmt.Sprintf("%s/runCask", dbLocation)

	// open the databases
	sdb, err := bitcask.Open(sampleDBname, bitcask.WithSync(useSync))
	if err != nil {
		return nil, err
	}
	edb, err := bitcask.Open(runDBname, bitcask.WithSync(useSync))
	if err != nil {
		return nil, err
	}

	// create the storage struct
	store := &Storage{
		sampleDB:   sdb,
		runDB:      edb,
		dbLocation: dbLocation,
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
	if err := storage.runDB.Sync(); err != nil {
		return err
	}
	if err := storage.runDB.Close(); err != nil {
		return err
	}
	return nil
}

// Wipe clears all entries from the samples and runs databases
func (storage *Storage) Wipe() error {
	if err := storage.runDB.DeleteAll(); err != nil {
		return err
	}
	return storage.sampleDB.DeleteAll()
}

// GetNumSamples returns the current number of samples in storage
func (storage *Storage) GetNumSamples() int {
	return storage.sampleDB.Len()
}

// GetNumRuns returns the current number of runs in storage
func (storage *Storage) GetNumRuns() int {
	return storage.runDB.Len()
}

// GetSampleLabels returns a channel of sample labels (keys) held in storage
func (storage *Storage) GetSampleLabels() chan []byte {
	return storage.sampleDB.Keys()
}

// GetRunLabels returns a channel of run names (keys) held in storage
func (storage *Storage) GetRunLabels() chan []byte {
	return storage.runDB.Keys()
}

// DeleteSample is a method to remove a sample from storage
func (storage *Storage) DeleteSample(sampleLabel string) error {
	return storage.sampleDB.Delete([]byte(sampleLabel))
}

// DeleteRun is a method to remove an run from storage
func (storage *Storage) DeleteRun(runName string) error {
	return storage.runDB.Delete([]byte(runName))
}

// AddSample is a method to marshal a sample and store it
func (storage *Storage) AddSample(sample *records.Sample) error {

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

// AddRun is a method to marshal an run and store it
func (storage *Storage) AddRun(run *records.Run) error {

	// check the DB limit hasn't been reached
	if storage.runDB.Len() == dbMaxEntries {
		return fmt.Errorf("database entry limit reached (%d)", dbMaxEntries)
	}

	// check the sample is not in the database already
	if storage.runDB.Has([]byte(run.Metadata.Label)) {
		return fmt.Errorf("duplicate run name can't be added to the database (%s)", run.Metadata.Label)
	}

	// marshal the sample
	data, err := proto.Marshal(run)
	if err != nil {
		return err
	}

	// add the sample
	if err := storage.runDB.Put([]byte(run.Metadata.Label), data); err != nil {
		return err
	}
	return nil
}

// GetSample is a method to retrieve a sample from storage and unmarshal it to a struct
func (storage *Storage) GetSample(sampleLabel string) (*records.Sample, error) {

	// get the sample from the bit cask
	dbData, err := storage.sampleDB.Get([]byte(sampleLabel))
	if err != nil {
		return nil, err
	}

	// unmarshal the sample
	sample := &records.Sample{}
	if err := proto.Unmarshal(dbData, sample); err != nil {
		return nil, err
	}
	return sample, nil
}

// GetRun is a method to retrieve an run from storage and unmarshal it to a struct
func (storage *Storage) GetRun(runName string) (*records.Run, error) {

	// get the run from the bit cask
	dbData, err := storage.runDB.Get([]byte(runName))
	if err != nil {
		return nil, err
	}

	// unmarshal the sample
	exp := &records.Run{}
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
	sample := &records.Sample{}
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
	sample := &records.Sample{}
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
