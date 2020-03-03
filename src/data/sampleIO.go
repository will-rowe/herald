// Package data adds some wrapper functions to the protobuf messages
package data

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
)

// InitSample will init a sample struct with the minimum required values
func InitSample(label, experiment string, barcode int32, comment string) *Sample {

	// create the sample
	sample := &Sample{
		Created:        ptypes.TimestampNow(),
		Label:          label,
		ExperimentName: experiment,
		Status:         1,
		Barcode:        barcode,
		History:        []*Comment{},
		Tags:           &Tags{},
	}

	// create the history and pin any comment
	sample.AddComment("sample created.")
	if len(comment) != 0 {
		sample.AddComment(fmt.Sprintf("user comment: %v", comment))
	}

	// return pointer to the sample
	return sample
}

// AddComment is a method to add a comment to the history of a sample
func (sample *Sample) AddComment(text string) error {
	if len(text) == 0 {
		return fmt.Errorf("no comment provided")
	}
	comment := &Comment{
		Timestamp: ptypes.TimestampNow(),
		Text:      text,
	}
	sample.History = append(sample.History, comment)
	return nil
}

// AddTags is a method to tag a sample
// TODO: this isn't great, and is re-used in the experimentIO, so come up with a better way which will do both
func (sample *Sample) AddTags(tags []string) error {
	if len(tags) == 0 {
		return fmt.Errorf("no tags provided")
	}

	// add all the tags and check for unknown tags
	for _, tag := range tags {
		switch tag {
		case "rampart":
			sample.Tags.Rampart = true
			if err := sample.AddComment("added rampart tag."); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unrecognised tag: %v", tag)
		}
	}

	// update the status to "tagged"
	sample.Status = 2
	return nil
}
