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

// Test release version check
//func TestCheckRelease(t *testing.T) {
//	if _, _, _, err := CheckLatestRelease(); err != nil {
//		t.Fatal(err)
//	}
//}
