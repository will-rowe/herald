package helpers

import (
	"testing"
)

// TestHelpers
func TestHelpers(t *testing.T) {

	// CheckDir
	dummyDir := "/badPath/to/nowhere"
	if err := CheckDirExists(dummyDir); err == nil {
		t.Fatal("CheckDirExists function allowed a non-existant directory to pass")
	}
}
