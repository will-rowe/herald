package herald

// GetDbPath returns the location of the storage on disk
func (herald *Herald) GetDbPath() string {
	herald.Lock()
	defer herald.Unlock()
	return herald.storeLocation
}

// GetExperimentCount returns the current number of experiments in storage
func (herald *Herald) GetExperimentCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.experimentCount
}

// GetSampleCount returns the current number of samples in storage
func (herald *Herald) GetSampleCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.sampleCount
}

// GetUntaggedRecordCount returns the current number of samples in storage that are untagged
func (herald *Herald) GetUntaggedRecordCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.untaggedRecordCount
}

// GetTaggedRecordCount returns the current number of samples in storage that are tagged with at least one process
func (herald *Herald) GetTaggedRecordCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.taggedRecordCount
}

// GetAnnouncementCount returns the current number of samples that have been announced
func (herald *Herald) GetAnnouncementCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.announcementCount
}

// PrintSampleToJSONstring collects a sample from the database and returns a string of the sample protobuf data in JSON
func (herald *Herald) PrintSampleToJSONstring(label string) string {
	herald.Lock()
	defer herald.Unlock()

	// TODO: check the error from GetSampleJSONDump method
	sampleString, _ := herald.store.GetSampleJSONDump(label)

	return sampleString
}

// GetSampleLabel is used by JS to collect a sample label from the runtime slice of sample data
// NOTE: this assumes the caller has already run GetSampleCount (or similar) to find the iterator range
// TODO: add error on return too (will require re-write of JS function)
func (herald *Herald) GetSampleLabel(iterator int) string {
	herald.Lock()
	defer herald.Unlock()
	return herald.sampleDetails[0][iterator]
}

// GetSampleCreation is used by JS to collect a sample created timestamp from the runtime slice of sample data
// NOTE: this assumes the caller has already run GetSampleCount (or similar) to find the iterator range
// TODO: add error on return too (will require re-write of JS function)
func (herald *Herald) GetSampleCreation(iterator int) string {
	herald.Lock()
	defer herald.Unlock()
	return herald.sampleDetails[1][iterator]
}

// GetSampleExperiment is used by JS to collect a sample experiment name from the runtime slice of sample data
// NOTE: this assumes the caller has already run GetSampleCount (or similar) to find the iterator range
// TODO: add error on return too (will require re-write of JS function)
func (herald *Herald) GetSampleExperiment(iterator int) string {
	herald.Lock()
	defer herald.Unlock()
	return herald.sampleDetails[2][iterator]
}

// GetLabel is used by JS to collect an experiment name from the runtime slice of experiment names
// NOTE: this assumes the caller has already run GetExperimentCount (or similar) to find the iterator range
// TODO: add error on return too (will require re-write of JS function)
func (herald *Herald) GetLabel(iterator int) string {
	herald.Lock()
	defer herald.Unlock()
	return herald.experimentLabels[iterator]
}
