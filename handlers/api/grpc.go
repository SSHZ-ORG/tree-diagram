package api

import (
	"github.com/SSHZ-ORG/tree-diagram/pb"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var internalError = status.Error(codes.Unknown, "Internal Server Error")

type treeDiagramService struct {
	pb.UnsafeTreeDiagramServiceServer
}

func GrpcServer() *grpcweb.WrappedGrpcServer {
	s := grpc.NewServer()
	pb.RegisterTreeDiagramServiceServer(s, &treeDiagramService{})
	return grpcweb.WrapServer(s, grpcweb.WithOriginFunc(func(origin string) bool {
		return true
	}))
}
