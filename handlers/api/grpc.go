package api

import (
	"github.com/SSHZ-ORG/tree-diagram/pb"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

type treeDiagramService struct {
	pb.UnimplementedTreeDiagramServiceServer
}

func GrpcServer() *grpcweb.WrappedGrpcServer {
	s := grpc.NewServer()
	pb.RegisterTreeDiagramServiceServer(s, &treeDiagramService{})
	return grpcweb.WrapServer(s, grpcweb.WithOriginFunc(func(origin string) bool {
		return true
	}))
}
