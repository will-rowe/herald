// Package storage wraps a bit cask as the disk-backed key-value store for sample information
package storage

import (
	"bytes"
	"fmt"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/prologic/bitcask"

	"github.com/will-rowe/herald/src/sample"
)

// dbMaxEntries is used to cap the number of samples that can be added
const dbMaxEntries = 10

// useSync will run sync on every bit cask transaction, improving stability at the expense of time
const useSync = true

// Storage holds a bitcask db and some extra stuff
type Storage struct {
	db         *bitcask.Bitcask // the key-value store
	dbLocation string           // where the store is stored
}

// OpenStorage will open up the database and return a storage struct and any error
func OpenStorage(dbLocation string) (*Storage, error) {

	// open the database
	db, err := bitcask.Open(dbLocation, bitcask.WithSync(useSync))
	if err != nil {
		return nil, err
	}

	// create the storage struct
	store := &Storage{
		db:         db,
		dbLocation: dbLocation,
	}
	return store, nil
}

// CloseStorage will flush and close the storage database
func (storage *Storage) CloseStorage() error {
	if err := storage.db.Sync(); err != nil {
		return err
	}
	if err := storage.db.Close(); err != nil {
		return err
	}
	return nil
}

// Wipe clears all entries from the bit cask
func (storage *Storage) Wipe() error {
	return storage.db.DeleteAll()
}

// GetNumEntries returns the current number of entries in the bit cask
func (storage *Storage) GetNumEntries() int {
	return storage.db.Len()
}

// GetLabels returns a channel of sample labels (keys) held in storage
func (storage *Storage) GetLabels() chan []byte {
	return storage.db.Keys()
}

// DeleteSample is a method to remove a protobuf encoded sample from the bit cask
func (storage *Storage) DeleteSample(sampleLabel string) error {
	return storage.db.Delete([]byte(sampleLabel))
}

// AddSample is a method to marshal a sample and add it to the bit cask
func (storage *Storage) AddSample(sample *sample.Sample) error {

	// check the DB limit hasn't been reached
	if storage.db.Len() == dbMaxEntries {
		return fmt.Errorf("database entry limit reached (%d)", dbMaxEntries)
	}

	// check the sample is not in the database already
	if storage.db.Has([]byte(sample.Label)) {
		return fmt.Errorf("duplicate label can't be added to the database (%s)", sample.Label)
	}

	// marshal the sample
	data, err := proto.Marshal(sample)
	if err != nil {
		return err
	}

	// add the sample
	if err := storage.db.Put([]byte(sample.Label), data); err != nil {
		return err
	}
	return nil
}

// GetSample is a method to retrieve a sample from the bit cask and unmarshal it to a struct
func (storage *Storage) GetSample(sampleLabel string) (*sample.Sample, error) {

	// get the sample from the bit cask
	data, err := storage.db.Get([]byte(sampleLabel))
	if err != nil {
		return nil, err
	}

	// unmarshal the sample
	sample := &sample.Sample{}
	if err := proto.Unmarshal(data, sample); err != nil {
		return nil, err
	}
	return sample, nil
}

// GetSampleProtoDump is a method to retrieve a sample from the bit cask and return a string dump of the protobuf message
func (storage *Storage) GetSampleProtoDump(sampleLabel string) (string, error) {

	// get the sample from the bit cask
	data, err := storage.db.Get([]byte(sampleLabel))
	if err != nil {
		return "", err
	}

	// unmarshal the sample
	sample := &sample.Sample{}
	if err := proto.Unmarshal(data, sample); err != nil {
		return "", err
	}

	return proto.MarshalTextString(sample), nil
}

// GetSampleJSONDump is a method to retrieve a sample from the bit cask and return a string dump of the protobuf message in JSON
func (storage *Storage) GetSampleJSONDump(sampleLabel string) (string, error) {

	// get the sample from the bit cask
	data, err := storage.db.Get([]byte(sampleLabel))
	if err != nil {
		return "", err
	}

	// unmarshal the sample
	sample := &sample.Sample{}
	if err := proto.Unmarshal(data, sample); err != nil {
		return "", err
	}

	// convert to JSON
	buf := &bytes.Buffer{}
	jsonMarshaller := jsonpb.Marshaler{
		EnumsAsInts:  false, // Whether to render enum values as integers, as opposed to string values.
		EmitDefaults: true,  // Whether to render fields with zero values
		Indent:       "\t",  // A string to indent each level by
		OrigName:     false, // Whether to use the original (.proto) name for fields
	}
	jsonMarshaller.Marshal(buf, sample)

	return string(buf.Bytes()), nil
}
