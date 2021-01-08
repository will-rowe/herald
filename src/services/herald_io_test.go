package services

import (
	"testing"

	"github.com/golang/protobuf/proto"
)

// TestProtobufSample tests the marshalling of a sample
func TestProtobufSample(t *testing.T) {

	// set up a basic sample
	test := InitSample("testSample", "testRun", 1)

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

// TestProtobufExp tests the marshalling of an run
func TestProtobufExp(t *testing.T) {

	// set up a basic run
	test := InitRun("testRun", "", "", "")

	// marshal it
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatalf("marshaling error: %v", err)
	}

	// unmarshal it
	newTest := &Run{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		t.Fatalf("unmarshaling error: %v", err)
	}

	// test for data match
	if test.Metadata.GetLabel() != newTest.Metadata.GetLabel() {
		t.Fatalf("data mismatch %q != %q", test.Metadata.GetLabel(), newTest.Metadata.GetLabel())
	}
}

// TestTaggingExp tests the tagging of an run
func TestTaggingExp(t *testing.T) {

	// set up a basic sample
	test := InitRun("testRun", "", "", "")

	// check that tags are required to method call
	if err := test.Metadata.AddTags(nil); err == nil {
		t.Fatal("AddTags method did not return error when called with no tags")
	}

	// tag it with an unrecognised tag
	if err := test.Metadata.AddTags([]string{"bogus"}); err == nil {
		t.Fatal("AddTags method did not return error when called with unrecognised tag")
	}

	// tag it with recognised tags
	if err := test.Metadata.AddTags([]string{"pipelineA", "basecall", "sequence"}); err != nil {
		t.Fatalf("AddTags method returned error when called with recognised tags (sequence, basecall, pipelineA): %v", err)
	}

	// check ordering
	if err := test.Metadata.createServiceDAG(); err != nil {
		t.Fatal(err)
	}

	// pipelineA depends on sequence and basecall, so should be moved last
	if test.GetMetadata().GetRequestOrder()[2] != "pipelineA" {
		t.Fatal("service list was not ordered correctly")
	}
}
