package data

import (
	"testing"
)

// test init
func Test_init(t *testing.T) {

	// check known processes populated
	if len(ProcessRegister) == 0 {
		t.Fatalf("init function did not register any processes")
	}
}
