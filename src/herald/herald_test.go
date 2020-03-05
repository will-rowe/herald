package herald

import (
	"os"
	"testing"

	"github.com/will-rowe/herald/src/data"
)

// TestHeraldInit
func TestHeraldInit(t *testing.T) {

	// setup the storage
	tmp, err := InitHerald("./tmp")
	if err != nil {
		t.Fatal(err)
	}

	// close the storage
	if err := tmp.Destroy(); err != nil {
		t.Fatal(err)
	}
}

// TestHeraldCreateSample
func TestHeraldCreateSample(t *testing.T) {

	// open the storage
	tmp, err := InitHerald("./tmp")
	if err != nil {
		t.Fatal(err)
	}

	// create an experiment
	testExpName := "test experiment"
	testExp := data.InitExperiment(testExpName, "", "", "")
	if err := tmp.store.AddExperiment(testExp); err != nil {
		t.Fatal(err)
	}

	// create and add a sample
	testLabel := "testLabel"
	if err := tmp.CreateSample(testLabel, testExpName, 1, "test comment", []string{"sequence"}); err != nil {
		t.Fatal(err)
	}

	// check runtime info was updated
	count := tmp.GetSampleCount()
	if count != 1 {
		t.Fatal("herald sample count not updated (should be 1)")
	}
	if storedLabel := tmp.GetSampleLabel(0); storedLabel != testLabel {
		t.Fatalf("stored label does not match that used during sample creation (%v vs %v)", storedLabel, testLabel)
	}

	// check the status of the sample
	status, err := tmp.GetSampleStatus(testLabel)
	if err != nil {
		t.Fatal(err)
	}
	if status != "tagged" {
		t.Fatalf("sample status should be tagged, not: %v", status)
	}

	// close the storage
	if err := tmp.Destroy(); err != nil {
		t.Fatal(err)
	}

	// clean up
	os.RemoveAll("./tmp/")

}

// TestHeraldCreateExperiment
func TestHeraldCreateExperiment(t *testing.T) {

	// open the storage
	tmp, err := InitHerald("./tmp")
	if err != nil {
		t.Fatal(err)
	}

	// create and add an experiment
	testName := "test experiment"
	if err := tmp.CreateExperiment(testName, "/tmp", "/tmp/fast5_pass", "/tmp/fastq_pass", "", []string{"sequence", "basecall"}); err != nil {
		t.Fatal(err)
	}

	// check runtime info was updated
	count := tmp.GetExperimentCount()
	if count != 1 {
		t.Fatal("herald experiment count not updated (should be 1)")
	}
	if storedName := tmp.GetLabel(0); storedName != testName {
		t.Fatalf("stored label does not match that used during sample creation (%v vs %v)", storedName, testName)
	}

	// close the storage
	if err := tmp.Destroy(); err != nil {
		t.Fatal(err)
	}

	// clean up
	os.RemoveAll("./tmp/")
}
