// Package minknow is
package minknow

import (
	"google.golang.org/grpc"
)

// CheckAPIstatus will determine if the MinKNOW service is running
func CheckAPIstatus() bool {

	channel, err := grpc.Dial("localhost:9501", grpc.WithInsecure())
	_ = channel
	if err != nil {
		return false
	}
	return true
}
