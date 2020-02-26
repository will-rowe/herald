// Package server is used to orchestrate the Herald message passing
package server

import (
	"net/http"
)

// Server
type Server struct {
	fieldA string
}

// NetworkActive is a helper function to check outgoing calls can be made
// TODO: this is only temporary as it's not a great check, we need to check API endpoints instead
func NetworkActive() bool {
	if _, err := http.Get("http://google.com/"); err != nil {
		return false
	}
	return true
}