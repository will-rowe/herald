// Package helpers contains some helper functions for Herald
package helpers

import (
	"fmt"
	"os"
)

// CheckDirExists is used to check a directory exists
func CheckDirExists(dirPath string) error {
	fh, err := os.Stat(dirPath)
	switch {
	case err != nil:
		return err
	case fh.IsDir():
		return nil
	default:
		return fmt.Errorf("not a directory: %v", dirPath)
	}
}

// DeduplicateStringSlice returns a slice with duplicate entries removed
func DeduplicateStringSlice(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}
