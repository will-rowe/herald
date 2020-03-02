// Package helpers contains some helper functions for Herald
package helpers

import (
	"fmt"
	"os"
)

// CheckDirExists is used to check a directory exists
func CheckDirExists(dirPath string) error {
	if _, err := os.Stat(dirPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist (%v)", dirPath)
		} else {
			return nil
		}
	}
	return nil
}
