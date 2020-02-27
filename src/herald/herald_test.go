package herald

import (
	"os"
	"testing"
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

	// create and add a sample
	testLabel := "testLabel"
	if err := tmp.CreateSample(testLabel, 1, "test comment", []string{"sequence"}); err != nil {
		t.Fatal(err)
	}

	// check runtime info was updated
	count, _, _, _ := tmp.GetSampleCounts()
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
