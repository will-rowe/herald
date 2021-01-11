package server

import (
	"context"
	"fmt"
	"log"

	"github.com/will-rowe/herald/src/services"
)

// SubmitMinionPipeline recieves a MinionPipelineRequest,
// checks and sends it to CLIMB, then returns
// a PipelineResponse.
func (x *HeraldServer) SubmitMinionPipeline(ctx context.Context, submission *services.MinionPipelineRequest) (*services.MinionPipelineResponse, error) {
	log.Printf("received pipeline submission request")
	if submission.GetPipeline() == "" {
		return nil, fmt.Errorf("no pipeline label sent")
	}
	log.Printf("run label: %s", submission.GetParam1())
	return &services.MinionPipelineResponse{Val: fmt.Sprintf("here is a server response for the client regarding: %v", submission.GetPipeline())}, nil
}
