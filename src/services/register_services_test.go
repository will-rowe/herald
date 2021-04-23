package services

import (
	"testing"
)

// test init will run some basic checks on the provided services.
func Test_init(t *testing.T) {

	// check known processes populated
	if len(ServiceRegister) == 0 {
		t.Fatalf("init function did not register any processes")
	}

	// check that the services are offline
	for name, service := range ServiceRegister {
		if service.GetServiceName() == "" {
			t.Fatal("nameless service found")
		}
		if name != service.GetServiceName() {
			t.Fatalf("name mismatch for: %s (%s)", name, service.GetServiceName())
		}
		if service.GetAddress() == "" {
			t.Fatalf("no address found for: %s", name)
		}
	}
}
