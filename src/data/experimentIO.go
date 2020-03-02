// Package data adds some wrapper functions to the protobuf messages
package data

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
)

// InitExperiment will init an experiment struct with the minimum required values
func InitExperiment(name, outputDir, fast5Dir, fastqDir string) *Experiment {

	// create the experiment
	experiment := &Experiment{
		Created:              ptypes.TimestampNow(),
		Name:                 name,
		History:              []*Comment{},
		OutputDirectory:      outputDir,
		Fast5OutputDirectory: fast5Dir,
		FastqOutputDirectory: fastqDir,
	}

	// create the history
	experiment.AddComment("experiment created.")

	// return pointer to the experiment
	return experiment
}

// AddComment is a method to add a comment to the history of an experiment
func (experiment *Experiment) AddComment(text string) error {
	if len(text) == 0 {
		return fmt.Errorf("no comment provided")
	}
	comment := &Comment{
		Timestamp: ptypes.TimestampNow(),
		Text:      text,
	}
	experiment.History = append(experiment.History, comment)
	return nil
}

// AddTags is a method to tag an experiment
// TODO: this isn't great, and is re-used in the sampleIO, so come up with a better way which will do both
func (experiment *Experiment) AddTags(tags []string) error {
	if len(tags) == 0 {
		return fmt.Errorf("no tags provided")
	}

	// add all the tags and check for unknown tags
	for _, tag := range tags {
		switch tag {
		case "experiment":
			experiment.Tags.Sequence = true
			if err := experiment.AddComment("added sequence tag."); err != nil {
				return err
			}
		case "basecall":
			experiment.Tags.Rampart = true
			if err := experiment.AddComment("added basecall tag."); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unrecognised tag: %v", tag)
		}
	}

	// update the status to "tagged"
	experiment.Status = 2
	return nil
}
