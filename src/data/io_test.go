package data

import (
	"testing"

	"github.com/golang/protobuf/proto"
)

// TestProtobufSample tests the marshalling of a sample
func TestProtobufSample(t *testing.T) {

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
	if test.Metadata.GetLabel() != newTest.Metadata.GetLabel() {
		t.Fatalf("data mismatch %q != %q", test.Metadata.GetLabel(), newTest.Metadata.GetLabel())
	}
}

// TestProtobufExp tests the marshalling of an experiment
func TestProtobufExp(t *testing.T) {

	// set up a basic experiment
	test := InitExperiment("testExperiment", "", "", "")

	// marshal it
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatalf("marshaling error: %v", err)
	}

	// unmarshal it
	newTest := &Experiment{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		t.Fatalf("unmarshaling error: %v", err)
	}

	// test for data match
	if test.Metadata.GetLabel() != newTest.Metadata.GetLabel() {
		t.Fatalf("data mismatch %q != %q", test.Metadata.GetLabel(), newTest.Metadata.GetLabel())
	}
}

// TestTaggingExp tests the tagging of an experiment
func TestTaggingExp(t *testing.T) {

	// set up a basic sample
	test := InitExperiment("testExperiment", "", "", "")

	// check that tags are required to method call
	if err := test.Metadata.AddTags(nil); err == nil {
		t.Fatal("AddTags method did not return error when called with no tags")
	}

	// tag it with an unrecognised tag
	if err := test.Metadata.AddTags([]string{"bogus"}); err == nil {
		t.Fatal("AddTags method did not return error when called with unrecognised tag")
	}

	// tag it with a recognised tag
	if err := test.Metadata.AddTags([]string{"sequence"}); err != nil {
		t.Fatal("AddTags method returned error when called with a recognised tag (sequence)")
	}
}
