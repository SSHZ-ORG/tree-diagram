package api

import (
	"context"
	"sync"

	"github.com/SSHZ-ORG/tree-diagram/apicache"
	"github.com/SSHZ-ORG/tree-diagram/models"
	"github.com/SSHZ-ORG/tree-diagram/pb"
	"google.golang.org/appengine/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func prepareRenderActor(ctx context.Context, aids []string) (map[string]*pb.RenderActorsResponse_ResponseItem, error) {
	m := make(map[string]*pb.RenderActorsResponse_ResponseItem)

	var missedIDs []string

	fromCache := apicache.GetRenderActor(ctx, aids)
	for _, id := range aids {
		if data, ok := fromCache[id]; ok {
			p := &pb.RenderActorsResponse_ResponseItem{}
			err := proto.Unmarshal(data, p)
			if err == nil {
				m[id] = p
				continue
			} else {
				log.Errorf(ctx, "proto.Unmarshal: %+v", err)
			}
		}

		missedIDs = append(missedIDs, id)
	}

	responses := make([]*pb.RenderActorsResponse_ResponseItem, len(missedIDs))
	errs := make([]error, len(missedIDs))
	wg := sync.WaitGroup{}
	wg.Add(len(missedIDs))

	for i, id := range missedIDs {
		go func(i int, id string) {
			defer wg.Done()

			res, err := models.PrepareRenderActorResponse(ctx, id)
			if err != nil {
				errs[i] = err
				return
			}

			responses[i] = res
		}(i, id)
	}

	wg.Wait()

	toCache := make(map[string][]byte)
	for i, res := range responses {
		if errs[i] != nil {
			return nil, errs[i]
		}

		id := missedIDs[i]
		s, err := proto.Marshal(res)
		if err != nil {
			return nil, err
		}
		toCache[id] = s
		m[id] = res
	}

	apicache.PutRenderActor(ctx, toCache)
	return m, nil
}

func (t treeDiagramService) RenderActors(ctx context.Context, req *pb.RenderActorsRequest) (*pb.RenderActorsResponse, error) {
	if len(req.GetId()) > 100 {
		return nil, status.Error(codes.InvalidArgument, "At most 100 actors can be fetched in one request")
	}

	for _, id := range req.GetId() {
		if id == "" {
			return nil, status.Error(codes.InvalidArgument, "Empty id")
		}
	}

	res, err := prepareRenderActor(ctx, req.GetId())
	if err != nil {
		log.Errorf(ctx, "prepareRenderActor: %+v", err)
		return nil, internalError
	}
	return &pb.RenderActorsResponse{Items: res}, nil
}

func (t treeDiagramService) ListActors(ctx context.Context, req *pb.ListActorsRequest) (*pb.ListActorsResponse, error) {
	for _, id := range req.GetId() {
		if id == "" {
			return nil, status.Error(codes.InvalidArgument, "Empty id")
		}
	}

	m, err := models.GetActorMap(ctx, req.GetId())
	if err != nil {
		log.Errorf(ctx, "models.GetActorMap: %+v", err)
		return nil, internalError
	}

	resp := &pb.ListActorsResponse{Items: make(map[string]*pb.ListActorsResponse_ResponseItem)}
	for _, a := range m {
		resp.GetItems()[a.ID] = &pb.ListActorsResponse_ResponseItem{
			FavoriteCount: proto.Int32(int32(a.LastFavoriteCount)),
		}
	}
	return resp, nil
}
