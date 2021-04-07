package services

import (
	"testing"
)

// test init
func Test_init(t *testing.T) {

	// check known processes populated
	if len(ServiceRegister) == 0 {
		t.Fatalf("init function did not register any processes")
	}

	// check that the services are offline
	for name, service := range ServiceRegister {
		if service.CheckAccess() == true {
			t.Fatalf("%s service should be marked offline", name)
		}
	}
}
