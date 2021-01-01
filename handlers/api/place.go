package api

import (
	"context"

	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/SSHZ-ORG/tree-diagram/pb"
)

func (t treeDiagramService) RenderPlace(ctx context.Context, req *pb.RenderPlaceRequest) (*pb.RenderPlaceResponse, error) {
	return models.PrepareRenderPlaceResponse(ctx, req.GetId())
}
