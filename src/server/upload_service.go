package server

import (
	"context"
	"fmt"
	"log"

	"github.com/will-rowe/herald/src/services"
)

// SubmitUpload receives an Upload request,
// checks if it has a label and then returns
// a response.
func (heraldServer *HeraldServer) SubmitUpload(ctx context.Context, submission *services.UploadRequest) (*services.UploadResponse, error) {
	log.Printf("received sequencing request")
	if submission.GetVal1() == "" {
		return nil, fmt.Errorf("no label sent")
	}
	log.Printf("run label: %s", submission.GetVal1())
	return &services.UploadResponse{Val2: fmt.Sprintf("here is a response for the client regarding run label: %v", submission.GetVal1())}, nil
}
