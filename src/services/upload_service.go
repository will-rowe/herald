package services

import (
	"context"
	"fmt"
	"log"
)

// SubmitUpload receives an Upload request,
// checks if it has a label and then returns
// a response.
func (x *Server) SubmitUpload(ctx context.Context, submission *UploadRequest) (*UploadResponse, error) {
	log.Printf("received upload request")
	if submission.GetVal1() == "" {
		return nil, fmt.Errorf("no label sent")
	}
	log.Printf("run label: %s", submission.GetVal1())
	return &UploadResponse{Val2: fmt.Sprintf("here is a response for the client regarding run label: %v", submission.GetVal1())}, nil
}
