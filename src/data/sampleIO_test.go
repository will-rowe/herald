package data

import (
	"testing"

	"github.com/golang/protobuf/proto"
)

// TestProtobuf tests the marshalling of the sample
func TestProtobuf(t *testing.T) {

	// set up a basic sample
	test := InitSample("testSample", "testExperiment", 1)

	// marshal it
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatalf("marshaling error: %v", err)
	}

	// unmarshal it
	newTest := &Sample{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		t.Fatalf("unmarshaling error: %v", err)
	}

	// test for data match
	if test.GetLabel() != newTest.GetLabel() {
		t.Fatalf("data mismatch %q != %q", test.GetLabel(), newTest.GetLabel())
	}
}

// TestTagging tests the tagging of the sample
func TestTagging(t *testing.T) {

	// set up a basic sample
	test := InitSample("test", "testExperiment", 1)

	// check that tags are required to method call
	if err := test.AddTags(nil); err == nil {
		t.Fatal("AddTags method did not return error when called with no tags")
	}

	// tag it with an unrecognised tag
	if err := test.AddTags([]string{"bogus"}); err == nil {
		t.Fatal("AddTags method did not return error when called with unrecognised tag")
	}

	// tag it with a recognised tag
	if err := test.AddTags([]string{"rampart"}); err != nil {
		t.Fatal("AddTags method returned error when called with a recognised tag (rampart)")
	}
}

// TestStatus tests the status assignment for a sample
func TestStatus(t *testing.T) {

	// set up a basic sample
	test := InitSample("test", "testExperiment", 1)

	// check the status is untagged
	if test.GetStatus().String() != "untagged" {
		t.Fatal("initialised sample should have status: untagged")
	}

	// tag it with a recognised tag
	if err := test.AddTags([]string{"rampart"}); err != nil {
		t.Fatal("AddTags method returned error when called with a recognised tag (rampart)")
	}

	// check the status is tagged
	if test.GetStatus().String() != "tagged" {
		t.Fatal("tagged sample should have status: tagged")
	}
}
