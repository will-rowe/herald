package services

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	toposort "github.com/philopon/go-toposort"
)

// InitRun will init a run struct with the minimum required values
func InitRun(owner *User, label, outputDir, fast5Dir, fastqDir string) *Run {

	// create the run
	run := &Run{
		Metadata: &HeraldData{
			Created:      ptypes.TimestampNow(),
			Owner:        owner,
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
	run.Metadata.AddComment("run created.")

	// return pointer to the run
	return run
}

// InitSample will init a sample struct with the minimum required values
func InitSample(sampleLabel, runLabel string, barcode int32) *Sample {

	// create the sample
	sample := &Sample{
		Metadata: &HeraldData{
			Created:      ptypes.TimestampNow(),
			Label:        sampleLabel,
			History:      []*Comment{},
			Status:       1,
			Tags:         make(map[string]bool),
			RequestOrder: []string{},
		},
		ParentRun: runLabel,
		Barcode:   barcode,
	}

	// create the history
	sample.Metadata.AddComment("sample created.")

	// return pointer to the run
	return sample
}

// AddComment is a method to add a comment to the history of an run or sample
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
	if len(heraldData.GetTags()) != 0 {
		return fmt.Errorf("data has already been tagged")
	}

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
	return heraldData.createServiceDAG()
}

// SetTag is a method to set a tag either true or false (complete or incomplete)
func (heraldData *HeraldData) SetTag(serviceName string, value bool) error {

	// check the tag is a recognised service
	if _, ok := ServiceRegister[serviceName]; !ok {
		return fmt.Errorf("unrecognised service name: %v", serviceName)
	}

	// check the data has this tag and set it
	if _, ok := heraldData.Tags[serviceName]; !ok {
		return fmt.Errorf("%v does not have tag: %v", heraldData.GetLabel(), serviceName)
	}
	heraldData.Tags[serviceName] = value
	heraldData.AddComment(fmt.Sprintf("%v tag marked as %v.", serviceName, value))
	return nil
}

// CheckStatus checks the tags and updates the status if all tags are now marked complete
// TODO: this func is incomplete - it only checks for tagged services atm
func (heraldData *HeraldData) CheckStatus() error {
	status := heraldData.GetStatus()
	switch status {

	// uninitialised
	case 0:
		return fmt.Errorf("encountered uninitialised data: %v", heraldData.Label)

	// untagged
	case 1:
		return nil

	// tagged
	case 2:

		// if there is an incomplete tag, nothing to do
		for _, complete := range heraldData.GetTags() {
			if !complete {
				return nil
			}
		}
		// set status to untagged
		heraldData.Status = 1
		return nil

	// announced
	//case 3:
	//	return nil

	// unknown
	default:
		return fmt.Errorf("unknown status: %d", status)
	}
}

// createServiceDAG creates a linear ordering of services that accounts for service dependencies
func (heraldData *HeraldData) createServiceDAG() error {

	// reset requestOrder
	heraldData.RequestOrder = []string{}

	// transfer service names from map to slice
	serviceList := make([]string, len(heraldData.Tags))
	numServices := 0
	for serviceName := range heraldData.GetTags() {
		serviceList[numServices] = serviceName
		numServices++
	}

	// create a dag
	dag := toposort.NewGraph(numServices)

	// create the nodes
	dag.AddNodes(serviceList...)

	// loop through input list and create edges
	for _, serviceName := range serviceList {

		// ignore services with no dependencies
		service := ServiceRegister[serviceName]
		if len(service.GetDeps()) == 0 {
			continue
		}

		// loop over the depencies and draw edges
		for _, dependencyName := range service.GetDeps() {
			dag.AddEdge(dependencyName, serviceName)
		}
	}

	// sort the graph and check for cycles
	result, ok := dag.Toposort()
	if !ok {
		return fmt.Errorf("service dependency cycle detected")
	}
	heraldData.RequestOrder = result
	return nil
}
