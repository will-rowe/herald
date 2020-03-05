package data

import (
	"testing"
)

// test init
func Test_init(t *testing.T) {

	// check known processes populated
	if len(ServiceRegister) == 0 {
		t.Fatalf("init function did not register any processes")
	}

	// try out the dag
	services := make(map[string]bool)
	services["pipelineA"] = false
	services["sequence"] = false
	services["basecall"] = false
	orderedList, err := createServiceDAG(services)
	if err != nil {
		t.Fatal(err)
	}

	// pipelineA depends on sequence and basecall, so should be moved last
	if orderedList[2] != "pipelineA" {
		t.Fatal("service list was not ordered correctly")
	}

}
