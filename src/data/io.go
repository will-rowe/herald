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
			Created:      ptypes.TimestampNow(),
			Label:        label,
			History:      []*Comment{},
			Status:       1,
			Tags:         make(map[string]bool),
			RequestOrder: []string{},
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
			Created:      ptypes.TimestampNow(),
			Label:        label,
			History:      []*Comment{},
			Status:       1,
			Tags:         make(map[string]bool),
			RequestOrder: []string{},
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

// AddTags is a method to tag an exeriment or sample with a service
func (heraldData *HeraldData) AddTags(tags []string) error {
	if len(tags) == 0 {
		return fmt.Errorf("no tags provided")
	}

	// reset requestOrder
	heraldData.RequestOrder = []string{}

	// range over the tags to be added
	for _, serviceName := range tags {

		// check the tag is a recognised service
		_, ok := ServiceRegister[serviceName]
		if !ok {
			return fmt.Errorf("unrecognised service name: %v", serviceName)
		}

		// make sure this data has not already been tagged with this service
		for existingTag := range heraldData.Tags {
			if existingTag == serviceName {
				return fmt.Errorf("data already tagged with service: %v", serviceName)
			}
		}

		// tag the sample
		heraldData.Tags[serviceName] = false
	}

	// update the status to "tagged"
	heraldData.Status = 2

	// generate new request order
	//
	//

	return nil
}

// CheckStatus checks for active tags and then determines if the service endpoints have been reached. It then updates tags and the status
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
		for serviceName, complete := range experiment.Metadata.GetTags() {

			// ignore completed tags
			if complete {
				continue
			}

			// check the endpoint
			// TODO: hardcoding sequence and basecalling basic test for now but I want to automate how we detect endpoints in each service
			if serviceName == "sequence" {
				if err := helpers.CheckDirExists(experiment.GetFast5OutputDirectory()); err != nil {
					experiment.Metadata.Tags[serviceName] = true
				}
			}
			if serviceName == "basecall" {
				if err := helpers.CheckDirExists(experiment.GetFastqOutputDirectory()); err != nil {
					experiment.Metadata.Tags[serviceName] = true
				}
			}
		}

		// loop over tags again and reset status to untagged if all processes finished
		for _, complete := range experiment.Metadata.GetTags() {
			if !complete {
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
