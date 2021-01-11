package services

import (
	"context"
	"fmt"
	"log"
)

// SubmitMinionPipeline recieves a MinionPipelineRequest,
// checks and sends it to CLIMB, then returns
// a PipelineResponse.
func (x *Server) SubmitMinionPipeline(ctx context.Context, submission *MinionPipelineRequest) (*MinionPipelineResponse, error) {
	log.Printf("received pipeline submission request")
	if submission.GetPipeline() == "" {
		return nil, fmt.Errorf("no pipeline label sent")
	}
	log.Printf("run label: %s", submission.GetParam1())
	return &MinionPipelineResponse{Val: fmt.Sprintf("here is a server response for the client regarding: %v", submission.GetPipeline())}, nil
}
