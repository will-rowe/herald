package storage

import (
	"fmt"
	"os"
	"testing"

	"github.com/will-rowe/herald/src/records"
)

// TestStorage open and close
func TestStorageIO(t *testing.T) {

	// setup the storage
	sampleStore, err := OpenStorage("./tmp")
	if err != nil {
		t.Fatal(err)
	}

	// close the storage
	if err := sampleStore.CloseStorage(); err != nil {
		t.Fatal(err)
	}
}

// TestAdd
func TestStorageAdd(t *testing.T) {

	// setup the storage
	sampleStore, err := OpenStorage("./tmp")
	if err != nil {
		t.Fatal(err)
	}

	// check the current number of samples
	if sampleStore.GetNumSamples() != 0 {
		t.Fatalf("database not empty: %d", sampleStore.GetNumSamples())
	}

	// add samples
	var i int32
	for i = 0; i < 9; i++ {
		sampleName := fmt.Sprintf("sample %d", i)
		sample := records.InitSample(sampleName, "testRun", i)
		if err := sampleStore.AddSample(sample); err != nil {
			t.Fatal(err)
		}
	}

	// check db now has correct number of samples
	if sampleStore.GetNumSamples() != 9 {
		t.Fatalf("incorrect number of samples added to db: %d", sampleStore.GetNumSamples())
	}

	// check you can't add a duplicate label
	sample := records.InitSample("sample 1", "testRun", 1)
	if err := sampleStore.AddSample(sample); err == nil {
		t.Fatal(err)
	}
	if sampleStore.GetNumSamples() != 9 {
		t.Fatalf("incorrect number of samples added to db: %d", sampleStore.GetNumSamples())
	}

	// check you can retrieve a sample
	sampleCopy, err := sampleStore.GetSample(sample.Metadata.GetLabel())
	if err != nil {
		t.Fatal(err)
	}
	if sample.Metadata.GetLabel() != sampleCopy.Metadata.GetLabel() {
		t.Fatal("retrieved sample label does not match inserted sample label")
	}

	// test a JSON dump
	jsonDump, err := sampleStore.GetSampleJSONDump(sample.Metadata.GetLabel())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(jsonDump)

	// check you can delete a sample
	if err := sampleStore.DeleteSample(sample.Metadata.GetLabel()); err != nil {
		t.Fatal(err)
	}
	if sampleStore.GetNumSamples() != 8 {
		t.Fatal("num entries does not updated after sample is deleted")
	}

	// check the wipe
	if err := sampleStore.Wipe(); err != nil {
		t.Fatal(err)
	}
	if sampleStore.GetNumSamples() != 0 {
		t.Fatalf("db was not wiped (%d keys left)", sampleStore.GetNumSamples())
	}

	// close the storage
	if err := sampleStore.CloseStorage(); err != nil {
		t.Fatal(err)
	}

	// clean up
	os.RemoveAll("./tmp/")
}
