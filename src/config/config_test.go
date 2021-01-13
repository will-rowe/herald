package config

import (
	"os"
	"testing"
)

// TestConfig tests the loading of a config from file to memory
func TestConfig(t *testing.T) {

	// load a copy of the config
	c, err := InitConfig(".")
	if err != nil {
		t.Fatal(err)
	}

	// check the config
	if c.GetFileformat() != DefaultConfigType {
		t.Fatal("config not inited with default value")
	}

	// delete the config we made
	if err := os.Remove(c.GetFilepath()); err != nil {
		t.Fatal(err)
	}

}
