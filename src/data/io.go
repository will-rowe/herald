// Package data adds some wrapper functions to the protobuf messages
package data

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
)

// InitExperiment will init an experiment struct with the minimum required values
func InitExperiment(label, outputDir, fast5Dir, fastqDir string) *Experiment {

	// create the experiment
	experiment := &Experiment{
		Metadata: &HeraldData{
			Created: ptypes.TimestampNow(),
			Label:   label,
			History: []*Comment{},
			Tags:    []*Process{},
		},
		OutputDirectory:      outputDir,
		Fast5OutputDirectory: fast5Dir,
		FastqOutputDirectory: fastqDir,
	}

	// create the history
	experiment.Metadata.AddComment("experiment created.")

	// return pointer to the experiment
	return experiment
}

// InitSample will init a sample struct with the minimum required values
func InitSample(label, expLabel string, barcode int32) *Sample {

	// create the sample
	sample := &Sample{
		Metadata: &HeraldData{
			Created: ptypes.TimestampNow(),
			Label:   label,
			History: []*Comment{},
			Tags:    []*Process{},
		},
		ExperimentLabel: expLabel,
		Barcode:         barcode,
	}

	// create the history
	sample.Metadata.AddComment("experiment created.")

	// return pointer to the experiment
	return sample
}

// AddComment is a method to add a comment to the history of an experiment or sample
func (heraldData *HeraldData) AddComment(text string) error {
	if len(text) == 0 {
		return fmt.Errorf("no comment provided")
	}
	comment := &Comment{
		Timestamp: ptypes.TimestampNow(),
		Text:      text,
	}
	heraldData.History = append(heraldData.History, comment)
	return nil
}

// AddTags is a method to tag a process to an exeriment or sample
func (heraldData *HeraldData) AddTags(tags []string) error {
	if len(tags) == 0 {
		return fmt.Errorf("no tags provided")
	}

	// range over the tags to be added
	for _, processName := range tags {

		// check the tags are recognised processes
		process, ok := ProcessRegister[processName]
		if !ok {
			return fmt.Errorf("unrecognised process name: %v", processName)
		}

		// make sure the sample is not already tagged with this process
		for _, existingTag := range heraldData.Tags {
			if existingTag.Name == processName {
				return fmt.Errorf("sample already tagged with process: %v", processName)
			}
		}

		// make a copy of the process and tag the sample
		heraldData.Tags = append(heraldData.Tags, process)
		if err := heraldData.AddComment(fmt.Sprintf("tagged with process: %v.", processName)); err != nil {
			return err
		}
	}

	// update the status to "tagged"
	heraldData.Status = 2
	return nil
}

// CheckStatus is a method to check the status of an experiment or sample and update it accordingly
func (experiment *Experiment) CheckStatus() error {

	return nil
}
