package herald

import (
	"os"
	"testing"
)

// TestHerald
func TestHerald(t *testing.T) {

	// open the storage
	tmp, err := InitHerald("./tmp")
	if err != nil {
		t.Fatal(err)
	}

	// create and add a run
	testExpName := "test run"
	if err := tmp.AddRun(testExpName, "/tmp", "/tmp/fast5_pass", "/tmp/fastq_pass", "", []string{"sequence", "basecall"}, false); err != nil {
		t.Fatal(err)
	}

	// check runtime info was updated
	expCount := tmp.GetRunCount()
	if expCount != 1 {
		t.Fatal("herald run count not updated (should be 1)")
	}
	if storedName := tmp.GetLabel(0); storedName != testExpName {
		t.Fatalf("stored label does not match that used during sample creation (%v vs %v)", storedName, testExpName)
	}

	// create and add a sample
	testSampleLabel := "test sample"
	if err := tmp.CreateSample(testSampleLabel, testExpName, 1, "test comment", []string{"sequence"}); err != nil {
		t.Fatal(err)
	}

	// check runtime info was updated
	sampleCount := tmp.GetSampleCount()
	if sampleCount != 1 {
		t.Fatal("herald sample count not updated (should be 1)")
	}
	if storedLabel := tmp.GetSampleLabel(0); storedLabel != testSampleLabel {
		t.Fatalf("stored label does not match that used during sample creation (%v vs %v)", storedLabel, testSampleLabel)
	}

	// close the storage
	if err := tmp.Destroy(); err != nil {
		t.Fatal(err)
	}

	// reopen and check counts
	tmp, err = InitHerald("./tmp")
	if err != nil {
		t.Fatal(err)
	}
	expCount2 := tmp.GetRunCount()
	if expCount2 != 1 {
		t.Fatal("herald run count not updated (should be 1)")
	}
	if storedName2 := tmp.GetLabel(0); storedName2 != testExpName {
		t.Fatalf("stored label does not match that used during sample creation (%v vs %v)", storedName2, testExpName)
	}
	sampleCount2 := tmp.GetSampleCount()
	if sampleCount2 != 1 {
		t.Fatal("herald sample count not updated (should be 1)")
	}
	if storedLabel2 := tmp.GetSampleLabel(0); storedLabel2 != testSampleLabel {
		t.Fatalf("stored label does not match that used during sample creation (%v vs %v)", storedLabel2, testSampleLabel)
	}

	// close the storage
	if err := tmp.Destroy(); err != nil {
		t.Fatal(err)
	}

	// clean up
	os.RemoveAll("./tmp/")
}
