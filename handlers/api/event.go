package api

import (
	"context"

	"cloud.google.com/go/civil"
	"github.com/SSHZ-ORG/tree-diagram/apicache"
	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/SSHZ-ORG/tree-diagram/pb"
	"github.com/SSHZ-ORG/tree-diagram/utils"
	"google.golang.org/appengine/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	queryPageSize = 10
)

func (t treeDiagramService) RenderEvent(ctx context.Context, req *pb.RenderEventRequest) (*pb.RenderEventResponse, error) {
	eid := req.GetId()
	if eid == "" {
		return nil, status.Error(codes.InvalidArgument, "Empty id")
	}

	if fromCache := apicache.GetRenderEvent(ctx, eid); fromCache != nil {
		r := &pb.RenderEventResponse{}
		if err := proto.Unmarshal(fromCache, r); err == nil {
			return r, nil
		} else {
			log.Errorf(ctx, "proto.Unmarshal: %+v", err)
			// Continue below
		}
	}

	res, err := models.PrepareRenderEventResponse(ctx, eid)
	if err == nil {
		if m, err := proto.Marshal(res); err == nil {
			apicache.PutRenderEvent(ctx, eid, m)
		} else {
			panic(err)
		}
	} else {
		log.Errorf(ctx, "models.PrepareRenderEventResponse: %+v", err)
		return nil, internalError
	}

	return res, nil
}

func (t treeDiagramService) QueryEvents(ctx context.Context, req *pb.QueryEventsRequest) (*pb.QueryEventsResponse, error) {
	events, err := models.QueryEvents(ctx, req.GetFilter(), queryPageSize, int(req.GetOffset()))
	if err != nil {
		log.Errorf(ctx, "models.QueryEvents: %+v", err)
		return nil, internalError
	}

	resp := &pb.QueryEventsResponse{}
	for _, e := range events {
		resp.Events = append(resp.Events, &pb.QueryEventsResponse_Event{
			Id:            &e.ID,
			Name:          &e.Name,
			Date:          utils.ToProtoDate(civil.DateOf(e.Date)),
			LastNoteCount: proto.Int32(int32(e.LastNoteCount)),
		})
	}
	return resp, nil
}
