package storage

import (
	"fmt"
	"os"
	"testing"

	"github.com/will-rowe/herald/src/sample"
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
	if sampleStore.GetNumEntries() != 0 {
		t.Fatalf("database not empty: %d", sampleStore.GetNumEntries())
	}

	// add samples
	var i int32
	for i = 0; i < 9; i++ {
		sample := &sample.Sample{
			Label:   fmt.Sprintf("sample %d", i),
			Barcode: i,
		}
		if err := sampleStore.AddSample(sample); err != nil {
			t.Fatal(err)
		}
	}

	// check db now has correct number of samples
	if sampleStore.GetNumEntries() != 9 {
		t.Fatalf("incorrect number of samples added to db: %d", sampleStore.GetNumEntries())
	}

	// check you can't add a duplicate label
	sample := &sample.Sample{
		Label:   "sample 1",
		Barcode: 666,
	}
	if err := sampleStore.AddSample(sample); err == nil {
		t.Fatal(err)
	}
	if sampleStore.GetNumEntries() != 9 {
		t.Fatalf("incorrect number of samples added to db: %d", sampleStore.GetNumEntries())
	}

	// check you can retrieve a sample
	sampleCopy, err := sampleStore.GetSample(sample.Label)
	if err != nil {
		t.Fatal(err)
	}
	if sample.GetLabel() != sampleCopy.GetLabel() {
		t.Fatal("retrieved sample label does not match inserted sample label")
	}
	t.Log(sampleStore.GetSampleDump(sample.Label))

	// check you can delete a sample
	if err := sampleStore.DeleteSample(sample.Label); err != nil {
		t.Fatal(err)
	}
	if sampleStore.GetNumEntries() != 8 {
		t.Fatal("num entries does not updated after sample is deleted")
	}

	// check the wipe
	if err := sampleStore.Wipe(); err != nil {
		t.Fatal(err)
	}
	if sampleStore.GetNumEntries() != 0 {
		t.Fatalf("db was not wiped (%d keys left)", sampleStore.GetNumEntries())
	}

	// close the storage
	if err := sampleStore.CloseStorage(); err != nil {
		t.Fatal(err)
	}

	// clean up
	os.RemoveAll("./tmp/")
}
