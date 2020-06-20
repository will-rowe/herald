// Package helpers contains some helper functions for Herald
package helpers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/will-rowe/herald/src/version"
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

// CheckLatestRelease will query the GitHub API check the latest released version of Herald against the one in use
// returns a bool indicating if there is a newer version available, the newest release tag, a link to the release, and any error
func CheckLatestRelease() (bool, string, string, error) {
	client := github.NewClient(nil)
	opt := &github.ListOptions{}

	// get the releases via GitHub API
	releases, _, err := client.Repositories.ListReleases(context.Background(), "will-rowe", "herald", opt)
	if err != nil {
		return false, "", "", err
	}

	tagName, url := "", ""
	updateBool := false

	// iterate over the release tags until the most recent, non pre-release is reached
	for _, release := range releases {

		// ignore pre-releases
		if release.GetPrerelease() {
			continue
		}

		// arrived at most recent release
		tagName = release.GetTagName()
		url = release.GetTarballURL()
		if release.GetTagName() != version.VERSION {
			updateBool = true
		}
		break
	}

	return updateBool, tagName, url, nil
}

// NetworkActive is a helper function to check outgoing calls can be made
// TODO: this is only temporary as it's not a great check, we need to check API endpoints instead
func NetworkActive() bool {
	if _, err := http.Get("http://google.com/"); err != nil {
		return false
	}
	return true
}
