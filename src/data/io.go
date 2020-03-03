// Package data adds some wrapper functions to the protobuf messages
package data

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/will-rowe/herald/src/helpers"
)

// InitExperiment will init an experiment struct with the minimum required values
func InitExperiment(label, outputDir, fast5Dir, fastqDir string) *Experiment {

	// create the experiment
	experiment := &Experiment{
		Metadata: &HeraldData{
			Created: ptypes.TimestampNow(),
			Label:   label,
			History: []*Comment{},
			Status:  1,
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
			Status:  1,
			Tags:    []*Process{},
		},
		ParentExperiment: expLabel,
		Barcode:          barcode,
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

// AddTags is a method to tag an exeriment or sample with a process
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

// CheckStatus checks for active tags and then determines if the process endpoints have been reached. It then updates tags and the status
func (experiment *Experiment) CheckStatus() error {
	status := experiment.Metadata.GetStatus()
	switch status {

	// uninitialised
	case 0:
		return fmt.Errorf("experiment is uninitialised: %v", experiment.GetMetadata().Label)

	// untagged
	case 1:
		return nil

	// tagged
	case 2:

		// check all the tagged processes for their endpoints
		for _, tag := range experiment.Metadata.GetTags() {

			// ignore completed tags
			if tag.GetComplete() {
				continue
			}

			// check the endpoint
			// TODO: hardcoding sequence and basecalling basic test for now but I want to automate how we detect endpoints in each process
			if tag.GetName() == "sequence" {
				if err := helpers.CheckDirExists(experiment.GetFast5OutputDirectory()); err != nil {
					tag.Complete = true
				}
			}
			if tag.GetName() == "basecall" {
				if err := helpers.CheckDirExists(experiment.GetFastqOutputDirectory()); err != nil {
					tag.Complete = true
				}
			}
		}

		// loop over tags again and reset status to untagged if all processes finished
		for _, tag := range experiment.Metadata.GetTags() {
			if tag.GetComplete() == false {
				break
			}
			experiment.Metadata.Status = 2
		}

		return nil

	// announced
	//case 3:
	//	return nil

	// unknown
	default:
		return fmt.Errorf("unknown status: %d", status)
	}
}
