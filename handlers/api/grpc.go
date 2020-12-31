package api

import (
	"context"

	"github.com/SSHZ-ORG/tree-diagram/pb"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/appengine/log"
	"google.golang.org/grpc"
)

type treeDiagramService struct {
	pb.UnimplementedTreeDiagramServiceServer
}

func (t treeDiagramService) Echo(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Infof(ctx, "%+v", req)
	return &pb.EchoResponse{Payload: req.Payload}, nil
}

func GrpcServer() *grpcweb.WrappedGrpcServer {
	s := grpc.NewServer()
	pb.RegisterTreeDiagramServiceServer(s, &treeDiagramService{})
	return grpcweb.WrapServer(s, grpcweb.WithOriginFunc(func(origin string) bool {
		return true
	}))
}
