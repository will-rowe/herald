package helpers

import (
	"testing"
)

// TestHelpers
func TestHelpers(t *testing.T) {

	// CheckDir
	dummyDir := "/badPath/to/nowhere"
	if err := CheckDir(dummyDir); err == nil {
		t.Fatal("CheckDir function allowed a non-existant directory to pass")
	}
}
