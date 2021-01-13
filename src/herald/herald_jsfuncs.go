package herald

// GetUser returns the name of the user from the config
func (herald *Herald) GetUser() string {
	herald.Lock()
	defer herald.Unlock()
	return herald.config.GetUser().GetName()
}

// GetServerLogfile returns the location of the server logfile
func (herald *Herald) GetServerLogfile() string {
	herald.Lock()
	defer herald.Unlock()
	return herald.config.GetServerlog()
}

// GetDbPath returns the location of the storage on disk
func (herald *Herald) GetDbPath() string {
	herald.Lock()
	defer herald.Unlock()
	return herald.storeLocation
}

// GetRunCount returns the current number of runs in storage
func (herald *Herald) GetRunCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.runCount
}

// GetSampleCount returns the current number of samples in storage
func (herald *Herald) GetSampleCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.sampleCount
}

// GetUntaggedCount returns the current number of runs/samples in storage that are untagged
func (herald *Herald) GetUntaggedCount(descriptor string) int {
	herald.Lock()
	defer herald.Unlock()
	switch descriptor {
	case "runs":
		return herald.untaggedCount[0]
	case "samples":
		return herald.untaggedCount[1]
	}
	return -1
}

// GetTaggedIncompleteCount returns the current number of runs/samples in storage that are tagged with incomplete service requests
func (herald *Herald) GetTaggedIncompleteCount(descriptor string) int {
	herald.Lock()
	defer herald.Unlock()
	switch descriptor {
	case "runs":
		return herald.taggedIncompleteCount[0]
	case "samples":
		return herald.taggedIncompleteCount[1]
	}
	return -1
}

// GetTaggedCompleteCount returns the current number of runs/samples in storage that are tagged with complete service requests
func (herald *Herald) GetTaggedCompleteCount(descriptor string) int {
	herald.Lock()
	defer herald.Unlock()
	switch descriptor {
	case "runs":
		return herald.taggedCompleteCount[0]
	case "samples":
		return herald.taggedCompleteCount[1]
	}
	return -1
}

// GetAnnouncementQueueSize returns the current number of items in the announcment queue
func (herald *Herald) GetAnnouncementQueueSize() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.announcementQueue.Len()
}

// GetAnnouncementCount returns the current number of announcements made
func (herald *Herald) GetAnnouncementCount() int {
	herald.Lock()
	defer herald.Unlock()
	return herald.announcementCount
}

// PrintConfigToJSONstring prints the current in-memory config to a JSON string
func (herald *Herald) PrintConfigToJSONstring() string {
	herald.Lock()
	defer herald.Unlock()

	// TODO: check for error
	return herald.config.GetJSONDump()
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

// GetSampleRun is used by JS to collect a sample run name from the runtime slice of sample data
// NOTE: this assumes the caller has already run GetSampleCount (or similar) to find the iterator range
// TODO: add error on return too (will require re-write of JS function)
func (herald *Herald) GetSampleRun(iterator int) string {
	herald.Lock()
	defer herald.Unlock()
	return herald.sampleDetails[2][iterator]
}

// GetLabel is used by JS to collect an run name from the runtime slice of run names
// NOTE: this assumes the caller has already run GetRunCount (or similar) to find the iterator range
// TODO: add error on return too (will require re-write of JS function)
func (herald *Herald) GetLabel(iterator int) string {
	herald.Lock()
	defer herald.Unlock()
	return herald.runLabels[iterator]
}
