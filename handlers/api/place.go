package api

import (
	"context"

	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/SSHZ-ORG/tree-diagram/pb"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (t treeDiagramService) RenderPlace(ctx context.Context, req *pb.RenderPlaceRequest) (*pb.RenderPlaceResponse, error) {
	id := req.GetId()
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty id")
	}
	resp, err := models.PrepareRenderPlaceResponse(ctx, id)
	if err != nil {
		log.Errorf(ctx, "models.PrepareRenderPlaceResponse: %+v", err)
		return nil, internalError
	}
	return resp, nil
}
