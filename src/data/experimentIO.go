// Package data adds some wrapper functions to the protobuf messages
package data

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
)

// InitExperiment will init an experiment struct with the minimum required values
func InitExperiment(name string, outputDir string) *Experiment {

	// create the experiment
	experiment := &Experiment{
		Created:         ptypes.TimestampNow(),
		Name:            name,
		History:         []*Comment{},
		OutputDirectory: outputDir,
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
